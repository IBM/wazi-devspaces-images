// //
// // Copyright (c) 2022 IBM Corporation
// // This program and the accompanying materials are made
// // available under the terms of the Eclipse Public License 2.0
// // which is available at https://www.eclipse.org/legal/epl-2.0/
// //
// // SPDX-License-Identifier: EPL-2.0
// //
// // Contributors:
// //   IBM Corporation - initial API and implementation
// //

package wazi

import (
	"context"
	"reflect"
	"strings"

	chev2 "github.com/eclipse-che/che-operator/api/v2"
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	test "github.com/eclipse-che/che-operator/pkg/common/test"
	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	"testing"
)

func TestReloadWaziLicense(t *testing.T) {

	objs := []runtime.Object{}
	ctx := test.GetWaziDeployContext(nil, objs)
	cli := ctx.ClusterAPI.Client

	waziLicense := &chev2.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: test.TestWaziLicenseNamespace,
			Name:      test.TestWaziLicenseName,
		},
	}

	deployContext := &chetypes.DeployContext{
		WaziLicense: waziLicense,
		ClusterAPI: chetypes.ClusterAPI{
			Client:           cli,
			NonCachingClient: cli,
			Scheme:           scheme.Scheme,
		},
	}

	err := ReloadWaziLicenseCR(deployContext)
	if err != nil {
		t.Errorf("Failed to reload wazilicense, %v", err)
	}
}

func TestGetWaziLicenseCR(t *testing.T) {

	ctx := test.GetWaziDeployContext(nil, []runtime.Object{})
	_, _, err := GetWaziLicenseCR(ctx)
	if err != nil {
		t.Errorf("Failed to get wazilicense, %v", err)
	}
}

func TestSyncWaziObject(t *testing.T) {

}

func TestAppendFinalizer(t *testing.T) {

	objs := []runtime.Object{}
	ctx := test.GetWaziDeployContext(nil, objs)
	cl := ctx.ClusterAPI.Client

	expectedCR := chev2.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: test.TestWaziLicenseNamespace,
			Name:      test.TestWaziLicenseName,
			Finalizers: []string{
				"sample-finalizer",
			},
		},
	}

	err := AppendFinalizer(ctx, "sample-finalizer")
	if err != nil {
		t.Errorf("Failed to reconcile operand finalizer, %v", err)
	}
	waziLicense := chev2.WaziLicense{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: test.TestWaziLicenseName, Namespace: test.TestWaziLicenseNamespace}, &waziLicense)
	if err != nil {
		t.Errorf("Failed to fetch the given Wazi-license spec")
	}
	if !reflect.DeepEqual(expectedCR.Finalizers, waziLicense.Finalizers) {
		t.Errorf("Expected Annotations and Annotations returned from API server for the given PodSpec are different (-want +got): %v", cmp.Diff(waziLicense, expectedCR))
	}
}

func TestDeleteFinalizer(t *testing.T) {

	objs := []runtime.Object{}
	ctx := test.GetWaziDeployContext(nil, objs)
	cl := ctx.ClusterAPI.Client

	expectedCR := chev2.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  test.TestWaziLicenseNamespace,
			Name:       test.TestWaziLicenseName,
			Finalizers: nil,
		},
	}
	ctx.WaziLicense.Finalizers = append(ctx.WaziLicense.Finalizers, "sample-finalizer")
	err := DeleteFinalizer(ctx, "sample-finalizer")
	if err != nil {
		t.Errorf("Failed to reconcile operand finalizer, %v", err)
	}
	waziLicense := chev2.WaziLicense{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: test.TestWaziLicenseName, Namespace: test.TestWaziLicenseNamespace}, &waziLicense)
	if err != nil {
		t.Errorf("Failed to fetch the given Wazi-license spec")
	}
	if !reflect.DeepEqual(expectedCR.Finalizers, waziLicense.Finalizers) {
		t.Errorf("Expected Annotations and Annotations returned from API server for the given PodSpec are different (-want +got): %v", cmp.Diff(waziLicense, expectedCR))
	}
}

func TestGetWaziAnnotations(t *testing.T) {
	type testCase struct {
		name                string
		expectedCr          chev2.WaziLicense
		expectedAnnotations map[string]string
	}

	testCases := []testCase{
		{
			name: "Check Annotations when the waziLicense is of use `Wazi` ",
			expectedCr: chev2.WaziLicense{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       test.TestWaziLicenseNamespace,
					Name:            test.TestWaziLicenseName,
					ResourceVersion: "1",
					Annotations:     WaziAnnotations,
				},
				Spec: chev2.WaziLicenseSpec{
					License: chev2.WaziLicenseSpecLicense{
						Accept: true,
						Use:    test.TestWaziLicenseUsageWazi,
					},
				},
			},
			expectedAnnotations: WaziAnnotations,
		},
		{
			name: "Check Annotations when the waziLicense is of use `IDzEE` ",
			expectedCr: chev2.WaziLicense{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       test.TestWaziLicenseNamespace,
					Name:            test.TestWaziLicenseName,
					ResourceVersion: "1",
					Annotations:     WaziAnnotations,
				},
				Spec: chev2.WaziLicenseSpec{
					License: chev2.WaziLicenseSpecLicense{
						Accept: true,
						Use:    test.TestWaziLicenseUsageIDzEE,
					},
				},
			},
			expectedAnnotations: IDzEEAnnotations,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			objs := []runtime.Object{}
			ctx := test.GetWaziDeployContext(nil, objs)

			ctx.WaziLicense = &testCase.expectedCr
			annotations := GetWaziAnnotations(ctx)
			if !reflect.DeepEqual(testCase.expectedAnnotations, annotations) {
				t.Errorf("Expected Annotations and Annotations returned from API server for the given PodSpec are different (-want +got): %v", cmp.Diff(testCase.expectedAnnotations, annotations))

			}
		})
	}
}

func TestAddWaziLicenseUsageEnvVar(t *testing.T) {
	type testCase struct {
		name             string
		expectedEnvVar   []v1.EnvVar
		inputWaziLicense chev2.WaziLicense
		ExpectedValue    string
	}

	testCases := []testCase{
		{
			name:           "Test for Wazi License Usage",
			expectedEnvVar: []v1.EnvVar{},
			inputWaziLicense: chev2.WaziLicense{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: test.TestWaziLicenseNamespace,
					Name:      test.TestWaziLicenseName,
				},
				Spec: chev2.WaziLicenseSpec{
					License: chev2.WaziLicenseSpecLicense{
						Accept: true,
						Use:    test.TestWaziLicenseUsageWazi,
					},
				},
			},
			ExpectedValue: strings.ToLower(test.TestWaziLicenseUsageWazi),
		},
		{
			name: "Test for IDzEE License Usage",
			expectedEnvVar: []v1.EnvVar{
				{Name: envVarWaziLicenseUsage,
					Value: test.TestWaziLicenseUsageIDzEE,
				},
			},
			inputWaziLicense: chev2.WaziLicense{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: test.TestWaziLicenseNamespace,
					Name:      test.TestWaziLicenseName,
				},
				Spec: chev2.WaziLicenseSpec{
					License: chev2.WaziLicenseSpecLicense{
						Accept: true,
						Use:    test.TestWaziLicenseUsageIDzEE,
					},
				},
			},
			ExpectedValue: test.TestWaziLicenseUsageIDzEE,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			objs := []runtime.Object{}
			ctx := test.GetWaziDeployContext(nil, objs)

			ctx.WaziLicense = &testCase.inputWaziLicense
			_, err := AddWaziLicenseUsageEnvVar(ctx, test.TestWaziLicenseNamespace, &testCase.expectedEnvVar)

			if err != nil {
				t.Errorf("AddWaziLicenseUsageEnvVar failed with the following error %v", err)
			}

			for _, expectedEnvVar := range testCase.expectedEnvVar {

				if testCase.ExpectedValue != expectedEnvVar.Value {
					t.Errorf("Expected value was %s and got %s instead", testCase.ExpectedValue, expectedEnvVar.Value)
				} else {
					break
				}
			}
		})
	}
}
