module github.com/openshift-hive/hypershift-installer

replace (
	github.com/apcera/gssapi => github.com/openshift/gssapi v0.0.0-20161010215902-5fb4217df13b
	github.com/containers/image => github.com/openshift/containers-image v0.0.0-20190130162827-4bc6d24282b1
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0 // Pin non-versioned import to v22.0.0
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20191121165722-d1d5f6476656
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20200206190020-71b826cc0f0a // Use OpenShift fork
	github.com/metal3-io/cluster-api-provider-baremetal => github.com/openshift/cluster-api-provider-baremetal v0.0.0-20190821174549-a2a477909c1d // Pin OpenShift fork
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200930075302-db52bc4ef99f
	//github.com/openshift/api => github.com/openshift/api v0.0.0-20200526144822-34f54f12813a
	github.com/openshift/oc => github.com/openshift/oc v0.0.0-alpha.0.0.20200527150310-06cc74b10770
	github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.41.1-openshift-3 // Pin to openshift fork with IPv6 fixes
	k8s.io/api => k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.2
	k8s.io/apiserver => k8s.io/apiserver v0.19.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.2
	k8s.io/client-go => k8s.io/client-go v0.19.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.2
	k8s.io/code-generator => k8s.io/code-generator v0.19.2
	k8s.io/component-base => k8s.io/component-base v0.19.2
	k8s.io/cri-api => k8s.io/cri-api v0.19.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.2
	k8s.io/kubectl => k8s.io/kubectl v0.19.2
	k8s.io/kubelet => k8s.io/kubelet v0.19.2
	k8s.io/kubernetes => k8s.io/kubernetes v1.18.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.2
	k8s.io/metrics => k8s.io/metrics v0.19.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.2
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20200506073438-9d49428ff837 // Pin OpenShift fork
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20200120114645-8a9592f1f87b // Pin OpenShift fork
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20200323110431-3311de91e078 // Pin OpenShift fork
)

go 1.13

require (
	github.com/Luzifer/go-dhparam v1.0.0
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/aws/aws-sdk-go v1.35.0 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.2.1
	github.com/google/uuid v1.1.2
	github.com/krishicks/yaml-patch v0.0.10
	github.com/openshift/api v3.9.1-0.20191111211345-a27ff30ebf09+incompatible
	github.com/openshift/client-go v0.0.0-20200827190008-3062137373b5
	github.com/openshift/cluster-api v0.0.0-20191129101638-b09907ac6668
	github.com/openshift/ibm-roks-toolkit v0.0.0-20201028154754-45600a282504
	github.com/openshift/installer v0.9.0-master.0.20200528191708-78bc158696f3
	github.com/openshift/oc v0.0.0-alpha.0.0.20200930115932-6bfa7f97c8d4
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/cli-runtime v0.19.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.19.2
	sigs.k8s.io/cluster-api v0.3.11-0.20201029150259-71ba9246f9cd
	sigs.k8s.io/cluster-api-provider-aws v0.0.0
	sigs.k8s.io/controller-runtime v0.7.0-alpha.3
	sigs.k8s.io/yaml v1.2.0
)
