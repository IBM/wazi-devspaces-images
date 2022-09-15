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
	"strconv"

	"reflect"

	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	licensingoperatorv1 "github.com/IBM/ibm-licensing-operator/api/v1"
	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	chev1alpha1 "github.com/che-incubator/kubernetes-image-puller-operator/api/v1alpha1"
	orgv1 "github.com/eclipse-che/che-operator/api/v1"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/google/go-cmp/cmp"
	operatorsv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	packagesv1 "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"

	"github.com/sirupsen/logrus"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	fakeDiscovery "k8s.io/client-go/discovery/fake"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"testing"
)

var (
	waziLicenseName = "wazi-license"
	namespace = "eclipse-che"
)

func TestNativeUserModeEnabled(t *testing.T) {
	type testCase struct {
		name                    string
		initObjects             []runtime.Object
		cheFlavor               string
		isLicenseAccept         bool
		expectedNativeUserValue bool
	}

	testCases := []testCase{
		{
			name:                    "Wazi License should be accepted from initial CR",
			cheFlavor:               waziLicenseName,
			isLicenseAccept:         true,
			expectedNativeUserValue: true,
		},
		{
			name:                    "Wazi License should not be accepted from initial CR",
			cheFlavor:               waziLicenseName,
			isLicenseAccept:         false,
			expectedNativeUserValue: false,
		},
		{
			name:                    "Wazi License should be accepted from initial CR",
			cheFlavor:               "",
			isLicenseAccept:         true,
			expectedNativeUserValue: true,
		},
		{
			name:                    "Wazi License should not be accepted from initial CR",
			cheFlavor:               "",
			isLicenseAccept:         false,
			expectedNativeUserValue: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logf.SetLogger(zap.New(zap.WriteTo(os.Stdout), zap.UseDevMode(true)))

			scheme := scheme.Scheme
			orgv1.SchemeBuilder.AddToScheme(scheme)
			scheme.AddKnownTypes(crdv1.SchemeGroupVersion, &crdv1.CustomResourceDefinition{})

			initCR := InitWaziWithSimpleCR().DeepCopy()
			initCR.Spec.License.Accept = testCase.isLicenseAccept
			testCase.initObjects = append(testCase.initObjects, initCR)

			cli := fake.NewFakeClientWithScheme(scheme, testCase.initObjects...)
			nonCachedClient := fake.NewFakeClientWithScheme(scheme, testCase.initObjects...)
			clientSet := fakeclientset.NewSimpleClientset()
			fakeDiscovery, ok := clientSet.Discovery().(*fakeDiscovery.FakeDiscovery)
			fakeDiscovery.Fake.Resources = []*metav1.APIResourceList{}

			if !ok {
				t.Fatal("Error creating fake discovery client")
			}

			r := NewReconciler(cli, nonCachedClient, fakeDiscovery, scheme, "")
			r.tests = true

			os.Setenv("CHE_FLAVOR", testCase.cheFlavor)
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      testCase.cheFlavor,
					Namespace: namespace,
				},
			}

			_, err := r.Reconcile(context.TODO(), req)
			if err != nil {
				t.Fatalf("Error reconciling: %v", err)
			}
			cr := &orgv1.WaziLicense{}
			if err := r.client.Get(context.TODO(), types.NamespacedName{Name: waziLicenseName, Namespace: namespace}, cr); err != nil {
				t.Errorf("CR not found")
			}

			if !reflect.DeepEqual(testCase.expectedNativeUserValue, cr.Spec.License.Accept) {
				expectedValue, actualValue := "nil", "nil"
				if !testCase.expectedNativeUserValue {
					expectedValue = strconv.FormatBool(testCase.expectedNativeUserValue)
				}
				if !cr.Spec.License.Accept {
					actualValue = strconv.FormatBool(cr.Spec.License.Accept)
				}

				t.Errorf("Expected nativeUserMode '%+v', but found '%+v' for input '%+v'",
					expectedValue, actualValue, testCase)
			}
		})
	}
}

func TestWaziController(t *testing.T) {
	var err error

	cl, dc, scheme := Init()

	// Create a ReconcileChe object with the scheme and fake client
	r := NewReconciler(cl, cl, dc, &scheme, "")
	r.tests = true

	// get CR
	waziCR := &orgv1.WaziLicense{
		Spec: orgv1.WaziLicenseSpec{
			License: orgv1.WaziLicenseSpecLicense{
				Accept: false,
			},
		},
	}

	if err := cl.Get(context.TODO(), types.NamespacedName{Name: waziLicenseName, Namespace: namespace}, waziCR); err != nil {
		t.Errorf("CR not found")
	}

	// update CR and make sure Wazi configmap has been updated
	waziCR.Spec.License.Accept = true
	if err := cl.Update(context.TODO(), waziCR); err != nil {
		t.Error("Failed to update WaziLicense custom resource")
	}

	if err = r.client.Get(context.TODO(), types.NamespacedName{Name: waziCR.Name, Namespace: waziCR.Namespace}, waziCR); err != nil {
		t.Errorf("Failed to get the Wazi custom resource %s: %s", waziCR.Name, err)
	}
}

func TestReconcileOFinalizers(t *testing.T) {
	type testCase struct {
		name              string
		deleteInitiated   bool
		expectedFinalizer []string
	}

	testCases := []testCase{
		{
			name:              "Condition: Deletion Time stamp activated",
			deleteInitiated:   true,
			expectedFinalizer: nil,
		},
		{
			name:              "Condition: Deleteion timestamp absent and test append finalizers",
			deleteInitiated:   false,
			expectedFinalizer: []string{util.WaziLicensingOperandRequestFinalizerName, util.WaziLicensingQuerySourceFinalizerName},
		},
	}
	for _, testCase := range testCases {
		cl, dc, scheme := Init()
		r := NewReconciler(cl, cl, dc, &scheme, "")
		deployContext := deploy.DeployContext{
			WaziLicense: initWaziWithTimeStamp(testCase.deleteInitiated),
			ClusterAPI: deploy.ClusterAPI{
				Client:           cl,
				NonCachingClient: cl,
				Scheme:           &scheme,
			},
		}
		r.reconcileFinalizers(&deployContext)
		testWaziCr := orgv1.WaziLicense{}
		if err := cl.Get(context.TODO(), types.NamespacedName{Name: waziLicenseName, Namespace: namespace}, &testWaziCr); err != nil {
			t.Errorf("Unable to Get Wazi-license from client error: %v", err)
		}
		if !reflect.DeepEqual(testCase.expectedFinalizer, testWaziCr.ObjectMeta.Finalizers) {
			t.Errorf("Expected Finalizers and Finalizers returned from API server for the given Wazi-License Spec are different (-want +got): %v", cmp.Diff(testCase.expectedFinalizer, testWaziCr.ObjectMeta.Finalizers))
		}
	}
}

func Init() (client.Client, discovery.DiscoveryInterface, runtime.Scheme) {
	objs, ds, scheme := createAPIObjects()

	// Register operator types with the runtime scheme

	// Create a fake client to mock API calls
	return fake.NewFakeClient(objs...), ds, scheme
}

func createAPIObjects() ([]runtime.Object, discovery.DiscoveryInterface, runtime.Scheme) {

	// A CheCluster custom resource with metadata and spec
	testWaziCR := InitWaziWithSimpleCR()
	testWaziPod := returnTestPodsWithAnnotations(map[string]string{})
	testSubscription := initWaziSubscription("ibm-wazi-developer-for-workspaces")
	testCSV := initWaziCSVWithVersion()

	// Objects to track in the fake client.
	objs := []runtime.Object{
		testWaziCR, testWaziPod, testSubscription, testCSV,
	}

	// Register operator types with the runtime scheme
	scheme := scheme.Scheme
	scheme.AddKnownTypes(orgv1.GroupVersion, testWaziCR)
	chev1alpha1.AddToScheme(scheme)
	packagesv1.AddToScheme(scheme)
	operatorsv1.AddToScheme(scheme)
	operatorsv1alpha1.AddToScheme(scheme)
	licensingoperatorv1.AddToScheme(scheme)
	odlmv1alpha1.AddToScheme(scheme)

	cli := fakeclientset.NewSimpleClientset()
	fakeDiscovery, ok := cli.Discovery().(*fakeDiscovery.FakeDiscovery)
	fakeDiscovery.Fake.Resources = []*metav1.APIResourceList{
		{TypeMeta: metav1.TypeMeta{
			Kind:       "operator.ibm.com",
			APIVersion: "v1",
		},
			GroupVersion: "operator.ibm.com/v1",
			APIResources: []metav1.APIResource{
				{
					Name:    "ibmlicensingquerysources",
					Group:   "operator.ibm.com",
					Version: "v1",
				},
			},
		},
	}
	if !ok {
		logrus.Error("Error creating fake discovery client")
		os.Exit(1)
	}

	// Create a fake client to mock API calls
	return objs, fakeDiscovery, *scheme
}

func InitWaziWithSimpleCR() *orgv1.WaziLicense {
	return &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Name:      waziLicenseName,
			Namespace: namespace,
		},
		Spec: orgv1.WaziLicenseSpec{
			License: orgv1.WaziLicenseSpecLicense{
				Accept: true,
			},
		},
	}
}

func initWaziWithOperandRequest(name string) *odlmv1alpha1.OperandRequest {
	return &odlmv1alpha1.OperandRequest{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func returnExpectedCR(opExists, opStatus bool) *odlmv1alpha1.OperandRequest {
	if !opExists {
		return nil
	}
	if opExists && !opStatus {
		return &odlmv1alpha1.OperandRequest{
			ObjectMeta: v1.ObjectMeta{
				Name:      testOperandRequestName,
				Namespace: namespace,
			},
		}
	}
	return &odlmv1alpha1.OperandRequest{
		ObjectMeta: v1.ObjectMeta{
			Name:      testOperandRequestName,
			Namespace: namespace,
		},
		Status: odlmv1alpha1.OperandRequestStatus{
			Members: []odlmv1alpha1.MemberStatus{
				{Name: sampleSubscriptionName,
					Phase: odlmv1alpha1.MemberPhase{
						OperatorPhase: odlmv1alpha1.OperatorPhase("Running"),
					},
				},
			},
		},
	}
}

func initWaziWithTimeStamp(testTimeStamp bool) *orgv1.WaziLicense {
	if !testTimeStamp {
		return &orgv1.WaziLicense{
			ObjectMeta: metav1.ObjectMeta{
				Name:      waziLicenseName,
				Namespace: namespace,
			},
			Spec: orgv1.WaziLicenseSpec{
				License: orgv1.WaziLicenseSpecLicense{
					Accept: true,
				},
			},
		}
	}
	return &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Name:              waziLicenseName,
			Namespace:         namespace,
			DeletionTimestamp: &v1.Time{Time: metav1.Now().Time},
			Finalizers:        []string{"sample-finalizer"},
		},
		Spec: orgv1.WaziLicenseSpec{
			License: orgv1.WaziLicenseSpecLicense{
				Accept: true,
			},
		},
	}
}
