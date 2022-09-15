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
package util

const (
	WaziPackageName                            = "ibm-wazi-developer-for-workspaces"
	WaziLicenseClusterPermissionsFinalizerName = "wazilicense.clusterpermissions.finalizers.che.eclipse.org"
	WaziLicensingOperandRequestFinalizerName   = "operandrequests.finalizers.che.eclipse.org"
	WaziLicensingQuerySourceFinalizerName      = "ibmlicensingquerysources.finalizers.che.eclipse.org"
)

var WaziAnnotations = map[string]string{
	"cloudpakId":               "9d41d2d8126f4200b62ba1acc0dffa2e",
	"cloudpakName":             "IBM Wazi for Dev Spaces",
	"cloudpakMetric":           "VIRTUAL_PROCESSOR_CORE",
	"productID":                "0e775d0d3f354a5ca074a6a4398045f3",
	"productId":                "0e775d0d3f354a5ca074a6a4398045f3", // Remove when license service is fixed
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
	"productId":                "3fbe9adb503d47a486bf5d138feb8fb9", // Remove when license service is fixed
	"productName":              "IBM Developer for z/OS Enterprise Edition",
	"productMetric":            "VU_VALUE_UNIT",
	"productChargedContainers": "All",
	"productCloudpakRatio":     "",
	"productVersion":           "15.0.5",
}

// List of Common Services to request from ODLM
var WaziCommonServices = []string{"ibm-licensing-operator"}

const (
	WaziLicenseClusterRoleName          = "-wazilicense-clusterrole"
	WaziLicenseEnvUsage                 = "WAZI_LICENSE_USAGE"
	WaziOperandRequestKind              = "OperandRequest"
	WaziOperandRequestAPIVersion        = "operator.ibm.com/v1alpha1"
	WaziOperandRequestRegistry          = "common-service"
	WaziOperandRequestRegistryNamespace = "ibm-common-services"
	WaziOperandRequestName              = "ibm-wazi-code-operand-request"
	WaziLicensingQuerySourceKind        = "IBMLicensingQuerySource"
	WaziLicensingQuerySourceAPIVersion  = "operator.ibm.com/v1"
	WaziLicensingQuerySourceName        = "ibm-wazi-code-licensing-query-source"
	WaziLicensingQuerySourceQuery       = "count(container_processes{namespace=~'wazi-.+-.+',container=~'wazi-plugins\\\\w\\\\w\\\\w'})"
)
