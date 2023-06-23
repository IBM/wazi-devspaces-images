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
	"reflect"
	"testing"

	licensequerysrc "github.com/IBM/ibm-licensing-operator/api/v1"
	test "github.com/eclipse-che/che-operator/pkg/common/test"
	"github.com/eclipse-che/che-operator/pkg/common/utils"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var (
	licenseReconciler = NewLicenseReconciler()
)

func TestLicenseReconcile(t *testing.T) {
	type testCase struct {
		name                     string
		WaziLicenseType          string
		expectedLicenseQuerySpec licensequerysrc.IBMLicensingQuerySource
	}

	testCases := []testCase{
		{
			name:            "Usage: Wazi",
			WaziLicenseType: test.TestWaziLicenseUsageWazi,
			expectedLicenseQuerySpec: licensequerysrc.IBMLicensingQuerySource{
				TypeMeta: metav1.TypeMeta{
					Kind:       waziLicensingKind,
					APIVersion: waziLicensingAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            waziLicensingName,
					Namespace:       test.TestWaziLicenseNamespace,
					ResourceVersion: "1",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "org.eclipse.che/v2",
							Kind:               test.TestWaziLicenseKind,
							Name:               test.TestWaziLicenseName,
							Controller:         test.TrueWaziPointer(),
							BlockOwnerDeletion: test.TrueWaziPointer(),
						},
					},
				},
				Spec: licensequerysrc.IBMLicensingQuerySourceSpec{
					Query:       waziLicensinQuery,
					Annotations: WaziAnnotations,
				},
			},
		},
		{
			name:            "Usage: Idzee",
			WaziLicenseType: test.TestWaziLicenseUsageIDzEE,
			expectedLicenseQuerySpec: licensequerysrc.IBMLicensingQuerySource{
				TypeMeta: metav1.TypeMeta{
					Kind:       waziLicensingKind,
					APIVersion: waziLicensingAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            waziLicensingName,
					Namespace:       test.TestWaziLicenseNamespace,
					ResourceVersion: "1",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         "org.eclipse.che/v2",
							Kind:               test.TestWaziLicenseKind,
							Name:               test.TestWaziLicenseName,
							Controller:         test.TrueWaziPointer(),
							BlockOwnerDeletion: test.TrueWaziPointer(),
						},
					},
				},
				Spec: licensequerysrc.IBMLicensingQuerySourceSpec{
					Query:       waziLicensinQuery,
					Annotations: IDzEEAnnotations,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			waziLicense := test.GetFakeWaziLicense(true, testCase.WaziLicenseType)
			objs := []runtime.Object{}
			ctx := test.GetWaziDeployContext(waziLicense, objs)

			licenseQuerySource := licenseReconciler.getLicensingSpec(ctx)
			test.SetLicenseQuerySourceOwnerReference(licenseQuerySource)

			objs = []runtime.Object{licenseQuerySource}
			ctx = test.GetWaziDeployContext(waziLicense, objs)

			_, _, err := licenseReconciler.Reconcile(ctx)
			if err != nil {
				t.Error("Reconciliation loop failed")
			}

			actualLicenseQuerySource := licensequerysrc.IBMLicensingQuerySource{}
			err = ctx.ClusterAPI.Client.Get(context.TODO(), types.NamespacedName{Namespace: test.TestWaziLicenseNamespace, Name: waziLicensingName}, &actualLicenseQuerySource)

			if err != nil {
				t.Error("Failed to fetch the IBM licensing query source")
			}

			if !reflect.DeepEqual(testCase.expectedLicenseQuerySpec, actualLicenseQuerySource) {
				t.Errorf("Expected IBM licensing query source vs generated IBM licensing query source (-want +got): %v", cmp.Diff(testCase.expectedLicenseQuerySpec, actualLicenseQuerySource))
			}
		})
	}

}

func TestLicenseFinalize(t *testing.T) {

	objs := []runtime.Object{}
	ctx := test.GetWaziDeployContext(nil, objs)
	licenseQuerySource := licenseReconciler.getLicensingSpec(ctx)
	objs = []runtime.Object{licenseQuerySource}
	ctx = test.GetWaziDeployContext(nil, objs)

	ctx.WaziLicense.Finalizers = append(ctx.WaziLicense.Finalizers, deploy.GetFinalizerName(waziLicensingResourceName))
	done := licenseReconciler.Finalize(ctx)

	if !done {
		t.Error("ODLM Finalize fucntion failed")
	}

	if utils.Contains(ctx.WaziLicense.Finalizers, deploy.GetFinalizerName(waziLicensingResourceName)) {
		t.Error("Finalizer failed to delete finalizer type: waziLicensingResourceName")
	}
}
