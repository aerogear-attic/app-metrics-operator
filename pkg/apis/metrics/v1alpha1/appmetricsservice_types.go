package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppMetricsServiceSpec defines the desired state of AppMetricsService
// +k8s:openapi-gen=true
type AppMetricsServiceSpec struct {
}

// AppMetricsServiceStatus defines the observed state of AppMetricsService
// +k8s:openapi-gen=true
type AppMetricsServiceStatus struct {
	Phase StatusPhase `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppMetricsService is the Schema for the appmetricsservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type AppMetricsService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppMetricsServiceSpec   `json:"spec,omitempty"`
	Status AppMetricsServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppMetricsServiceList contains a list of AppMetricsService
type AppMetricsServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppMetricsService `json:"items"`
}

type StatusPhase string

var (
	PhaseEmpty     StatusPhase = ""
	PhaseComplete  StatusPhase = "Complete"
	PhaseProvision StatusPhase = "Provisioning"
)

func init() {
	SchemeBuilder.Register(&AppMetricsService{}, &AppMetricsServiceList{})
}
