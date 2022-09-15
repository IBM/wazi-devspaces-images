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

package deploy

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	WaziLicenseStatusSuccessReason                     string = "WaziLicenseReconciled"
	WaziLicenseStatusSuccessMessage                    string = "Wazi License Successfully Reconciled"
	WaziLicenseStatusFailedPermissionsReason           string = "WaziLicenseFailedPermissions"
	WaziLicenseStatusFailedPermissionsMessage          string = "Wazi License Failed Permissions"
	WaziLicenseStatusFailedAnnotationsReason           string = "WaziLicenseFailedAnnotations"
	WaziLicenseStatusFailedAnnotationsMessage          string = "Wazi License Failed Annotations"
	WaziLicenseStatusFailedODLMReason                  string = "WaziLicenseFailedODLM"
	WaziLicenseStatusFailedODLMMessage                 string = "Wazi License Failed Operand Request"
	WaziLicenseStatusFailedLicensingQuerySourceReason  string = "WaziLicenseFailedLicensingQuerySource"
	WaziLicenseStatusFailedLicensingQuerySourceMessage string = "Wazi License Failed Licensing Query Source"
)

func SetWaziLicenseStatusSuccess(deployContext *DeployContext) error {

	patch := client.MergeFrom(deployContext.WaziLicense.DeepCopy())

	meta.RemoveStatusCondition(&deployContext.WaziLicense.Status.Conditions, metav1.StatusFailure)
	meta.SetStatusCondition(&deployContext.WaziLicense.Status.Conditions, metav1.Condition{
		Status:  metav1.ConditionTrue,
		Type:    metav1.StatusSuccess,
		Reason:  WaziLicenseStatusSuccessReason,
		Message: WaziLicenseStatusSuccessMessage,
	})

	return deployContext.ClusterAPI.Client.Status().Patch(context.TODO(), deployContext.WaziLicense, patch)
}

func SetWaziLicenseStatusFailedPermissions(deployContext *DeployContext) error {

	return setWaziLicenseStatusFailed(deployContext, WaziLicenseStatusFailedPermissionsReason, WaziLicenseStatusFailedPermissionsMessage)
}

func SetWaziLicenseStatusFailedAnnotations(deployContext *DeployContext) error {

	return setWaziLicenseStatusFailed(deployContext, WaziLicenseStatusFailedAnnotationsReason, WaziLicenseStatusFailedAnnotationsMessage)
}

func SetWaziLicenseStatusFailedODLM(deployContext *DeployContext) error {

	return setWaziLicenseStatusFailed(deployContext, WaziLicenseStatusFailedODLMReason, WaziLicenseStatusFailedODLMMessage)
}

func SetWaziLicenseStatusFailedLicensingQuerySource(deployContext *DeployContext) error {

	return setWaziLicenseStatusFailed(deployContext, WaziLicenseStatusFailedLicensingQuerySourceReason, WaziLicenseStatusFailedLicensingQuerySourceMessage)
}

func setWaziLicenseStatusFailed(deployContext *DeployContext, reason string, message string) error {

	patch := client.MergeFrom(deployContext.WaziLicense.DeepCopy())

	meta.SetStatusCondition(&deployContext.WaziLicense.Status.Conditions, metav1.Condition{
		Status:  metav1.ConditionFalse,
		Type:    metav1.StatusFailure,
		Reason:  reason,
		Message: message,
	})

	return deployContext.ClusterAPI.Client.Status().Patch(context.TODO(), deployContext.WaziLicense, patch)
}
