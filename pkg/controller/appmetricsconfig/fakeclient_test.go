package appmetricsconfig


import (
	metricsv1alpha1 "github.com/aerogear/app-metrics-operator/pkg/apis/metrics/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	routev1 "github.com/openshift/api/route/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

//buildReconcileWithFakeClientWithMocks return reconcile with fake client, schemes and mock objects
func buildReconcileWithFakeClientWithMocks(objs []runtime.Object, t *testing.T) *ReconcileAppMetricsConfig {
	s := scheme.Scheme

	// Add route Openshift scheme
	if err := routev1.AddToScheme(s); err != nil {
		t.Fatalf("Unable to add route scheme: (%v)", err)
	}

	s.AddKnownTypes(metricsv1alpha1.SchemeGroupVersion, &metricsv1alpha1.AppMetricsConfig{})

	// create a fake client to mock API calls with the mock objects
	cl := fake.NewFakeClient(objs...)

	// create a ReconcileMobileSecurityService object with the scheme and fake client
	return &ReconcileAppMetricsConfig{client: cl, scheme: s}
}

