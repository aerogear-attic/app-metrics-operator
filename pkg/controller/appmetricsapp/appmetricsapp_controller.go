package appmetricsapp

import (
	"context"
	"fmt"
	"os"
	"strings"

	metricsv1alpha1 "github.com/aerogear/app-metrics-operator/pkg/apis/metrics/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_appmetricsapp")

// Add creates a new AppMetricsApp Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAppMetricsApp{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("appmetricsapp-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AppMetricsApp
	err = c.Watch(&source.Kind{Type: &metricsv1alpha1.AppMetricsApp{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources and requeue the owner AppMetricsApp
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsApp{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAppMetricsApp implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAppMetricsApp{}

// ReconcileAppMetricsApp reconciles a AppMetricsApp object
type ReconcileAppMetricsApp struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AppMetricsApp object and makes changes based on the state read
// and what is in the AppMetricsApp.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAppMetricsApp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AppMetricsApp")

	// Fetch the AppMetricsApp instance
	instance := &metricsv1alpha1.AppMetricsApp{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	err = isValidAppNamespace(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	routeList := &routev1.RouteList{}
	listOptions := &client.ListOptions{}

	ns, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return reconcile.Result{}, err
	}

	listOptions.InNamespace(ns)
	listOptions.MatchingLabels(map[string]string{"component": "aerogear-app-metrics"})

	err = r.client.List(context.TODO(), listOptions, routeList)
	if err != nil {
		return reconcile.Result{}, err
	}

	if len(routeList.Items) != 1 {
		err = fmt.Errorf("Found more or less than one Route in the namespace for the metrics service")
		return reconcile.Result{}, err
	}

	// Define a new ConfigMap object
	configMap := newConfigMapForCR(instance, routeList.Items[0].Spec.Host)

	// Set AppMetricsApp instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, configMap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this ConfigMap already exists
	found := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name, "cr.Spec.Name", instance.Spec.Name)
		err = r.client.Create(context.TODO(), configMap)
		if err != nil {
			reqLogger.Info("Error creating the new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name, "Error", err)
			return reconcile.Result{}, err
		}

		// ConfigMap created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// ConfigMap already exists - don't requeue
	reqLogger.Info("Skip reconcile: ConfigMap already exists", "ConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	return reconcile.Result{}, nil
}

func newConfigMapForCR(cr *metricsv1alpha1.AppMetricsApp, host string) *corev1.ConfigMap {
	configmapname := fmt.Sprintf("%s-metrics", cr.Spec.Name)
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configmapname,
			Namespace: cr.Namespace,
		},
		Data: map[string]string{
			"SDKConfig": fmt.Sprintf("{\"url\": \"https://%s\"}", host),
		},
	}
}

// isValidAppNamespace returns an error when the namespace passed is not present in the APP_NAMESPACES environment variable provided to the operator.
func isValidAppNamespace(instance *metricsv1alpha1.AppMetricsApp) error {
	appNamespacesEnvVar, found := os.LookupEnv("APP_NAMESPACES")
	if !found {
		return fmt.Errorf("APP_NAMESPACES environment variable is required for the creation of the app cr")
	}

	for _, ns := range strings.Split(appNamespacesEnvVar, ";") {
		if ns == instance.Namespace {
			return nil
		}
	}
	return fmt.Errorf("The app cr %s was created in a namespace which is not present in the APP_NAMESPACES environment variable provided to the operator", instance.Name)
}
