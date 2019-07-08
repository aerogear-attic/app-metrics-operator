package config

import "os"

type Config struct {
	AppMetricsContainerName string
	PostgresContainerName   string

	AppMetricsImageStreamName string
	AppMetricsImageStreamTag  string

	PostgresImageStreamNamespace string
	PostgresImageStreamName      string
	PostgresImageStreamTag       string

	AppMetricsImageStreamInitialImage string
	PostgresImageStreamInitialImage   string

	BackupImage string
}

func New() Config {
	return Config{
		AppMetricsContainerName: getEnv("APP_METRICS_CONTAINER_NAME", "appmetrics"),
		PostgresContainerName:   getEnv("POSTGRES_CONTAINER_NAME", "postgresql"),

		AppMetricsImageStreamName: getEnv("APP_METRICS_IMAGE_STREAM_NAME", "appmetrics-imagestream"),
		AppMetricsImageStreamTag:  getEnv("APP_METRICS_IMAGE_STREAM_TAG", "0.0"),

		PostgresImageStreamNamespace: getEnv("POSTGRES_IMAGE_STREAM_NAMESPACE", "openshift"),
		PostgresImageStreamName:      getEnv("POSTGRES_IMAGE_STREAM_NAME", "postgresql"),
		// Used both to set the tag, and also for the "POSTGRES_VERSION" in the Secret
		PostgresImageStreamTag: getEnv("POSTGRES_IMAGE_STREAM_TAG", "10"),

		// these are used when the image stream does not exist and created for the first time by the operator
		AppMetricsImageStreamInitialImage: getEnv("APP_METRICS_IMAGE_STREAM_INITIAL_IMAGE", "docker.io/aerogear/aerogear-app-metrics:0.0.13"),

		BackupImage: getEnv("BACKUP_IMAGE", "quay.io/integreatly/backup-container:1.0.8"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
