package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&HostedControlPlane{})
	SchemeBuilder.Register(&HostedControlPlaneList{})
}

// +kubebuilder:resource:path=hostedcontrolplanes,shortName=hcp,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=hcp;hcps
// HostedControlPlane is the Schema for the KubeadmControlPlane API.
type HostedControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HostedControlPlaneSpec   `json:"spec,omitempty"`
	Status HostedControlPlaneStatus `json:"status,omitempty"`
}

type HostedControlPlaneSpec struct {
	Name *string `json:"name,omitempty"`

	// SSHKey is the public Secure Shell (SSH) key to provide access to instances.
	// +optional
	SSHKey string `json:"sshKey,omitempty"`

	BaseDomain string `json:"baseDomain"`

	PullSecret string `json:"pullSecret,omitempty"`

	*Networking `json:"networking,omitempty"`

	// Version defines the desired Kubernetes version. If no version number
	// is supplied then the latest version of Kubernetes that EKS supports
	// will be used.
	// +kubebuilder:validation:MinLength:=2
	// +optional
	Version *string `json:"version,omitempty"`
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	// NetworkType is the type of network to install. The default is OpenShiftSDN
	//
	// +kubebuilder:default=OpenShiftSDN
	// +optional
	NetworkType string `json:"networkType,omitempty"`

	// ClusterNetwork is the list of IP address pools for pods.
	// Default is 10.128.0.0/14 and a host prefix of /23.
	//
	// +optional
	ClusterNetwork []ClusterNetworkEntry `json:"clusterNetwork,omitempty"`

	// ServiceNetwork is the list of IP address pools for services.
	// Default is 172.30.0.0/16.
	// NOTE: currently only one entry is supported.
	//
	// +kubebuilder:validation:MaxItems=1
	// +optional
	ServiceNetwork []string `json:"serviceNetwork,omitempty"`
}

// ClusterNetworkEntry is a single IP address block for pod IP blocks. IP blocks
// are allocated with size 2^HostSubnetLength.
type ClusterNetworkEntry struct {
	// CIDR is the IP block address pool.
	CIDR string `json:"cidr"`

	// HostPrefix is the prefix size to allocate to each node from the CIDR.
	// For example, 24 would allocate 2^8=256 adresses to each node.
	HostPrefix int32 `json:"hostPrefix"`

	// The size of blocks to allocate from the larger pool.
	// This is the length in bits - so a 9 here will allocate a /23.
	// +optional
	DeprecatedHostSubnetLength int32 `json:"hostSubnetLength,omitempty"`
}

type HostedControlPlaneStatus struct {
	// Ready denotes that the HostedControlPlane API Server is ready to
	// receive requests
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// ErrorMessage indicates that there is a terminal problem reconciling the
	// state, and will be set to a descriptive error message.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// HostedControlPlaneList contains a list of HostedControlPlanes.
type HostedControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HostedControlPlane `json:"items"`
}
