//
// Copyright (c) 2023 IBM Corporation
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

	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	"github.com/devfile/devworkspace-operator/pkg/infrastructure"
	chev2 "github.com/eclipse-che/che-operator/api/v2"
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/eclipse-che/che-operator/pkg/common/test"
	"github.com/eclipse-che/che-operator/pkg/common/utils"
	"github.com/eclipse-che/che-operator/pkg/deploy/wazi"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type WaziLicenseReconciler struct {
	Log              logr.Logger
	Scheme           *k8sruntime.Scheme
	client           client.Client
	nonCachedClient  client.Client
	discoveryClient  discovery.DiscoveryInterface
	reconcileManager *wazi.WaziReconcileManager
	namespace        string
}

func NewReconciler(
	k8sclient client.Client,
	noncachedClient client.Client,
	discoveryClient discovery.DiscoveryInterface,
	scheme *k8sruntime.Scheme,
	namespace string) *WaziLicenseReconciler {

	reconcileManager := wazi.NewReconcileManager()

	if !test.IsTestMode() {
		reconcileManager.RegisterReconciler(NewWaziLicenseValidator())
	}

	if infrastructure.IsOpenShift() {
		reconcileManager.RegisterReconciler(wazi.NewAnnotationsReconciler())
		reconcileManager.RegisterReconciler(wazi.NewODLMReconciler())
		reconcileManager.RegisterReconciler(wazi.NewLicenseReconciler())
	}

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

func (r *WaziLicenseReconciler) SetupWithManager(mgr ctrl.Manager) error {

	controllerBuilder := ctrl.NewControllerManagedBy(mgr).
		Watches(&source.Kind{Type: &chev2.WaziLicense{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &odlmv1alpha1.OperandRequest{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &chev2.WaziLicense{},
		})

	if r.namespace != "" {
		controllerBuilder = controllerBuilder.WithEventFilter(utils.InNamespaceEventFilter(r.namespace))
	}

	return controllerBuilder.
		For(&chev2.WaziLicense{}).
		Complete(r)
}

func (r *WaziLicenseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("wazilicense", req.NamespacedName)

	clusterAPI := chetypes.ClusterAPI{
		Client:           r.client,
		NonCachingClient: r.nonCachedClient,
		DiscoveryClient:  r.discoveryClient,
		Scheme:           r.Scheme,
	}

	wazilicense, _, err := wazi.GetWaziLicenseCRWithNamespacedName(r.client, req.NamespacedName)

	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("WaziLicense Custom Resource not found.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	deployContext := &chetypes.DeployContext{
		ClusterAPI:  clusterAPI,
		WaziLicense: wazilicense,
	}

	if deployContext.WaziLicense.ObjectMeta.DeletionTimestamp.IsZero() {

		result, done, err := r.reconcileManager.ReconcileAll(deployContext)

		if !done {
			return result, err
		} else {
			logrus.Info("Successfully reconciled.")
			return ctrl.Result{}, nil
		}

	} else {
		done := r.reconcileManager.FinalizeAll(deployContext)
		return ctrl.Result{Requeue: !done}, nil
	}
}
