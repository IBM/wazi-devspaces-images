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

	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *WaziLicenseReconciler) ReconcileWaziOperandRequest(deployContext *deploy.DeployContext) (done bool, err error) {

	operandRequestSpec := getWaziOperandRequestSpec(deployContext, util.WaziOperandRequestName)
	_, err = deploy.WaziCreateIfNotExists(deployContext, operandRequestSpec)
	if err != nil { // Issue creating resource
		return false, err
	}

	// Get the Operand Request for status data
	operandRequest, err := getWaziOperandRequest(deployContext, util.WaziOperandRequestName)
	if err != nil { // No Operand Request found
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	if len(operandRequest.Status.Members) == 0 {
		return false, nil // No members to iterate over
	}

	for _, member := range operandRequest.Status.Members {
		// Only care about the Operator to be running, not the instance
		if member.Phase.OperatorPhase != odlmv1alpha1.OperatorRunning {
			return false, nil
		}
	}

	return true, nil
}

func getWaziOperandRequestSpec(deployContext *deploy.DeployContext, name string) *odlmv1alpha1.OperandRequest {

	operands := []odlmv1alpha1.Operand{}
	for _, commonService := range util.WaziCommonServices {
		operands = append(operands, odlmv1alpha1.Operand{Name: commonService})
	}

	operandRequest := &odlmv1alpha1.OperandRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       util.WaziOperandRequestKind,
			APIVersion: util.WaziOperandRequestAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: deployContext.WaziLicense.Namespace,
		},
		Spec: odlmv1alpha1.OperandRequestSpec{
			Requests: []odlmv1alpha1.Request{
				{
					Registry:          util.WaziOperandRequestRegistry,
					RegistryNamespace: util.WaziOperandRequestRegistryNamespace,
					Operands:          operands,
				},
			},
		},
	}

	return operandRequest
}

func getWaziOperandRequest(deployContext *deploy.DeployContext, name string) (waziOperandRequest *odlmv1alpha1.OperandRequest, err error) {
	waziOperandRequest = &odlmv1alpha1.OperandRequest{}
	err = deployContext.ClusterAPI.Client.Get(context.TODO(), types.NamespacedName{Namespace: deployContext.WaziLicense.Namespace, Name: name}, waziOperandRequest)
	if err != nil {
		return nil, err
	}
	return waziOperandRequest, nil
}
