package contollers

import (
	"context"
	"fmt"

	capiv1 "github.com/openshift-hive/hypershift-installer/pkg/capi/api/v1alpha4"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsprovider/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/yaml"
)

type GuestClusterReconciler struct {
	Client     client.Client
	scheme     *runtime.Scheme
	controller controller.Controller
	recorder   record.EventRecorder
}

func (r *GuestClusterReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	// TODO (alberto): watch hostedControlPlane events too.
	// So when controlPlane.Status.Ready it triggers a reconcile here.
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&capiv1.GuestCluster{}).
		WithOptions(options).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	r.scheme = mgr.GetScheme()
	r.controller = c
	r.recorder = mgr.GetEventRecorderFor("guest-cluster-controller")

	return nil
}

func (r *GuestClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("reconcile")

	// Fetch the GuestCluster instance
	guestCluster := &capiv1.GuestCluster{}
	err := r.Client.Get(ctx, req.NamespacedName, guestCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("GuestCluster not found")
			return ctrl.Result{}, nil
		}
		log.Error(err, "error getting guestCluster")
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, guestCluster.ObjectMeta)
	if err != nil {
		log.Error(err, "error getting owner cluster")
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if util.IsPaused(cluster, guestCluster) {
		log.Info("GuestCluster or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	if !guestCluster.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Implement delete")
		// TODO (alberto): Implement deletion reconciliation loop.
	}

	log = log.WithValues("cluster", cluster.Name)

	patchHelper, err := patch.NewHelper(guestCluster, r.Client)
	if err != nil {
		log.Error(err, "error building patchHelper")
		return ctrl.Result{}, err
	}

	controlPlane := &capiv1.HostedControlPlane{}
	controlPlaneRef := types.NamespacedName{
		Name:      cluster.Spec.ControlPlaneRef.Name,
		Namespace: cluster.Namespace,
	}

	if err := r.Client.Get(ctx, controlPlaneRef, controlPlane); err != nil {
		log.Error(err, "failed to get control plane ref")
		return reconcile.Result{}, err
	}

	if guestCluster.Status.Ready != true {
		// TODO (alberto): populate the API and create/consume infrastructure via aws sdk
		// role profile, sg, vpc, subnets.
		log.Info("guestCLuster not ready yet")
		if !controlPlane.Status.Ready {
			log.Info("Control plane is not ready yet. Requeuing")
			return reconcile.Result{Requeue: true}, nil
		}

		machineSet, err := r.generateMachineset(ctx, *controlPlane, *guestCluster)
		if err != nil {
			log.Error(err, "failed to generate machineSet")
			return reconcile.Result{}, err
		}
		if err := r.Client.Create(ctx, machineSet); err != nil {
			log.Error(err, "failed to create machineSet")
			return reconcile.Result{}, err
		}

		guestCluster.Status.Ready = true
	}

	// Set the values for upper level controller
	// TODO (alberto): Infer this from the hostedControlPlane resource.
	// This is required by the core Cluster resource.
	// https://github.com/kubernetes-sigs/cluster-api/blob/2dd66d6e1559f47302c263cd75bd20265ce6403a/controllers/cluster_controller_phases.go#L183-L188
	guestCluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
		Host: "fake",
		Port: int32(1234),
	}

	if err := patchHelper.Patch(ctx, guestCluster); err != nil {
		log.Error(err, "failed to patch guestCluster")
		return ctrl.Result{}, fmt.Errorf("failed to patch GuestCluster: %w", err)
	}

	log.Info("Successfully reconciled GuestCluster")

	return ctrl.Result{}, nil
}

func (r *GuestClusterReconciler) generateMachineset(ctx context.Context,
	hcp capiv1.HostedControlPlane, guestCluster capiv1.GuestCluster) (*machineapi.MachineSet, error) {
	// TODO (alberto): generate compute nodes with CAPI machineSet
	// resources rather than MAPIError getting list
	var machineSetList machineapi.MachineSetList
	if err := r.Client.List(ctx, &machineSetList); err != nil {
		return nil, err
	}
	if len(machineSetList.Items) < 1 {
		return nil, nil
	}
	ms := machineSetList.Items[0]
	providerConfigRaw := ms.Spec.Template.Spec.ProviderSpec.Value.Raw
	var providerConfig *awsprovider.AWSMachineProviderConfig

	if err := yaml.Unmarshal(providerConfigRaw, &providerConfig); err != nil {
		fmt.Errorf("error unmarshalling providerSpec: %v", err)
	}

	providerConfig.UserDataSecret.Name = fmt.Sprintf("%s-user-data", hcp.GetName())
	machineSet := &machineapi.MachineSet{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      guestCluster.GetName(),
			Namespace: "openshift-machine-api",
		},
		Spec: machineapi.MachineSetSpec{
			Replicas: &guestCluster.Spec.InitialReplicas,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"machine.openshift.io/cluster-api-machineset":   guestCluster.GetName(),
					"machine.openshift.io/cluster-api-machine-role": "compute",
				},
			},
			Template: machineapi.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"machine.openshift.io/cluster-api-machineset":   guestCluster.GetName(),
						"machine.openshift.io/cluster-api-machine-role": "compute",
					},
				},
				Spec: machineapi.MachineSpec{
					ProviderSpec: machineapi.ProviderSpec{
						Value: &runtime.RawExtension{Object: providerConfig},
					},
				},
			},
		},
	}
	return machineSet, nil
}
