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
	"testing"

	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	testOperandRequestName = "ibm-wazi-code-operand-request"
	sampleSubscriptionName = "sample-subscription"
)

func TestReconcileWaziOperandRequest(t *testing.T) {
	type testCase struct {
		name                string
		operandExists       bool
		operandStatusExists bool
		expectedOutcome     bool
	}

	testCases := []testCase{
		{
			name:                "Check Reconcile Condition: creation of OperandRequest if absent",
			operandExists:       false,
			operandStatusExists: false,
			expectedOutcome:     false,
		},
		{
			name:                "Check Reconcile Condition: Member status absent and OperandRequest present ",
			operandExists:       true,
			operandStatusExists: false,
			expectedOutcome:     false,
		},
		{
			name:                "Check Reconcile Condition: Member status and OperandRequest absent ",
			operandExists:       true,
			operandStatusExists: true,
			expectedOutcome:     true,
		},
	}
	for _, testCase := range testCases {
		objs, dc, scheme := createAPIObjects()
		if expectedCR := returnExpectedCR(testCase.operandExists, testCase.operandStatusExists); expectedCR != nil {
			objs = append(objs, returnExpectedCR(testCase.operandExists, testCase.operandStatusExists))
		}
		cl := fake.NewFakeClient(objs...)
		r := NewReconciler(cl, cl, dc, &scheme, "")
		deployContext := deploy.DeployContext{
			WaziLicense: InitWaziWithSimpleCR(),
			ClusterAPI: deploy.ClusterAPI{
				Client:           r.client,
				NonCachingClient: r.client,
				Scheme:           r.Scheme,
			},
		}
		done, err := r.ReconcileWaziOperandRequest(&deployContext)
		if err != nil {
			t.Fatalf("ReconcileWaziOperandRequest failed: err= %v", err)
		}
		if done != testCase.expectedOutcome {
			t.Fatalf("Expected outcome: %v, Recieved outcome: %v", testCase.expectedOutcome, done)
		}
		if !testCase.operandExists {
			testODLMCr := odlmv1alpha1.OperandRequest{}
			if err := cl.Get(context.TODO(), types.NamespacedName{Name: testOperandRequestName, Namespace: namespace}, &testODLMCr); err != nil {
				t.Fatalf("The test case %v failed with err %v", testCase.name, err)
			}
		}
	}
}
