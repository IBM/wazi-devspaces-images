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
	"testing"

	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	"github.com/eclipse-che/che-operator/pkg/common/test"
	"github.com/eclipse-che/che-operator/pkg/common/utils"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var odlmReconciler = NewODLMReconciler()

func TestIsOperandRequestComplete(t *testing.T) {
	var tests = []struct {
		phase odlmv1alpha1.OperatorPhase
		want  bool
	}{
		{odlmv1alpha1.OperatorRunning, true},
		{odlmv1alpha1.OperatorReady, false},
		{odlmv1alpha1.OperatorInstalling, false},
		{odlmv1alpha1.OperatorUpdating, false},
		{odlmv1alpha1.OperatorFailed, false},
		{odlmv1alpha1.OperatorInit, false},
		{odlmv1alpha1.OperatorNotFound, false},
		{odlmv1alpha1.OperatorNone, false},
	}

	for _, tt := range tests {
		testname := string(tt.phase)
		var operandRequest *odlmv1alpha1.OperandRequest = &odlmv1alpha1.OperandRequest{
			Status: odlmv1alpha1.OperandRequestStatus{
				Members: []odlmv1alpha1.MemberStatus{{
					Phase: odlmv1alpha1.MemberPhase{
						OperatorPhase: tt.phase}}},
			},
		}
		t.Run(testname, func(t *testing.T) {
			result := odlmReconciler.isOperandRequestRunning(operandRequest)
			assert.Equal(t, result, tt.want)
		})
	}

}

func TestGetOperandSpec(t *testing.T) {

	ctx := test.GetWaziDeployContext(nil, []runtime.Object{})
	if operandRequest := odlmReconciler.getOperandRequestSpec(ctx); operandRequest == nil {
		t.Error("getOperandSpec returned unexpected nil result")
	}
}

func TestODLMFinalize(t *testing.T) {

	ctx := test.GetWaziDeployContext(nil, []runtime.Object{})
	ctx.WaziLicense.Finalizers = append(ctx.WaziLicense.Finalizers, deploy.GetFinalizerName("operandrequests"))
	done := odlmReconciler.Finalize(ctx)
	if !done {
		t.Error("ODLM Finalize fucntion failed")
	}

	if utils.Contains(ctx.WaziLicense.Finalizers, deploy.GetFinalizerName("operandrequests")) {
		t.Error("Finalizer failed to delete finalizer type: operandrequests")
	}
}

func TestReconcileWaziOperandRequest(t *testing.T) {
	type testCase struct {
		name                string
		operandExists       bool
		operandStatusExists bool
		checkFinalizer      bool
		expectedOutcome     bool
	}
	testCases := []testCase{
		{
			name:                "Check Reconcile Condition: Operand already exists with status set to running",
			operandExists:       true,
			operandStatusExists: true,
			checkFinalizer:      false,
			expectedOutcome:     true,
		},
		{
			name:                "Check Reconcile Condition: Check Operand creation with reconcile method",
			operandExists:       false,
			operandStatusExists: false,
			checkFinalizer:      false,
			expectedOutcome:     false,
		},
		{
			name:                "Check Reconcile Condition: Check for finalizer append to Wazi License",
			operandExists:       true,
			operandStatusExists: true,
			checkFinalizer:      true,
			expectedOutcome:     true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			objs := []runtime.Object{}
			if operandRequest := test.GetOperandRequest(testCase.operandExists, testCase.operandStatusExists); operandRequest != nil {
				objs = []runtime.Object{operandRequest}
			}

			ctx := test.GetWaziDeployContext(nil, objs)

			_, done, err := odlmReconciler.Reconcile(ctx)
			if err != nil {
				t.Fatalf("ReconcileWaziOperandRequest failed: err= %v", err)
			}
			if done != testCase.expectedOutcome {
				t.Fatalf("Expected outcome: %v, Received outcome: %v", testCase.expectedOutcome, done)
			}
			if !testCase.operandExists {
				testODLMCr := odlmv1alpha1.OperandRequest{}
				if err := ctx.ClusterAPI.Client.Get(context.TODO(), types.NamespacedName{Name: waziOperandRequestName, Namespace: "eclipse-che"}, &testODLMCr); err != nil {
					t.Fatalf("The test case %v failed with err %v", testCase.name, err)
				}
			}
			if testCase.checkFinalizer {
				finalizers := ctx.WaziLicense.GetFinalizers()
				present := checkFinalizerExists(finalizers, deploy.GetFinalizerName("operandrequests"))
				if !present {
					t.Fatalf("Unable to find the required finalizer in the Wazi License")
				}
			}
		})
	}
}

func checkFinalizerExists(elements []string, key string) bool {
	for _, element := range elements {
		if element == key {
			return true
		}
	}
	return false
}
