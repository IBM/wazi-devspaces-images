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
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetWaziLicenseStatusSuccess(t *testing.T) {
	_, deployContext := initWaziDeployContext()
	err := SetWaziLicenseStatusSuccess(deployContext)

	if err != nil {
		t.Fatalf("Error when setting wazi to success: %v", err)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Type != metav1.StatusSuccess {
		t.Fatalf("Failed to set wazi license type to %v, was: %v", metav1.StatusSuccess, deployContext.WaziLicense.Status.Conditions[0].Type)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Status != metav1.ConditionTrue {
		t.Fatalf("Failed to set wazi license status to %v, was : %v", metav1.ConditionTrue, deployContext.WaziLicense.Status.Conditions[0].Status)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Reason != WaziLicenseStatusSuccessReason {
		t.Fatalf("Failed to set wazi license reason to %v, was: %v", WaziLicenseStatusSuccessReason, deployContext.WaziLicense.Status.Conditions[0].Reason)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Message != WaziLicenseStatusSuccessMessage {
		t.Fatalf("Failed to set wazi license message to %v, was: %v", WaziLicenseStatusSuccessMessage, deployContext.WaziLicense.Status.Conditions[0].Message)
	}
}
func TestSetWaziLicenseStatusFailedPermissions(t *testing.T) {
	_, deployContext := initWaziDeployContext()
	err := SetWaziLicenseStatusFailedPermissions(deployContext)

	if err != nil {
		t.Fatalf("Error when setting wazi to failed permissions: %v", err)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Type != metav1.StatusFailure {
		t.Fatalf("Failed to set wazi license type to %v, was: %v", metav1.StatusFailure, deployContext.WaziLicense.Status.Conditions[0].Type)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Status != metav1.ConditionFalse {
		t.Fatalf("Failed to set wazi license status to %v, was : %v", metav1.ConditionFalse, deployContext.WaziLicense.Status.Conditions[0].Status)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Reason != WaziLicenseStatusFailedPermissionsReason {
		t.Fatalf("Failed to set wazi license reason to %v, was: %v", WaziLicenseStatusFailedPermissionsReason, deployContext.WaziLicense.Status.Conditions[0].Reason)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Message != WaziLicenseStatusFailedPermissionsMessage {
		t.Fatalf("Failed to set wazi license message to %v, was: %v", WaziLicenseStatusFailedPermissionsMessage, deployContext.WaziLicense.Status.Conditions[0].Message)
	}
}
func TestSetWaziLicenseStatusFailedAnnotations(t *testing.T) {
	_, deployContext := initWaziDeployContext()
	err := SetWaziLicenseStatusFailedAnnotations(deployContext)

	if err != nil {
		t.Fatalf("Error when setting wazi to failed annotations: %v", err)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Type != metav1.StatusFailure {
		t.Fatalf("Failed to set wazi license type to %v, was: %v", metav1.StatusFailure, deployContext.WaziLicense.Status.Conditions[0].Type)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Status != metav1.ConditionFalse {
		t.Fatalf("Failed to set wazi license status to %v, was : %v", metav1.ConditionFalse, deployContext.WaziLicense.Status.Conditions[0].Status)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Reason != WaziLicenseStatusFailedAnnotationsReason {
		t.Fatalf("Failed to set wazi license reason to %v, was: %v", WaziLicenseStatusFailedAnnotationsReason, deployContext.WaziLicense.Status.Conditions[0].Reason)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Message != WaziLicenseStatusFailedAnnotationsMessage {
		t.Fatalf("Failed to set wazi license message to %v, was: %v", WaziLicenseStatusFailedAnnotationsMessage, deployContext.WaziLicense.Status.Conditions[0].Message)
	}
}
func TestSetWaziLicenseStatusFailedODLM(t *testing.T) {
	_, deployContext := initWaziDeployContext()
	err := SetWaziLicenseStatusFailedODLM(deployContext)

	if err != nil {
		t.Fatalf("Error when setting wazi to failed dolm: %v", err)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Type != metav1.StatusFailure {
		t.Fatalf("Failed to set wazi license type to %v, was: %v", metav1.StatusFailure, deployContext.WaziLicense.Status.Conditions[0].Type)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Status != metav1.ConditionFalse {
		t.Fatalf("Failed to set wazi license status to %v, was : %v", metav1.ConditionFalse, deployContext.WaziLicense.Status.Conditions[0].Status)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Reason != WaziLicenseStatusFailedODLMReason {
		t.Fatalf("Failed to set wazi license reason to %v, was: %v", WaziLicenseStatusFailedODLMReason, deployContext.WaziLicense.Status.Conditions[0].Reason)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Message != WaziLicenseStatusFailedODLMMessage {
		t.Fatalf("Failed to set wazi license message to %v, was: %v", WaziLicenseStatusFailedODLMMessage, deployContext.WaziLicense.Status.Conditions[0].Message)
	}
}
func TestSetWaziLicenseStatusFailedLicensingQuerySource(t *testing.T) {
	_, deployContext := initWaziDeployContext()
	err := SetWaziLicenseStatusFailedLicensingQuerySource(deployContext)

	if err != nil {
		t.Fatalf("Error when setting wazi license to failed licenseing query source: %v", err)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Type != metav1.StatusFailure {
		t.Fatalf("Failed to set wazi license type to %v, was: %v", metav1.StatusFailure, deployContext.WaziLicense.Status.Conditions[0].Type)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Status != metav1.ConditionFalse {
		t.Fatalf("Failed to set wazi license status to %v, was : %v", metav1.ConditionFalse, deployContext.WaziLicense.Status.Conditions[0].Status)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Reason != WaziLicenseStatusFailedLicensingQuerySourceReason {
		t.Fatalf("Failed to set wazi license reason to %v, was: %v", WaziLicenseStatusFailedLicensingQuerySourceReason, deployContext.WaziLicense.Status.Conditions[0].Reason)
	}
	if deployContext.WaziLicense.Status.Conditions[0].Message != WaziLicenseStatusFailedLicensingQuerySourceMessage {
		t.Fatalf("Failed to set wazi license message to %v, was: %v", WaziLicenseStatusFailedLicensingQuerySourceMessage, deployContext.WaziLicense.Status.Conditions[0].Message)
	}
}
