package appmetricsconfig

import (
	metricsv1alpha1 "github.com/aerogear/app-metrics-operator/pkg/apis/metrics/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	appMetricsConfigInstance = metricsv1alpha1.AppMetricsConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-app",
			Namespace: "app-metrics",
		},
	}
)
