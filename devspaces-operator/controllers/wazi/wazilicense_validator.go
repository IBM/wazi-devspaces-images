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
	"fmt"

	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/deploy/wazi"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type WaziLicenseValidator struct {
	deploy.Reconcilable
}

func NewWaziLicenseValidator() *WaziLicenseValidator {
	return &WaziLicenseValidator{}
}

func (v *WaziLicenseValidator) Reconcile(ctx *chetypes.DeployContext) (reconcile.Result, bool, error) {

	ns := ""

	if ctx.WaziLicense != nil {
		ns = ctx.WaziLicense.Namespace
	} else if ctx.CheCluster != nil {
		ns = ctx.CheCluster.Namespace
	} else {
		return reconcile.Result{}, false, fmt.Errorf("unable to determine wazilicense namespace")
	}

	waziLicense, done, err := wazi.GetWaziLicenseCRInNamespace(ctx.ClusterAPI.Client, ns)
	if !done {
		if err != nil {
			return reconcile.Result{}, false, err
		} else {
			return reconcile.Result{}, false, fmt.Errorf("wazilicense not found")
		}
	}

	licenseAccept := waziLicense.Spec.License.Accept
	if !licenseAccept {
		return reconcile.Result{}, false, fmt.Errorf("required field \"spec.license.accept\" is not accepted")
	}

	return reconcile.Result{}, true, nil
}

func (v *WaziLicenseValidator) Finalize(ctx *chetypes.DeployContext) bool {
	return true
}
