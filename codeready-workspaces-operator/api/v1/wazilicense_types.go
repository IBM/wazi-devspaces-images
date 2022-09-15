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

package v1

// Important: You must regenerate some generated code after modifying this file. At the root of the project:
// Run `make generate`. It will perform required changes:
// - update `api/v1/zz_generatedxxx` files;
// - In the updated `config/crd/bases/org_v1_ibm_crd.yaml`: Delete all the `required:` openAPI rules in the CRD OpenApi schema;
// - Rename the new `config/crd/bases/org_v1_ibm_crd.yaml` to `config/crd/bases/org_v1_ibm_crd.yaml` to override it.
// IMPORTANT These 2 last steps are important to ensure backward compatibility with already existing `CheCluster` CRs that were created when no schema was provided.

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
// Personal z/OS cloud IDE for development and testing - License.
type WaziLicenseSpec struct {
	// General configuration settings related to the License
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="License"
	License WaziLicenseSpecLicense `json:"license"`
}

// License contains licensing information that must be accepted
type WaziLicenseSpecLicense struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:checkbox"}
	// +kubebuilder:default=false
	Accept bool `json:"accept"`
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="License Usage",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:select:Wazi","urn:alm:descriptor:com.tectonic.ui:select:IDzEE"}
	Use string `json:"use,omitempty"`
}

// WaziLicenseStatus defines the observed state of the Wazi License
type WaziLicenseStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// The `WaziLicense` custom resource allows defining and managing a License installation
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +operator-sdk:csv:customresourcedefinitions:displayName="Desired configuration of the IBM licensing"
// +operator-sdk:csv:customresourcedefinitions:resources={{OperandRequest,v1alpha1},{IBMLicensingQuerySource,v1}}
// +kubebuilder:storageversion
type WaziLicense struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WaziLicenseSpec   `json:"spec,omitempty"`
	Status            WaziLicenseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// WaziLicenseList contains a list of WaziLicense
type WaziLicenseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WaziLicense `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WaziLicense{}, &WaziLicenseList{})
}
