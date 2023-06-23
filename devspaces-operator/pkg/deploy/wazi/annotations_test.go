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

	test "github.com/eclipse-che/che-operator/pkg/common/test"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var (
	annotationsReconciler   = NewAnnotationsReconciler()
	defaultLicenseUsageType = test.TestWaziLicenseUsageWazi
	defaultPodName          = "fake-pod"
)

func TestFinalize(t *testing.T) {
	testPod := getFakePodWithAnnotations(defaultLicenseUsageType)
	waziLicense := test.GetFakeWaziLicense(true, defaultLicenseUsageType)
	objs := []runtime.Object{testPod}
	ctx := test.GetWaziDeployContext(waziLicense, objs)
	_, err := annotationsReconciler.finalize(ctx, testPod)
	if err != nil {
		t.Errorf("Annotaions Reconciler finalize func failed: %v", err)
	}
	recievedPod := corev1.Pod{}
	err = ctx.ClusterAPI.Client.Get(context.TODO(), types.NamespacedName{Name: defaultPodName, Namespace: test.TestWaziLicenseNamespace}, &recievedPod)
	if err != nil {
		t.Errorf("Annotaions Reconciler finalize func failed to get Pod: %v", err)
	}
	if !reflect.DeepEqual(&recievedPod, testPod) {
		t.Errorf("Expected Wazi Operator Pod vs generated Wazi Operator Pod (-want +got): %v", cmp.Diff(recievedPod, testPod))
	}
}

func TestUpdateAnnotations(t *testing.T) {
	type testCase struct {
		name        string
		licenseType string
	}
	testCases := []testCase{
		{
			name:        "Test Wazi Annotations update",
			licenseType: test.TestWaziLicenseUsageWazi,
		},
		{
			name:        "Test IDzEE Annotations update",
			licenseType: test.TestWaziLicenseUsageIDzEE,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			ctx := test.GetWaziDeployContext(test.GetFakeWaziLicense(true, FromLicenseUsageString(testCase.licenseType).String()), []runtime.Object{})
			samplePod := corev1.Pod{}
			annotationsReconciler.updatePodAnnotations(ctx, &samplePod)
			if !reflect.DeepEqual(samplePod.Annotations, getFakePodWithAnnotations(testCase.licenseType).Annotations) {
				diff := cmp.Diff(samplePod.Annotations, annotationsReconciler.mergeAnnotations(GetWaziAnnotations(ctx), getFakePodWithAnnotations(testCase.licenseType).Annotations))
				t.Errorf("Expected Wazi Operator Pod Annotations vs generated Wazi Operator Pod Annotations (-want +got): %v", diff)
			}
		})
	}
}

func TestMergeAnnotations(t *testing.T) {
	type testCase struct {
		name                string
		annotationExists    bool
		ExistingAnnotations map[string]string
		ExpectedAnnotation  map[string]string
	}
	testCases := []testCase{
		{
			name:                "Merge Annotations with existing annotations",
			annotationExists:    true,
			ExistingAnnotations: map[string]string{"sample-annotation-key": "sample-annotation-value"},
		},
		{
			name:                "Merge Annotations without existing annotations",
			annotationExists:    false,
			ExistingAnnotations: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recievedAnnotations := annotationsReconciler.mergeAnnotations(WaziAnnotations, testCase.ExistingAnnotations)
			if !reflect.DeepEqual(recievedAnnotations, annotationsReconciler.mergeAnnotations(WaziAnnotations, testCase.ExistingAnnotations)) {
				t.Errorf("Expected Wazi Operator Pod Annotations vs generated Wazi Operator Pod Annotations (-want +got): %v", cmp.Diff(recievedAnnotations, annotationsReconciler.mergeAnnotations(WaziAnnotations, testCase.ExistingAnnotations)))
			}
		})
	}
}

func getFakePodWithAnnotations(usageType string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      defaultPodName,
			Namespace: test.TestWaziLicenseNamespace,
		},
	}
	if usageType == defaultLicenseUsageType {
		pod.Annotations = WaziAnnotations
		return pod
	}
	pod.Annotations = IDzEEAnnotations
	return pod

}
