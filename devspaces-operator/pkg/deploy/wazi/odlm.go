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

	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var waziCommonServices = []string{"ibm-licensing-operator"}

const (
	waziOperandRequestName              = "ibm-wazi-code-operand-request"
	waziOperandRequestKind              = "OperandRequest"
	waziOperandRequestAPIVersion        = "operator.ibm.com/v1alpha1"
	waziOperandRequestRegistry          = "common-service"
	waziOperandRequestRegistryNamespace = "ibm-common-services"
)

type ODLMReconciler struct {
	deploy.Reconcilable
}

func NewODLMReconciler() *ODLMReconciler {
	return &ODLMReconciler{}
}

func (d *ODLMReconciler) Reconcile(ctx *chetypes.DeployContext) (reconcile.Result, bool, error) {

	operandRequestSpec := d.getOperandRequestSpec(ctx)

	if _, err := WaziCreateIfNotExists(ctx, operandRequestSpec); err != nil {
		return reconcile.Result{}, false, err
	}

	operandRequest, _, err := d.getOperandRequest(ctx)
	if err != nil {
		return reconcile.Result{}, false, err
	}

	done := d.isOperandRequestRunning(operandRequest)
	if !done {
		return reconcile.Result{}, false, nil
	}

	err = AppendFinalizer(ctx, deploy.GetFinalizerName("operandrequests"))
	if err != nil {
		return reconcile.Result{}, false, err
	}

	return reconcile.Result{}, true, nil
}

func (d *ODLMReconciler) Finalize(ctx *chetypes.DeployContext) bool {
	done := true
	if err := DeleteFinalizer(ctx, deploy.GetFinalizerName("operandrequests")); err != nil {
		done = false
		logrus.Errorf("Error deleting finalizer: %v", err)
	}
	return done
}

func (d *ODLMReconciler) getOperandRequestSpec(ctx *chetypes.DeployContext) *odlmv1alpha1.OperandRequest {

	ns := ctx.WaziLicense.Namespace
	operands := []odlmv1alpha1.Operand{}

	for _, commonService := range waziCommonServices {
		operands = append(operands, odlmv1alpha1.Operand{Name: commonService})
	}

	operandRequest := &odlmv1alpha1.OperandRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       waziOperandRequestKind,
			APIVersion: waziOperandRequestAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      waziOperandRequestName,
			Namespace: ns,
		},
		Spec: odlmv1alpha1.OperandRequestSpec{
			Requests: []odlmv1alpha1.Request{
				{
					Registry:          waziOperandRequestRegistry,
					RegistryNamespace: waziOperandRequestRegistryNamespace,
					Operands:          operands,
				},
			},
		},
	}

	return operandRequest
}

func (d *ODLMReconciler) getOperandRequest(ctx *chetypes.DeployContext) (*odlmv1alpha1.OperandRequest, bool, error) {

	cl := ctx.ClusterAPI.Client
	operandRequest := &odlmv1alpha1.OperandRequest{}
	namespacedName := types.NamespacedName{
		Name:      waziOperandRequestName,
		Namespace: ctx.WaziLicense.Namespace,
	}

	err := cl.Get(context.TODO(),
		namespacedName,
		operandRequest)
	if err != nil {
		return nil, false, err
	}

	return operandRequest, true, nil
}

func (d *ODLMReconciler) isOperandRequestRunning(operandRequest *odlmv1alpha1.OperandRequest) bool {
	for _, member := range operandRequest.Status.Members {
		if member.Phase.OperatorPhase == odlmv1alpha1.OperatorRunning {
			return true
		}
	}
	return false
}
