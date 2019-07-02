package appmetricsservice

import (
	"fmt"
	"strings"

	metricsv1alpha1 "github.com/aerogear/app-metrics-operator/pkg/apis/metrics/v1alpha1"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func labels(cr *metricsv1alpha1.AppMetricsService, suffix string) map[string]string {
	return map[string]string{
		"app":       cr.Name,
		"service":   fmt.Sprintf("%s-%s", cr.Name, suffix),
		"component": "aerogear-app-metrics",
	}
}

// objectMeta returns the default ObjectMeta for all the other objects here
func objectMeta(cr *metricsv1alpha1.AppMetricsService, suffix string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", cr.Name, suffix),
		Namespace: cr.Namespace,
		Labels:    labels(cr, suffix),
	}
}

func generatePassword() (string, error) {
	generatedPassword, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "error generating password")
	}
	return strings.Replace(generatedPassword.String(), "-", "", -1), nil
}
