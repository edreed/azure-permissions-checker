// +k8s:deepcopy-gen=package
// +groupName=permissions.azure.com

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzPermission enables checking resource and resource group access permissions for the configured Azure Cloud Provider principal.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type AzPermission struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzPermissionSpec    `json:"spec"`
	Status *AzPermissionStatus `json:"status,omitempty"`
}

type AzPermissionSpec struct {
	ResourcePath string `json:"resource_path,omitempty"`
}

type AzPermissionState string

const (
	Pending AzPermissionState = "Pending"
	Ready   AzPermissionState = "Ready"
	Failed  AzPermissionState = "Failed"
)

type AzPermissionStatus struct {
	State        AzPermissionState `json:"state"`
	StateReason  string            `json:"state_reason,omitempty"`
	StateMessage string            `json:"state_message,omitempty"`

	Principal string   `json:"principal,omitempty"`
	Allowed   []string `json:"allowed,omitempty"`
	Denied    []string `json:"denied,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzPermissionList is a list of AzPermission resources.
type AzPermissionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AzPermission `json:"items"`
}
