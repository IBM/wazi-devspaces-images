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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
type WaziLicenseSpec struct {
	License WaziLicenseSpecLicense `json:"license"`
}

type WaziLicenseSpecLicense struct {
	Accept bool   `json:"accept"`
	Use    string `json:"use,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +k8s:openapi-gen=true
// +kubebuilder:storageversion
type WaziLicense struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WaziLicenseSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
type WaziLicenseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WaziLicense `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WaziLicense{}, &WaziLicenseList{})
}
