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
package util

import (
	"os"
	"reflect"
	"strings"
	"testing"

	orgv1 "github.com/eclipse-che/che-operator/api/v1"
	"github.com/google/go-cmp/cmp"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	testNamespace    = "wazi-devspaces"
	subscriptionName = "ibm-wazi-developer-for-workspaces"
	waziCsvName      = "wazi-operator-csv"
	waziLicenseName        = "wazi-license"
)

func TestFromIDzEELicenseUsageString(t *testing.T) {
	licenseUsage := FromLicenseUsageString("idzee")
	if !strings.Contains(licenseUsage.String(), "idzee") {
		t.Fatalf("IDzEE license usage did not contain \"idzee\", was: %v", licenseUsage.String())
	}
}
func TestFromWaziLicenseUsageString(t *testing.T) {
	licenseUsage := FromLicenseUsageString("wazi")
	if !strings.Contains(licenseUsage.String(), "wazi") {
		t.Fatalf("Wazi license usage did not contain \"wazi\", was: %v", licenseUsage.String())
	}
}

func TestGetEmptyWaziLicenseUsage(t *testing.T) {
	waziLicense := &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "eclipse-che",
			Name:      "eclipse-che",
		},
	}
	usage := GetWaziLicenseUsage(waziLicense)
	if usage != "wazi" {
		t.Fatalf("Wazi license usage should be \"wazi\", was: %v", usage)
	}
}

func TestAddWaziLicenseEnv(t *testing.T) {
	waziLicense := &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       waziLicenseName,
			Name:            waziLicenseName,
			ResourceVersion: "1",
		},
	}

	orgv1.SchemeBuilder.AddToScheme(scheme.Scheme)
	cli := fake.NewFakeClientWithScheme(scheme.Scheme, waziLicense)
	env := []corev1.EnvVar{}

	sizeBeforeAdding := len(env)
	AddWaziLicenseEnv(cli, waziLicense.Namespace, &env)
	if len(env) < sizeBeforeAdding + 1 {
		t.Fatalf("Did not add wazi license to envvironment")
	}
	found := false
	for _, ele := range env {
		if ele.Value == "wazi" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Could not find wazi license within environment: %v", env)
	}
}

func TestGetFilledWaziLicenseUsage(t *testing.T) {
	waziLicense := &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "eclipse-che",
			Name:      "eclipse-che",
		},
	}
	waziLicense.Spec.License.Use = "wazi"
	usage := GetWaziLicenseUsage(waziLicense)
	if usage != "wazi" {
		t.Fatalf("Wazi license usage should be \"wazi\", was: %v", usage)
	}
}

func TestFindWaziLicenseCRinNamespace(t *testing.T) {
	type testCase struct {
		name              string
		initObjects       []runtime.Object
		watchNamespace    string
		expectedNumber    int
		expectedNamespace string
		expectedErr       bool
	}

	testCases := []testCase{
		{
			name: "CR in 'wazi-devspaces' namespace",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: namespace}},
			},
			watchNamespace:    namespace,
			expectedNumber:    1,
			expectedErr:       false,
			expectedNamespace: namespace,
		},
		{
			name: "CR in 'default' namespace",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: "default"}},
			},
			watchNamespace: namespace,
			expectedNumber: 0,
			expectedErr:    true,
		},
		{
			name: "several CR in 'wazi-devspaces' namespace",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: namespace}},
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: "test-wazi-devspaces", Namespace: namespace}},
			},
			watchNamespace: namespace,
			expectedNumber: 2,
			expectedErr:    true,
		},
		{
			name: "several CR in different namespaces",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: namespace}},
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: "default"}},
			},
			watchNamespace:    namespace,
			expectedNumber:    1,
			expectedErr:       false,
			expectedNamespace: namespace,
		},
		{
			name: "CR in 'wazi-devspaces' namespace, all-namespace mode",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: namespace}},
			},
			watchNamespace:    "",
			expectedNumber:    1,
			expectedErr:       false,
			expectedNamespace: namespace,
		},
		{
			name: "CR in 'default' namespace, all-namespace mode",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: "default"}},
			},
			watchNamespace:    "",
			expectedNumber:    1,
			expectedErr:       false,
			expectedNamespace: "default",
		},
		{
			name: "several CR in 'wazi-devspaces' namespace, all-namespace mode",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: namespace}},
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: "test-wazi-devspaces", Namespace: namespace}},
			},
			watchNamespace: "",
			expectedNumber: 2,
			expectedErr:    true,
		},
		{
			name: "several CR in different namespaces, all-namespace mode",
			initObjects: []runtime.Object{
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: namespace}},
				&orgv1.WaziLicense{ObjectMeta: metav1.ObjectMeta{Name: namespace, Namespace: "default"}},
			},
			watchNamespace: "",
			expectedNumber: 2,
			expectedErr:    true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			scheme := scheme.Scheme
			orgv1.SchemeBuilder.AddToScheme(scheme)
			cli := fake.NewFakeClientWithScheme(scheme, testCase.initObjects...)

			WaziLicense, num, err := FindWaziLicenseCRInNamespace(cli, testCase.watchNamespace)
			assert.Equal(t, testCase.expectedNumber, num)
			if testCase.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			if num == 1 {
				assert.Equal(t, testCase.expectedNamespace, WaziLicense.Namespace)
			}
		})
	}
}

func TestFindWaziPodInNamespace(t *testing.T) {
	type testCase struct {
		name           string
		initObjects    []runtime.Object
		watchNamespace string
		expectedPod    *corev1.Pod
	}

	testCases := []testCase{
		{
			name: "Extract pod from specific namespace with match in label selector",
			initObjects: []runtime.Object{
				&corev1.PodList{
					TypeMeta: metav1.TypeMeta{Kind: "PodList", APIVersion: "v1"},
					Items: []corev1.Pod{
						*CreateTestPod("wazi-license-operator", namespace, map[string]string{"app": "wazi-license-operator"}),
						*CreateTestPod("wazi-license-operator", "sample-namespace", map[string]string{"app": "wazi-license-operator"}),
					},
				},
			},
			watchNamespace: namespace,
			expectedPod:    CreateTestPod("wazi-license-operator", namespace, map[string]string{"app": "wazi-license-operator"}),
		},
		{
			name: "Extract pod in same namespace using match in specific label selector",
			initObjects: []runtime.Object{
				&corev1.PodList{
					TypeMeta: metav1.TypeMeta{Kind: "PodList", APIVersion: "v1"},
					Items: []corev1.Pod{
						*CreateTestPod("wazi-license-operator-1", namespace, map[string]string{"app": "wazi-license-operator"}),
						*CreateTestPod("wazi-license-operator-2", namespace, map[string]string{"app": "sample-label"}),
					},
				},
			},
			watchNamespace: namespace,
			expectedPod:    CreateTestPod("wazi-license-operator-1", namespace, map[string]string{"app": "wazi-license-operator"}),
		},
		{
			name: "No match returned with mismatched label selector in required namespace",
			initObjects: []runtime.Object{
				&corev1.PodList{
					TypeMeta: metav1.TypeMeta{Kind: "PodList", APIVersion: "v1"},
					Items: []corev1.Pod{
						*CreateTestPod("wazi-license-operator-1", namespace, map[string]string{"app": "sample-label-1"}),
						*CreateTestPod("wazi-license-operator-2", namespace, map[string]string{"app": "sample-label-2"}),
					},
				},
			},
			watchNamespace: namespace,
			expectedPod:    nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			scheme := scheme.Scheme
			orgv1.SchemeBuilder.AddToScheme(scheme)
			corev1.AddToScheme(scheme)
			cli := fake.NewFakeClientWithScheme(scheme, testCase.initObjects...)
			os.Setenv("CHE_FLAVOR", waziLicenseName)

			testPod, _, _ := FindWaziPodInNamespace(cli, testCase.watchNamespace)
			if !reflect.DeepEqual(testCase.expectedPod, testPod) {
				t.Errorf("Expected Pod and Pod returned from API server are different (-want +got): %v", cmp.Diff(testCase.expectedPod, testPod))
			}

		})
	}
}

func TestFindWaziSubscriptionInNamespace(t *testing.T) {
	type testCase struct {
		name           string
		initObjects    []runtime.Object
		watchNamespace string
		expectedCR     *operatorsv1alpha1.Subscription
	}

	testCases := []testCase{
		{
			name: "Test Client to check if Wazi subscription present in desired namespace",
			initObjects: []runtime.Object{
				&operatorsv1alpha1.Subscription{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
						Name:      subscriptionName,
					},
				},
			},
			watchNamespace: namespace,
			expectedCR: &operatorsv1alpha1.Subscription{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Subscription",
					APIVersion: "operators.coreos.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       namespace,
					Name:            subscriptionName,
					ResourceVersion: "999",
				},
			},
		},
		{
			name: "Test Client to check if Wazi subsciption is absent in desired namespace",
			initObjects: []runtime.Object{
				&operatorsv1alpha1.Subscription{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "random-namespace",
						Name:      subscriptionName,
					},
				},
			},
			watchNamespace: namespace,
			expectedCR:     nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			scheme := scheme.Scheme
			orgv1.SchemeBuilder.AddToScheme(scheme)
			operatorsv1alpha1.SchemeBuilder.AddToScheme(scheme)
			cli := fake.NewFakeClientWithScheme(scheme, testCase.initObjects...)
			Subscription, _ := FindWaziSubscriptionInNamespace(cli, testCase.watchNamespace)
			if !reflect.DeepEqual(testCase.expectedCR, Subscription) {
				t.Errorf("Expected Subscription and Subscription returned from API server are different (-want +got): %v", cmp.Diff(testCase.expectedCR, Subscription))
			}

		})
	}
}

func TestFindWaziClusterServiceVersionInNamespace(t *testing.T) {
	type testCase struct {
		name           string
		initObjects    []runtime.Object
		watchNamespace string
		watchName      string
		expectedCR     *operatorsv1alpha1.ClusterServiceVersion
	}

	testCases := []testCase{
		{
			name: "Test Client to check if Wazi clusterServiceVersion is present in desired namespace",
			initObjects: []runtime.Object{
				&operatorsv1alpha1.ClusterServiceVersion{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
						Name:      waziCsvName,
					},
				},
			},
			watchNamespace: namespace,
			watchName:      waziCsvName,
			expectedCR: &operatorsv1alpha1.ClusterServiceVersion{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterServiceVersion",
					APIVersion: "operators.coreos.com/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace:       namespace,
					Name:            waziCsvName,
					ResourceVersion: "999",
				},
			},
		},
		{
			name: "Test Client to check if Wazi ClusterServiceVersion is absent in desired namespace",
			initObjects: []runtime.Object{
				&operatorsv1alpha1.Subscription{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "random-namespace",
						Name:      waziCsvName,
					},
				},
			},
			watchNamespace: namespace,
			watchName:      waziCsvName,
			expectedCR:     nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			scheme := scheme.Scheme
			orgv1.SchemeBuilder.AddToScheme(scheme)
			operatorsv1alpha1.SchemeBuilder.AddToScheme(scheme)
			cli := fake.NewFakeClientWithScheme(scheme, testCase.initObjects...)
			testWaziCSV, _ := FindWaziClusterServiceVersionInNamespace(cli, testCase.watchNamespace, testCase.watchName)
			if !reflect.DeepEqual(testCase.expectedCR, testWaziCSV) {
				t.Errorf("Expected CSV and CSV returned from API server are different (-want +got): %v", cmp.Diff(testCase.expectedCR, testWaziCSV))
			}

		})
	}
}

// CreateTestPod is a helper function to create the Pod Spec for test resources.
func CreateTestPod(name, namespace string, labels map[string]string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       namespace,
			Labels:          labels,
			ResourceVersion: "999",
		},
	}
}
