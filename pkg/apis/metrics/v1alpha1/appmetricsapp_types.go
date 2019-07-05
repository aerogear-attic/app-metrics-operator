package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppMetricsAppSpec defines the desired state of AppMetricsApp
// +k8s:openapi-gen=true
type AppMetricsAppSpec struct {
	Name string `json:"appName"`
}

// AppMetricsAppStatus defines the observed state of AppMetricsApp
// +k8s:openapi-gen=true
type AppMetricsAppStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppMetricsApp is the Schema for the appmetricsapps API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type AppMetricsApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppMetricsAppSpec   `json:"spec,omitempty"`
	Status AppMetricsAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppMetricsAppList contains a list of AppMetricsApp
type AppMetricsAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppMetricsApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppMetricsApp{}, &AppMetricsAppList{})
}
