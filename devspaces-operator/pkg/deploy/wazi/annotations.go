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
	"os"
	"strings"

	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	k8shelper "github.com/eclipse-che/che-operator/pkg/common/k8s-helper"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const WaziPackageName = "ibm-wazi-code"

var WaziAnnotations = map[string]string{
	"cloudpakId":               "9d41d2d8126f4200b62ba1acc0dffa2e",
	"cloudpakName":             "IBM Wazi for Dev Spaces",
	"cloudpakMetric":           "VIRTUAL_PROCESSOR_CORE",
	"productID":                "0e775d0d3f354a5ca074a6a4398045f3",
	"productName":              "Wazi Code",
	"productMetric":            "AUTHORIZED_USER",
	"productChargedContainers": "All",
	"productCloudpakRatio":     "5:1",
	"productVersion":           "", // Derived from the CSV spec.version
}

var IDzEEAnnotations = map[string]string{
	"cloudpakId":               "",
	"cloudpakName":             "",
	"cloudpakMetric":           "",
	"productID":                "3fbe9adb503d47a486bf5d138feb8fb9",
	"productName":              "IBM Developer for z/OS Enterprise Edition",
	"productMetric":            "VU_VALUE_UNIT",
	"productChargedContainers": "All",
	"productCloudpakRatio":     "",
	"productVersion":           "16.0.4",
}

type LicenseUsage struct {
	value string
}

var (
	Wazi  = LicenseUsage{"wazi"}
	IDzEE = LicenseUsage{"idzee"}
)

func (t LicenseUsage) String() string {
	return t.value
}

func FromLicenseUsageString(s string) LicenseUsage {

	if strings.EqualFold(s, IDzEE.value) {
		return IDzEE
	}
	return Wazi
}

type AnnotationsReconciler struct {
	deploy.Reconcilable
}

func NewAnnotationsReconciler() *AnnotationsReconciler {
	return &AnnotationsReconciler{}
}

func (d *AnnotationsReconciler) Reconcile(ctx *chetypes.DeployContext) (reconcile.Result, bool, error) {

	operatorPod, done, err := d.getOperatorPod(ctx)
	if !done {
		return reconcile.Result{}, false, err
	}

	d.updatePodAnnotations(ctx, operatorPod)
	done, err = d.finalize(ctx, operatorPod)

	return reconcile.Result{}, done, err
}

func (d *AnnotationsReconciler) Finalize(ctx *chetypes.DeployContext) bool {
	return true
}

func (d *AnnotationsReconciler) finalize(ctx *chetypes.DeployContext, instance *v1.Pod) (bool, error) {

	if err := ctx.ClusterAPI.Client.Update(context.TODO(), instance); err != nil {
		return false, err
	}
	return true, nil
}

func (d *AnnotationsReconciler) getOperatorPod(ctx *chetypes.DeployContext) (*v1.Pod, bool, error) {

	api := k8shelper.New().GetClientset().CoreV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "app=" + os.Getenv("CHE_FLAVOR") + "-operator",
	}

	podList, err := api.Pods(ctx.WaziLicense.Namespace).List(context.TODO(), listOptions)
	if err != nil {
		return nil, false, err
	}

	return &podList.Items[0], true, nil
}

func (d *AnnotationsReconciler) updatePodAnnotations(ctx *chetypes.DeployContext, pod *v1.Pod) {

	annotations := GetWaziAnnotations(ctx)
	mergedAnnotations := d.mergeAnnotations(pod.GetAnnotations(), annotations)
	pod.SetAnnotations(mergedAnnotations)
}

func (d *AnnotationsReconciler) mergeAnnotations(maps ...map[string]string) map[string]string {

	annotations := make(map[string]string)

	for _, aMap := range maps {
	sMap:
		for key, value := range aMap {
			if _, ok := annotations[key]; ok {
				annotations[key] = value
				continue sMap
			}

			annotations[key] = value
		}
	}
	return annotations
}
