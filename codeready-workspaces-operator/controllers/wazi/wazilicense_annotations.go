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

	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *WaziLicenseReconciler) ReconcileWaziAnnotations(deployContext *deploy.DeployContext) (done bool, err error) {

	annotations := GetWaziAnnotations(deployContext)

	operatorPod, _, err := util.FindWaziPodInNamespace(r.client, r.namespace)
	if err != nil {
		return true, err
	}

	otherAnnoEntries := make(map[string]string)

	for annoKey, annoValue := range annotations {
		if _, ok := operatorPod.Annotations[annoKey]; ok {
			operatorPod.Annotations[annoKey] = annoValue
		} else {
			otherAnnoEntries[annoKey] = annoValue
		}
	}

	if len(otherAnnoEntries) > 0 {
		podAnnotations := operatorPod.GetAnnotations()
		for k, v := range podAnnotations {
			otherAnnoEntries[k] = v
		}
		operatorPod.SetAnnotations(otherAnnoEntries)
	}

	if err := deployContext.ClusterAPI.Client.Update(context.TODO(), operatorPod); err != nil {
		if errors.IsConflict(err) {
			return false, err
		}
		if errors.IsNotFound(err) {
			return false, err
		}
		r.Log.Error(err, "Unable to update annotations on the Operator Pod")
		return true, err
	}

	return true, nil
}

func GetWaziAnnotations(deployContext *deploy.DeployContext) map[string]string {

	licenseUsage := util.GetWaziLicenseUsage(deployContext.WaziLicense)
	if util.FromLicenseUsageString(licenseUsage) == util.IDzEE {
		return util.IDzEEAnnotations
	}

	waziAnnotations := util.WaziAnnotations

	// Update Product Version based upon the CSV
	waziSubscription, _ := util.FindWaziSubscriptionInNamespace(deployContext.ClusterAPI.Client, deployContext.WaziLicense.Namespace)
	if waziSubscription != nil {

		waziCSVName := waziSubscription.Status.InstalledCSV

		if len(waziCSVName) > 0 {

			waziCSV, _ := util.FindWaziClusterServiceVersionInNamespace(deployContext.ClusterAPI.Client, deployContext.WaziLicense.Namespace, waziCSVName)
			if waziCSV != nil {
				waziAnnotations["productVersion"] = waziCSV.Spec.Version.String()
			}
		}
	}

	return waziAnnotations
}
