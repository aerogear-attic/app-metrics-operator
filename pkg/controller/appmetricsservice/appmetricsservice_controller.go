package appmetricsservice

import (
	"context"

	metricsv1alpha1 "github.com/aerogear/app-metrics-operator/pkg/apis/metrics/v1alpha1"
	"github.com/aerogear/app-metrics-operator/pkg/config"

	openshiftappsv1 "github.com/openshift/api/apps/v1"
	imagev1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

var (
	cfg = config.New()
	log = logf.Log.WithName("controller_appmetricsservice")
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AppMetricsService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAppMetricsService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("appmetricsservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AppMetricsService
	err = c.Watch(&source.Kind{Type: &metricsv1alpha1.AppMetricsService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource DeploymentConfig and requeue the owner AppMetricsService
	err = c.Watch(&source.Kind{Type: &openshiftappsv1.DeploymentConfig{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsService{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource ImageStream and requeue the owner AppMetricsService
	err = c.Watch(&source.Kind{Type: &imagev1.ImageStream{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsService{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Secret and requeue the owner AppMetricsService
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsService{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource PersistentVolumeClaim and requeue the owner AppMetricsService
	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsService{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner AppMetricsService
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsService{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource ServiceAccount and requeue the owner AppMetricsService
	err = c.Watch(&source.Kind{Type: &corev1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsService{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Route and requeue the owner AppMetricsService
	err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &metricsv1alpha1.AppMetricsService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAppMetricsService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAppMetricsService{}

// ReconcileAppMetricsService reconciles a AppMetricsService object
type ReconcileAppMetricsService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AppMetricsService object and makes changes based on the state read
// and what is in the AppMetricsService.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAppMetricsService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AppMetricsService")

	// Fetch the AppMetricsService instance
	instance := &metricsv1alpha1.AppMetricsService{}
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

	// look for other appMetricsService resources and don't provision a new one if there is another one with Phase=Complete
	existingInstances := &metricsv1alpha1.AppMetricsServiceList{}
	opts := &client.ListOptions{Namespace: instance.Namespace}
	err = r.client.List(context.TODO(), opts, existingInstances)
	if err != nil {
		reqLogger.Error(err, "Failed to list AppMetricsService resources", "AppMetricsService.Namespace", instance.Namespace)
		return reconcile.Result{}, err
	} else if len(existingInstances.Items) > 1 { // check if > 1 since there's the current one already in that list.
		for _, existingInstance := range existingInstances.Items {
			if existingInstance.Name == instance.Name {
				continue
			}
			if existingInstance.Status.Phase == metricsv1alpha1.PhaseProvision || existingInstance.Status.Phase == metricsv1alpha1.PhaseComplete {
				reqLogger.Info("There is already an AppMetricsService resource in Complete phase. Doing nothing for this CR.", "AppMetricsService.Namespace", instance.Namespace, "AppMetricsService.Name", instance.Name)
				return reconcile.Result{}, nil
			}
		}
	}

	if instance.Status.Phase == metricsv1alpha1.PhaseEmpty {
		instance.Status.Phase = metricsv1alpha1.PhaseProvision
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update AppMetricsService resource status phase", "AppMetricsService.Namespace", instance.Namespace, "AppMetricsService.Name", instance.Name)
			return reconcile.Result{}, err
		}
	}

	//#region Postgres PVC
	persistentVolumeClaim, err := newPostgresqlPersistentVolumeClaim(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set AppMetricsService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, persistentVolumeClaim, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this PersistentVolumeClaim already exists
	foundPersistentVolumeClaim := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: persistentVolumeClaim.Name, Namespace: persistentVolumeClaim.Namespace}, foundPersistentVolumeClaim)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", persistentVolumeClaim.Namespace, "PersistentVolumeClaim.Name", persistentVolumeClaim.Name)
		err = r.client.Create(context.TODO(), persistentVolumeClaim)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	//#region Postgres Secret
	postgresqlSecret, err := newPostgresqlSecret(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set AppMetricsService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, postgresqlSecret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Secret already exists
	foundPostgresqlSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: postgresqlSecret.Name, Namespace: postgresqlSecret.Namespace}, foundPostgresqlSecret)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Secret", "Secret.Namespace", postgresqlSecret.Namespace, "Secret.Name", postgresqlSecret.Name)
		err = r.client.Create(context.TODO(), postgresqlSecret)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	//#region Postgres DeploymentConfig
	postgresqlDeploymentConfig, err := newPostgresqlDeploymentConfig(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set AppMetricsService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, postgresqlDeploymentConfig, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this DeploymentConfig already exists
	foundPostgresqlDeploymentConfig := &openshiftappsv1.DeploymentConfig{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: postgresqlDeploymentConfig.Name, Namespace: postgresqlDeploymentConfig.Namespace}, foundPostgresqlDeploymentConfig)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new DeploymentConfig", "DeploymentConfig.Namespace", postgresqlDeploymentConfig.Namespace, "DeploymentConfig.Name", postgresqlDeploymentConfig.Name)
		err = r.client.Create(context.TODO(), postgresqlDeploymentConfig)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	//#region Postgres Service
	postgresqlService, err := newPostgresqlService(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set AppMetricsService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, postgresqlService, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Service already exists
	foundPostgresqlService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: postgresqlService.Name, Namespace: postgresqlService.Namespace}, foundPostgresqlService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", postgresqlService.Namespace, "Service.Name", postgresqlService.Name)
		err = r.client.Create(context.TODO(), postgresqlService)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	//#region AppMetrics Service
	service, err := newAppMetricsServiceService(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set AppMetricsService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Service already exists
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	//#region AppMetrics Route
	route, err := newAppMetricsServiceRoute(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set AppMetricsService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, route, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Route already exists
	foundRoute := &routev1.Route{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: route.Name, Namespace: route.Namespace}, foundRoute)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		err = r.client.Create(context.TODO(), route)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	//#region AppMetrics ImageStream
	imageStream, err := newAppMetricsServiceImageStream(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set AppMetricsService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, imageStream, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this ImageStream already exists
	foundImageStream := &imagev1.ImageStream{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: imageStream.Name, Namespace: imageStream.Namespace}, foundImageStream)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ImageStream", "ImageStream.Namespace", imageStream.Namespace, "ImageStream.Name", imageStream.Name)
		err = r.client.Create(context.TODO(), imageStream)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	//#region AppMetrics DeploymentConfig
	deploymentConfig, err := newAppMetricsServiceDeploymentConfig(instance)

	if err := controllerutil.SetControllerReference(instance, deploymentConfig, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this DeploymentConfig already exists
	foundDeploymentConfig := &openshiftappsv1.DeploymentConfig{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deploymentConfig.Name, Namespace: deploymentConfig.Namespace}, foundDeploymentConfig)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new DeploymentConfig", "DeploymentConfig.Namespace", deploymentConfig.Namespace, "DeploymentConfig.Name", deploymentConfig.Name)
		err = r.client.Create(context.TODO(), deploymentConfig)
		if err != nil {
			return reconcile.Result{}, err
		}

		// DeploymentConfig created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}
	//#endregion

	if foundDeploymentConfig.Status.ReadyReplicas > 0 && instance.Status.Phase != metricsv1alpha1.PhaseComplete {
		instance.Status.Phase = metricsv1alpha1.PhaseComplete
		r.client.Status().Update(context.TODO(), instance)
	}

	// Resources already exist - don't requeue
	reqLogger.Info("Skip reconcile: resources already exist")
	return reconcile.Result{}, nil
}
