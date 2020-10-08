package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

func init() {
	SchemeBuilder.Register(&GuestCluster{})
	SchemeBuilder.Register(&GuestClusterList{})
}

// +kubebuilder:resource:path=guestclusters,shortName=gc,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=gc;gcs
// GuestCluster is the Schema for the GuestCluster API
type GuestCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestClusterSpec   `json:"spec,omitempty"`
	Status GuestClusterStatus `json:"status,omitempty"`
}

type GuestClusterSpec struct {
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// +optional
	InitialReplicas int32 `json:"initialReplicas,omitempty"`

	// TODO (alberto): populate the API and create/consume infrastructure via aws sdk
	// role profile, sg, vpc, subnets.
}

type GuestClusterStatus struct {
	// Ready is when the GuestClusterStatus has a API server URL.
	// +optional
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// HostedControlPlaneList contains a list of GuestClusters.
type GuestClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GuestCluster `json:"items"`
}
