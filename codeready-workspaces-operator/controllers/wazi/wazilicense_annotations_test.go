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
	"os"
	"reflect"
	"testing"

	licensingoperatorv1 "github.com/IBM/ibm-licensing-operator/api/v1"
	"github.com/blang/semver/v4"
	orgv1 "github.com/eclipse-che/che-operator/api/v1"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/google/go-cmp/cmp"
	"github.com/operator-framework/api/pkg/lib/version"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestReconcileWaziAnnotations(t *testing.T) {
	type testCase struct {
		name       string
		expectedCr *corev1.Pod
		crUse      string
	}

	testCases := []testCase{
		{
			name:       "Check Annotations when the waziLicense is of use `Wazi` ",
			crUse:      "wazi",
			expectedCr: returnTestPodsWithAnnotations(util.WaziAnnotations),
		},
		{
			name:       "Check Anootations when the waziLicense if of use Upper Case `Wazi`",
			crUse:      "WAZI",
			expectedCr: returnTestPodsWithAnnotations(util.WaziAnnotations),
		},
		{
			name:       "Check Annotations when the waziLicense is of use `Idzee` ",
			crUse:      "idzee",
			expectedCr: returnTestPodsWithAnnotations(util.IDzEEAnnotations),
		},
		{
			name:       "Check Annotations when the waziLicense is of use `Wazi` with varied upper and lower case variables ",
			crUse:      "IDZee",
			expectedCr: returnTestPodsWithAnnotations(util.IDzEEAnnotations),
		},
	}
	for _, testCase := range testCases {
		cl, dc, scheme := Init()
		r := NewReconciler(cl, cl, dc, &scheme, "")
		deployContext := deploy.DeployContext{
			WaziLicense: initWaziorIdzeeWithSimpleCR(testCase.crUse),
			ClusterAPI: deploy.ClusterAPI{
				Client:           r.client,
				NonCachingClient: r.client,
				Scheme:           r.Scheme,
			},
		}
		os.Setenv("CHE_FLAVOR", waziLicenseName)
		done, err := r.ReconcileWaziAnnotations(&deployContext)
		if !done {
			t.Errorf("ReconcileWaziAnnotations failed with err: %v", err)
		}
		os.Unsetenv("CHE_FLAVOR")
		annotatedPod := &corev1.Pod{}
		if err := cl.Get(context.TODO(), types.NamespacedName{Name: "fake-pg-pod", Namespace: namespace}, annotatedPod); err != nil {
			t.Fatalf("Could not get the given Annotated Pod: err=%v", err)
		}

		if !reflect.DeepEqual(testCase.expectedCr.Annotations, annotatedPod.Annotations) {
			t.Errorf("Expected Annotations and Annotations returned from API server for the given PodSpec are different (-want +got): %v", cmp.Diff(testCase.expectedCr, annotatedPod))
		}

	}
}

func returnTestPodsWithAnnotations(annotations map[string]string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-pg-pod",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "wazi-license-operator",
			},
			Annotations: annotations,
		},
	}
}

func initWaziWithIBMLicensingQuerySource(name string) *licensingoperatorv1.IBMLicensingQuerySource {
	return &licensingoperatorv1.IBMLicensingQuerySource{
		TypeMeta: metav1.TypeMeta{
			Kind:       "operator.ibm.com",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

}

func initWaziCSVWithVersion() *operatorsv1alpha1.ClusterServiceVersion {
	return &operatorsv1alpha1.ClusterServiceVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      waziLicenseName + "-csv",
			Namespace: namespace,
		},
		Spec: operatorsv1alpha1.ClusterServiceVersionSpec{
			Version: version.OperatorVersion{
				Version: semver.SpecVersion,
			},
		},
	}

}
func initWaziorIdzeeWithSimpleCR(licenseUsage string) *orgv1.WaziLicense {
	return &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Name:      waziLicenseName,
			Namespace: namespace,
		},
		Spec: orgv1.WaziLicenseSpec{
			License: orgv1.WaziLicenseSpecLicense{
				Accept: true,
				Use:    licenseUsage,
			},
		},
	}
}

func initWaziSubscription(name string) *operatorsv1alpha1.Subscription {
	return &operatorsv1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: operatorsv1alpha1.SubscriptionStatus{
			InstalledCSV: waziLicenseName + "-csv",
		},
	}
}
