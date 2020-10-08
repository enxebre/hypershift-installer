package contollers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	capiv1 "github.com/openshift-hive/hypershift-installer/pkg/capi/api/v1alpha4"
	"github.com/openshift-hive/hypershift-installer/pkg/installer"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

type HostedControlPlaneReconciler struct {
	Client     client.Client
	Log        logr.Logger
	scheme     *runtime.Scheme
	controller controller.Controller
	recorder   record.EventRecorder
}

func (r *HostedControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&capiv1.HostedControlPlane{}).
		WithOptions(options).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	r.scheme = mgr.GetScheme()
	r.controller = c
	r.recorder = mgr.GetEventRecorderFor("hosted-control-plane-controller")

	return nil
}

func (r *HostedControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("reconcile")

	// Fetch the hostedControlPlane instance
	hostedControlPlane := &capiv1.HostedControlPlane{}
	err := r.Client.Get(ctx, req.NamespacedName, hostedControlPlane)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, hostedControlPlane.ObjectMeta)
	if err != nil {
		//log.WithValues("Eerror", err)
		return ctrl.Result{}, err
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if util.IsPaused(cluster, hostedControlPlane) {
		log.Info("HostedControlPlane or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	if !hostedControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Delete")
		// TODO (alberto): Implement deletion reconciliation loop.
	}

	log = log.WithValues("cluster", cluster.Name)

	patchHelper, err := patch.NewHelper(hostedControlPlane, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	if hostedControlPlane.Status.Ready != true {
		log.Info("HostedControlPlane is not ready. Reconciling")
		// TODO (alberto): move installer.Reconcile into this package
		//
		// May be eventually just run a deployment with a CVO running a hostedControlPlane profile
		// passing the hostedControlPlane.spec.version through?
		if err := installer.Reconcile(r.Client, req); err != nil {
			log.Error(err, "error reconciling")
			return ctrl.Result{}, fmt.Errorf("failed to reconcile hypershift: %w", err)
		}
	}

	// Set the values for upper level controller
	hostedControlPlane.Status.Ready = true

	if err := patchHelper.Patch(ctx, hostedControlPlane); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to patch GuestCluster: %w", err)
	}

	log.Info("Successfully reconciled HostedControlPlane")

	return ctrl.Result{}, nil
}
