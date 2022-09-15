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
	orgv1 "github.com/eclipse-che/che-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"testing"
)

var (
	waziLicenseName = "wazi-license"
)

func TestUpdateWaziCRSpec(t *testing.T) {
	_, deployContext := initWaziLicenseDeployContext()

	updatedName := "wazi-lisense-updated"
	err := UpdateWaziCRSpec(deployContext, deployContext.WaziLicense.Name, updatedName)
	if err != nil {
		t.Errorf("CR spec not updated, %v", err)
	}

}

func TestReloadWaziLicense(t *testing.T) {
	cli, _ := initWaziLicenseDeployContext()

	waziLicense := &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       waziLicenseName,
			Name:            waziLicenseName,
			ResourceVersion: "2",
		},
	}

	deployContext := &DeployContext{
		WaziLicense: waziLicense,
		ClusterAPI: ClusterAPI{
			Client:           cli,
			NonCachingClient: cli,
			Scheme:           scheme.Scheme,
		},
	}

	err := ReloadWaziLicenseCR(deployContext)
	if err != nil {
		t.Errorf("Failed to reload wazilicense, %v", err)
	}

	if waziLicense.ObjectMeta.ResourceVersion != "1" {
		t.Errorf("Failed to reload wazilicense")
	}
}

func TestReconcileOperandRequestFinalizer(t *testing.T) {
	_, deployContext := initWaziLicenseDeployContext()

	err := ReconcileOperandRequestFinalizer(deployContext)
	if err != nil {
		t.Errorf("Failed to reconcile operand finalizer, %v", err)
	}
}

func TestReconcileLicensingQuerySourceFinalizer(t *testing.T) {
	_, deployContext := initWaziLicenseDeployContext()

	err := ReconcileLicensingQuerySourceFinalizer(deployContext)
	if err != nil {
		t.Errorf("Failed to reconcile query finalizer, %v", err)
	}
}

func initWaziLicenseDeployContext() (client.Client, *DeployContext) {
	waziLicense := &orgv1.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       waziLicenseName,
			Name:            waziLicenseName,
			ResourceVersion: "1",
		},
	}

	orgv1.SchemeBuilder.AddToScheme(scheme.Scheme)
	cli := fake.NewFakeClientWithScheme(scheme.Scheme, waziLicense)
	deployContext := &DeployContext{
		WaziLicense: waziLicense,
		ClusterAPI: ClusterAPI{
			Client:           cli,
			NonCachingClient: cli,
			Scheme:           scheme.Scheme,
		},
	}

	return cli, deployContext
}
