package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppMetricsServiceSpec defines the desired state of AppMetricsService
// +k8s:openapi-gen=true
type AppMetricsServiceSpec struct {
	// Backups is an array of configs that will be used to create CronJob resource instances
	Backups []AppMetricsServiceBackup `json:"backups,omitempty"`
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

// Backup contains the info needed to configure a CronJob for backups
type AppMetricsServiceBackup struct {
	// Name is the name that will be given to the resulting
	// CronJob
	Name string `json:"name"`

	// Schedule is the schedule that the job will be run at, in
	// cron format
	Schedule string `json:"schedule"`

	// EncryptionKeySecretName is the name of a secret containing
	// PGP/GPG details, including "GPG_PUBLIC_KEY",
	// "GPG_TRUST_MODEL", and "GPG_RECIPIENT"
	EncryptionKeySecretName string `json:"encryptionKeySecretName,omitempty"`

	// EncryptionKeySecretNamespace is the name of the namespace
	// that the secret referenced in EncryptionKeySecretName
	// resides in
	EncryptionKeySecretNamespace string `json:"encryptionKeySecretNamespace,omitempty"`

	// BackendSecretName is the name of a secret containing
	// storage backend details, such as "AWS_S3_BUCKET_NAME",
	// "AWS_ACCESS_KEY_ID", and "AWS_SECRET_ACCESS_KEY"
	BackendSecretName string `json:"backendSecretName"`

	// BackendSecretNamespace is the name of the namespace that
	// the secret referenced in BackendSecretName resides in
	BackendSecretNamespace string `json:"backendSecretNamespace,omitempty"`
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
