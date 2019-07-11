package appmetricsservice

import (
	"fmt"

	metricsv1alpha1 "github.com/aerogear/app-metrics-operator/pkg/apis/metrics/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func backups(cr *metricsv1alpha1.AppMetricsService) ([]batchv1beta1.CronJob, error) {
	cronjobs := []batchv1beta1.CronJob{}
	for _, backup := range cr.Spec.Backups {
		cronJobLabels := labels(cr, "backup")
		jobLabels := cronJobLabels
		jobLabels["cronjob-name"] = backup.Name

		cronjobs = append(cronjobs, batchv1beta1.CronJob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      backup.Name,
				Namespace: cr.Namespace,
				Labels:    cronJobLabels,
			},
			Spec: batchv1beta1.CronJobSpec{
				Schedule: backup.Schedule,
				JobTemplate: batchv1beta1.JobTemplateSpec{
					Spec: batchv1.JobSpec{
						Template: corev1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: jobLabels,
							},
							Spec: corev1.PodSpec{
								// This SA needs to be created beforehand
								// https://github.com/integr8ly/backup-container-image/tree/master/templates/openshift/rbac
								ServiceAccountName: "backupjob",
								Containers: []corev1.Container{
									{
										Name:            backup.Name + "-appmetrics-backup",
										Image:           cfg.BackupImage,
										ImagePullPolicy: "Always",
										Command:         buildBackupContainerCommand(backup, cr.Namespace),
										Env:             buildBackupCronJobEnvVars(backup, cr.Name, cr.Namespace),
									},
								},
								RestartPolicy: corev1.RestartPolicyOnFailure,
							},
						},
					},
				},
			},
		})
	}
	return cronjobs, nil
}

func buildBackupContainerCommand(backup metricsv1alpha1.AppMetricsServiceBackup, crNamespace string) []string {
	command := []string{"/opt/intly/tools/entrypoint.sh", "-c", "postgres", "-n", crNamespace}

	// If there is no encryption secret, we need to inhibit the
	// encryption behaviour
	if backup.EncryptionKeySecretName == "" {
		command = append(command, "-e", "")
	}

	return command
}

func buildBackupCronJobEnvVars(backup metricsv1alpha1.AppMetricsServiceBackup, crName string, crNamespace string) []corev1.EnvVar {

	envVars := []corev1.EnvVar{
		{
			Name:  "PRODUCT_NAME",
			Value: "app-metrics",
		},
		{
			Name:  "COMPONENT_SECRET_NAME",
			Value: fmt.Sprintf("%s-%s", crName, "postgresql"),
		},
		{
			Name:  "COMPONENT_SECRET_NAMESPACE",
			Value: crNamespace,
		},
	}

	backendSecretNamespace := backup.BackendSecretNamespace
	if backendSecretNamespace == "" {
		backendSecretNamespace = crNamespace
	}

	encryptionKeySecretNamespace := backup.EncryptionKeySecretNamespace
	if encryptionKeySecretNamespace == "" {
		encryptionKeySecretNamespace = crNamespace
	}

	if backup.BackendSecretName != "" {
		envVars = append(envVars,
			corev1.EnvVar{
				Name:  "BACKEND_SECRET_NAME",
				Value: backup.BackendSecretName,
			},
			corev1.EnvVar{
				Name:  "BACKEND_SECRET_NAMESPACE",
				Value: backendSecretNamespace,
			},
		)
	}

	if backup.EncryptionKeySecretName != "" {
		envVars = append(envVars,
			corev1.EnvVar{
				Name:  "ENCRYPTION_SECRET_NAME",
				Value: backup.EncryptionKeySecretName,
			},
			corev1.EnvVar{
				Name:  "ENCRYPTION_SECRET_NAMESPACE",
				Value: encryptionKeySecretNamespace,
			},
		)
	}

	return envVars
}
