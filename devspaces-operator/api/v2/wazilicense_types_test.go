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

package v2

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWaziStructs(t *testing.T) {

	type testCase struct {
		name            string
		licenseAccept   bool
		licenseUsage    string
		expectedOutcome bool
	}

	testCases := []testCase{

		{
			name:            "Wazi License Default",
			licenseAccept:   true,
			licenseUsage:    "Wazi",
			expectedOutcome: true,
		},
		{
			name:            "Wazi License IDzEE",
			licenseAccept:   true,
			licenseUsage:    "IDzEE",
			expectedOutcome: true,
		},
		{
			name:            "Wazi License Blank",
			licenseAccept:   true,
			licenseUsage:    "",
			expectedOutcome: true,
		},
	}

	waziLicenses := getFakeWaziLicenseList()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			actualOutcome := contains(waziLicenses, testCase.licenseAccept, testCase.licenseUsage)
			if actualOutcome != testCase.expectedOutcome {
				t.Errorf("'%s' failed with actualOutcome %t versus expectedOutcome %t", testCase.name, actualOutcome, testCase.expectedOutcome)
			}
		})
	}
}

func contains(wazilicenses *WaziLicenseList, accept bool, usage string) bool {

	for _, wazilicense := range wazilicenses.Items {

		if accept && (usage != "") {
			if (wazilicense.Spec.License.Accept == accept) && (wazilicense.Spec.License.Use == usage) {
				return true
			}
		} else {
			if (wazilicense.Spec.License.Accept == accept) || (wazilicense.Spec.License.Use == usage) {
				return true
			}
		}
	}
	return false
}

func getFakeWaziLicenseList() *WaziLicenseList {
	waziLicenses := []WaziLicense{
		*getFakeWaziLicense(true, "Wazi"),
		*getFakeWaziLicense(true, "IDzEE"),
		*getFakeWaziLicense(true, ""),
		*getFakeWaziLicense(false, ""),
	}

	return &WaziLicenseList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "WaziLicenseList",
			APIVersion: "v2",
		},
		ListMeta: metav1.ListMeta{},
		Items:    waziLicenses,
	}

}

func getFakeWaziLicense(accept bool, usageType string) *WaziLicense {
	return &WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "eclipse-che",
			Name:      "wazi-devspaces-license",
		},
		Spec: WaziLicenseSpec{
			License: WaziLicenseSpecLicense{
				Accept: accept,
				Use:    usageType,
			},
		},
	}
}
