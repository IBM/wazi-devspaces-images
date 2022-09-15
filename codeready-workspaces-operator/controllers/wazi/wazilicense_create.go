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
	"strconv"

	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
)

func (r *WaziLicenseReconciler) GenerateAndSaveFields(deployContext *deploy.DeployContext) (err error) {

	if util.IsOpenShift && !deployContext.WaziLicense.Spec.License.Accept {
		newLicenseAcceptValue := false
		deployContext.WaziLicense.Spec.License.Accept = newLicenseAcceptValue
		if err := deploy.UpdateWaziCRSpec(deployContext, "licenseAccept", strconv.FormatBool(newLicenseAcceptValue)); err != nil {
			return err
		}
	}

	if deployContext.WaziLicense.Spec.License.Use == "" {
		licenseUsage := util.GetWaziLicenseUsage(deployContext.WaziLicense)
		deployContext.WaziLicense.Spec.License.Use = licenseUsage
		if err := deploy.UpdateWaziCRSpec(deployContext, "licenseUsage", licenseUsage); err != nil {
			return err
		}
	}

	return nil
}
