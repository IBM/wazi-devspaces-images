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
	"os"
	"testing"

	"github.com/eclipse-che/che-operator/pkg/deploy"
)

func TestReconcileWaziLicensePermissions(t *testing.T) {

	cl, dc, scheme := Init()
	r := NewReconciler(cl, cl, dc, &scheme, "")
	deployContext := deploy.DeployContext{
		WaziLicense: initWaziorIdzeeWithSimpleCR(""),
		ClusterAPI: deploy.ClusterAPI{
			Client:           r.client,
			NonCachingClient: r.client,
			Scheme:           r.Scheme,
		},
	}
	os.Setenv("CHE_FLAVOR", waziLicenseName)

	done, err := r.reconcileWaziLicensePermissions(&deployContext)
	if err != nil {
		t.Fatalf("error when resting reconcile wazi license permissions")
	}
	if !done {
		t.Fatalf("did not actually complete reconcile wazi license permissions")
	}
}
