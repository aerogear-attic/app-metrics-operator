package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppMetricsConfigSpec defines the desired state of AppMetricsConfig
// +k8s:openapi-gen=true
type AppMetricsConfigSpec struct {
	Name string `json:"appName"`
}

// AppMetricsConfigStatus defines the observed state of AppMetricsConfig
// +k8s:openapi-gen=true
type AppMetricsConfigStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppMetricsConfig is the Schema for the appmetricsaconfigs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type AppMetricsConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppMetricsConfigSpec   `json:"spec,omitempty"`
	Status AppMetricsConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppMetricsConfigList contains a list of AppMetricsConfig
type AppMetricsConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppMetricsConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppMetricsConfig{}, &AppMetricsConfigList{})
}
