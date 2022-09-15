//
// Copyright (c) 2022 IBM Corporation
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   IBM Corporation - initial API and implementation
//

package wazi

import (
	"context"
	"time"

	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	"github.com/eclipse-che/che-operator/pkg/deploy"

	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	orgv1 "github.com/eclipse-che/che-operator/api/v1"
)

// CheClusterReconciler reconciles a CheCluster object
type WaziLicenseReconciler struct {
	Log    logr.Logger
	Scheme *k8sruntime.Scheme

	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client

	// This client, is a simple client
	// that reads objects without using the cache,
	// to simply read objects thta we don't intend
	// to further watch
	nonCachedClient client.Client
	// A discovery client to check for the existence of certain APIs registered
	// in the API Server
	discoveryClient  discovery.DiscoveryInterface
	tests            bool
	reconcileManager *deploy.ReconcileManager
	// the namespace to which to limit the reconciliation. If empty, all namespaces are considered
	namespace string
}

// NewReconciler returns a new WaziLicenseReconciler
func NewReconciler(
	k8sclient client.Client,
	noncachedClient client.Client,
	discoveryClient discovery.DiscoveryInterface,
	scheme *k8sruntime.Scheme,
	namespace string) *WaziLicenseReconciler {

	reconcileManager := deploy.NewReconcileManager()

	return &WaziLicenseReconciler{
		Scheme: scheme,
		Log:    ctrl.Log.WithName("controllers").WithName("WaziLicense"),

		client:           k8sclient,
		nonCachedClient:  noncachedClient,
		discoveryClient:  discoveryClient,
		namespace:        namespace,
		reconcileManager: reconcileManager,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *WaziLicenseReconciler) SetupWithManager(mgr ctrl.Manager) error {

	controllerBuilder := ctrl.NewControllerManagedBy(mgr).

		// Watch for changes to primary resource WaziLicense
		Watches(&source.Kind{Type: &orgv1.WaziLicense{}}, &handler.EnqueueRequestForObject{}).
		// Watch for changes to secondary resources and requeue the owner WaziLicese
		Watches(&source.Kind{Type: &odlmv1alpha1.OperandRequest{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &orgv1.WaziLicense{},
		})

	if r.namespace != "" {
		controllerBuilder = controllerBuilder.WithEventFilter(util.InNamespaceEventFilter(r.namespace))
	}

	return controllerBuilder.
		For(&orgv1.WaziLicense{}).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.5/pkg/reconcile
func (r *WaziLicenseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("wazilicense", req.NamespacedName)

	clusterAPI := deploy.ClusterAPI{
		Client:           r.client,
		NonCachingClient: r.nonCachedClient,
		DiscoveryClient:  r.discoveryClient,
		Scheme:           r.Scheme,
	}

	// Fetch the WaziLicense instance
	wazilicense, err := r.GetWaziCR(req)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	deployContext := &deploy.DeployContext{
		ClusterAPI:  clusterAPI,
		WaziLicense: wazilicense,
	}

	// Note: Proxy Configuration is for cluster wide operations

	if deployContext.WaziLicense.ObjectMeta.DeletionTimestamp.IsZero() {
		result, done, err := r.reconcileManager.ReconcileAll(deployContext)
		if !done {
			return result, err
		}
	} else {
		r.reconcileManager.FinalizeAll(deployContext)
	}

	// Reconcile finalizers before CR is deleted
	r.reconcileFinalizers(deployContext)

	// Reconcile Permissions
	done, err := r.reconcileWaziLicensePermissions(deployContext)
	if !done {
		if err != nil {
			_ = deploy.SetWaziLicenseStatusFailedPermissions(deployContext)
			logrus.Error(err)
		}
		// reconcile after 1 seconds since we deal with cluster objects
		return ctrl.Result{RequeueAfter: time.Second}, err
	}

	if err := r.GenerateAndSaveFields(deployContext); err != nil {
		_ = deploy.ReloadWaziLicenseCR(deployContext)
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 1}, err
	}

	// Reconcile Annotations
	done, err = r.ReconcileWaziAnnotations(deployContext)
	if !done {
		if err != nil {
			logrus.Error(err)
			_ = deploy.SetWaziLicenseStatusFailedAnnotations(deployContext)
			return ctrl.Result{RequeueAfter: time.Second}, err
		}
	}

	// Reconcile Operand Request
	done, err = r.ReconcileWaziOperandRequest(deployContext)
	if !done {
		if err != nil {
			_ = deploy.SetWaziLicenseStatusFailedODLM(deployContext)
			logrus.Error(err)
		}
		return ctrl.Result{}, nil // Reconcile will fire again when Operand Request changes
	}

	done, err = r.ReconcileWaziLicensingQuerySource(deployContext)
	if !done {
		if err != nil {
			_ = deploy.SetWaziLicenseStatusFailedLicensingQuerySource(deployContext)
			logrus.Error(err)
		}
		return ctrl.Result{}, nil
	}

	_ = deploy.SetWaziLicenseStatusSuccess(deployContext)

	return ctrl.Result{}, nil
}

func (r *WaziLicenseReconciler) reconcileFinalizers(deployContext *deploy.DeployContext) {

	if _, err := r.reconcileWaziLicensePermissionsFinalizers(deployContext); err != nil {
		logrus.Error(err)
	}

	if err := deploy.ReconcileOperandRequestFinalizer(deployContext); err != nil {
		logrus.Error(err)
	}

	if err := deploy.ReconcileLicensingQuerySourceFinalizer(deployContext); err != nil {
		logrus.Error(err)
	}
}

func (r *WaziLicenseReconciler) GetWaziCR(request ctrl.Request) (instance *orgv1.WaziLicense, err error) {
	instance = &orgv1.WaziLicense{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		r.Log.Error(err, "Failed to get %s CR: %s", "Cluster name", instance.Name)
		return nil, err
	}
	return instance, nil
}
