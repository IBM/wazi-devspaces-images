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

package test

import (
	licensequerysrc "github.com/IBM/ibm-licensing-operator/api/v1"
	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
	chev1 "github.com/eclipse-che/che-operator/api/v1"
	chev2 "github.com/eclipse-che/che-operator/api/v2"
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	securityv1 "github.com/openshift/api/security/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakeDiscovery "k8s.io/client-go/discovery/fake"
	fakeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	TestWaziLicenseName       = "wazi-devspaces-license"
	TestWaziLicenseKind       = "WaziLicense"
	TestWaziLicenseVersion    = "v2"
	TestWaziLicenseNamespace  = "eclipse-che"
	TestWaziLicenseUsageWazi  = "Wazi"
	TestWaziLicenseUsageIDzEE = "IDzEE"
	testOperandRequestName    = "ibm-wazi-code-operand-request"
	sampleSubscriptionName    = "sample-subscription"
)

func GetFakeWaziLicense(accept bool, usageType string) *chev2.WaziLicense {
	return &chev2.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: TestWaziLicenseNamespace,
			Name:      TestWaziLicenseName,
		},
		Spec: chev2.WaziLicenseSpec{
			License: chev2.WaziLicenseSpecLicense{
				Accept: accept,
				Use:    usageType,
			},
		},
	}
}

func GetOperandRequest(opExists, opStatus bool) *odlmv1alpha1.OperandRequest {
	if !opExists {
		return nil
	}
	if opExists && !opStatus {
		return &odlmv1alpha1.OperandRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testOperandRequestName,
				Namespace: TestWaziLicenseNamespace,
			},
		}
	}
	return &odlmv1alpha1.OperandRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testOperandRequestName,
			Namespace: TestWaziLicenseNamespace,
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

func SetLicenseQuerySourceOwnerReference(licenseQuerySource *licensequerysrc.IBMLicensingQuerySource) {
	licenseQuerySource.ObjectMeta.ResourceVersion = "1"
	licenseQuerySource.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         "org.eclipse.che/v2",
			Kind:               TestWaziLicenseKind,
			Name:               TestWaziLicenseName,
			Controller:         TrueWaziPointer(),
			BlockOwnerDeletion: TrueWaziPointer(),
		},
	}
}

func ReturnTestPodsWithAnnotations(annotations map[string]string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake-pg-pod",
			Namespace: TestWaziLicenseNamespace,
			Labels: map[string]string{
				"app": "wazi-devspaces-operator",
			},
			Annotations: annotations,
		},
	}
}

func GetWaziDeployContext(waziLicense *chev2.WaziLicense, initObjs []runtime.Object) *chetypes.DeployContext {

	if waziLicense == nil {
		waziLicense = GetFakeWaziLicense(true, TestWaziLicenseUsageWazi)
	}

	scheme := scheme.Scheme
	chev1.SchemeBuilder.AddToScheme(scheme)
	chev2.SchemeBuilder.AddToScheme(scheme)
	odlmv1alpha1.SchemeBuilder.AddToScheme(scheme)
	licensequerysrc.SchemeBuilder.AddToScheme(scheme)
	corev1.SchemeBuilder.AddToScheme(scheme)

	scheme.AddKnownTypes(crdv1.SchemeGroupVersion, &crdv1.CustomResourceDefinition{})
	scheme.AddKnownTypes(operatorsv1alpha1.SchemeGroupVersion, &operatorsv1alpha1.Subscription{})
	scheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Pod{})
	securityv1.Install(scheme)

	initObjs = append(initObjs, waziLicense)
	cli := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects(initObjs...).Build()
	clientSet := fakeclientset.NewSimpleClientset()
	fakeDiscovery, _ := clientSet.Discovery().(*fakeDiscovery.FakeDiscovery)

	return &chetypes.DeployContext{
		WaziLicense: waziLicense,
		ClusterAPI: chetypes.ClusterAPI{
			Client:           cli,
			NonCachingClient: cli,
			Scheme:           scheme,
			DiscoveryClient:  fakeDiscovery,
		},
		Proxy: &chetypes.Proxy{},
	}
}

func TrueWaziPointer() *bool {
	b := true
	return &b
}
