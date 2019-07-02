package appmetricsservice

import (
	"fmt"

	metricsv1alpha1 "github.com/aerogear/app-metrics-operator/pkg/apis/metrics/v1alpha1"
	openshiftappsv1 "github.com/openshift/api/apps/v1"
	imagev1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newAppMetricsServiceService(cr *metricsv1alpha1.AppMetricsService) (*corev1.Service, error) {
	serviceObjectMeta := objectMeta(cr, "appmetrics")
	serviceObjectMeta.Annotations = map[string]string{
		"org.aerogear.metrics/plain_endpoint": "/rest/prometheus/metrics",
	}
	serviceObjectMeta.Labels["mobile"] = "enabled"

	return &corev1.Service{
		ObjectMeta: serviceObjectMeta,
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":     cr.Name,
				"service": "appmetrics",
			},
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:     "web",
					Protocol: corev1.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 3000,
					},
				},
			},
		},
	}, nil
}

func newAppMetricsServiceRoute(cr *metricsv1alpha1.AppMetricsService) (*routev1.Route, error) {
	return &routev1.Route{
		ObjectMeta: objectMeta(cr, "appmetrics"),
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: fmt.Sprintf("%s-%s", cr.Name, "appmetrics"),
			},
			TLS: &routev1.TLSConfig{
				Termination:                   routev1.TLSTerminationEdge,
				InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyNone,
			},
		},
	}, nil
}

func newAppMetricsServiceImageStream(cr *metricsv1alpha1.AppMetricsService) (*imagev1.ImageStream, error) {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      cfg.AppMetricsImageStreamName,
			Labels:    labels(cr, cfg.AppMetricsImageStreamName),
		},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				{
					Name: cfg.AppMetricsImageStreamTag,
					From: &corev1.ObjectReference{
						Kind: "DockerImage",
						Name: cfg.AppMetricsImageStreamInitialImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						Scheduled: false,
					},
				},
			},
		},
	}, nil
}

func newAppMetricsServiceDeploymentConfig(cr *metricsv1alpha1.AppMetricsService) (*openshiftappsv1.DeploymentConfig, error) {
	labels := map[string]string{
		"app":       cr.Name,
		"service":   "appmetrics",
		"component": "aerogear-app-metrics",
	}

	return &openshiftappsv1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: openshiftappsv1.DeploymentConfigSpec{
			Replicas: 1,
			Selector: labels,
			Triggers: openshiftappsv1.DeploymentTriggerPolicies{
				openshiftappsv1.DeploymentTriggerPolicy{
					Type: openshiftappsv1.DeploymentTriggerOnConfigChange,
				},
				openshiftappsv1.DeploymentTriggerPolicy{
					Type: openshiftappsv1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &openshiftappsv1.DeploymentTriggerImageChangeParams{
						Automatic:      true,
						ContainerNames: []string{cfg.AppMetricsContainerName},
						From: corev1.ObjectReference{
							Kind: "ImageStreamTag",
							Name: cfg.AppMetricsImageStreamName + ":" + cfg.AppMetricsImageStreamTag,
						},
					},
				},
				openshiftappsv1.DeploymentTriggerPolicy{
					Type: openshiftappsv1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &openshiftappsv1.DeploymentTriggerImageChangeParams{
						Automatic:      true,
						ContainerNames: []string{cfg.PostgresContainerName},
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Namespace: cfg.PostgresImageStreamNamespace,
							Name:      cfg.PostgresImageStreamName + ":" + cfg.PostgresImageStreamTag,
						},
					},
				},
			},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:            cfg.PostgresContainerName,
							Image:           cfg.PostgresImageStreamName + ":" + cfg.PostgresImageStreamTag,
							ImagePullPolicy: corev1.PullAlways,
							Env: []corev1.EnvVar{
								{
									Name:  "POSTGRES_SERVICE_HOST",
									Value: fmt.Sprintf("%s-postgresql", cr.Name),
								},
							},
							Command: []string{
								"/bin/sh",
								"-c",
								"source /opt/rh/rh-postgresql10/enable && until pg_isready -h $POSTGRES_SERVICE_HOST; do echo waiting for database; sleep 2; done;",
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            cfg.AppMetricsContainerName,
							Image:           cfg.AppMetricsImageStreamName + ":" + cfg.AppMetricsImageStreamTag,
							ImagePullPolicy: corev1.PullAlways,
							Env: []corev1.EnvVar{
								{
									Name:  "PGHOST",
									Value: fmt.Sprintf("%s-postgresql", cr.Name),
								},
								{
									Name: "PGUSER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: "POSTGRES_USERNAME",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: fmt.Sprintf("%s-postgresql", cr.Name),
											},
										},
									},
								},
								{
									Name: "PGPASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: "POSTGRES_PASSWORD",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: fmt.Sprintf("%s-postgresql", cr.Name),
											},
										},
									},
								},
								{
									Name: "PGDATABASE",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: "POSTGRES_DATABASE",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: fmt.Sprintf("%s-postgresql", cr.Name),
											},
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          cfg.AppMetricsContainerName,
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 3000,
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3000,
										},
									},
								},
								InitialDelaySeconds: 15,
								TimeoutSeconds:      1,
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 3000,
										},
									},
								},
								InitialDelaySeconds: 30,
								TimeoutSeconds:      1,
							},
						},
					},
				},
			},
		},
	}, nil
}
