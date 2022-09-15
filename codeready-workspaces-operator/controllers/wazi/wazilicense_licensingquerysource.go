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
	"context"

	licensingoperatorv1 "github.com/IBM/ibm-licensing-operator/api/v1"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/google/go-cmp/cmp/cmpopts"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

const (
	LicensingQuerySourceResourceName    = "ibmlicensingquerysources"
	LicensingQuerySourceResourceGroup   = "operator.ibm.com"
	LicensingQuerySourceResourceVersion = "v1"
)

var (
	licQuerySrcDiffOpts = cmpopts.IgnoreFields(licensingoperatorv1.IBMLicensingQuerySource{}, "TypeMeta", "ObjectMeta")
)

func (r *WaziLicenseReconciler) ReconcileWaziLicensingQuerySource(deployContext *deploy.DeployContext) (done bool, err error) {

	licensingQuerySourceResourceDefinitionExists, err := isLicensingQuerySourceResourceDefinitionExist(deployContext)

	if err != nil {
		return false, err
	}

	if licensingQuerySourceResourceDefinitionExists {
		licensingQuerySourceSpec := getWaziLicensingQuerySourceSpec(deployContext, util.WaziLicensingQuerySourceName)
		return deploy.WaziSync(deployContext, licensingQuerySourceSpec, licQuerySrcDiffOpts)
	}

	return false, nil
}

func getWaziLicensingQuerySourceSpec(deployContext *deploy.DeployContext, name string) *licensingoperatorv1.IBMLicensingQuerySource {

	annotations := GetWaziAnnotations(deployContext)
	licensingQuerySource := &licensingoperatorv1.IBMLicensingQuerySource{
		TypeMeta: metav1.TypeMeta{
			Kind:       util.WaziLicensingQuerySourceKind,
			APIVersion: util.WaziLicensingQuerySourceAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: deployContext.WaziLicense.Namespace,
		},
		Spec: licensingoperatorv1.IBMLicensingQuerySourceSpec{

			Query:       util.WaziLicensingQuerySourceQuery,
			Annotations: annotations,
		},
	}

	return licensingQuerySource
}

func getWaziLicensingQuerySource(deployContext *deploy.DeployContext, name string) (waziLicensingQuerySource *licensingoperatorv1.IBMLicensingQuerySource, err error) {
	waziLicensingQuerySource = &licensingoperatorv1.IBMLicensingQuerySource{}
	err = deployContext.ClusterAPI.Client.Get(context.TODO(), types.NamespacedName{Namespace: deployContext.WaziLicense.Namespace, Name: name}, waziLicensingQuerySource)
	if err != nil {
		return nil, err
	}
	return waziLicensingQuerySource, nil
}

func isLicensingQuerySourceResourceDefinitionExist(deployContext *deploy.DeployContext) (exists bool, err error) {

	resources, err := deployContext.ClusterAPI.DiscoveryClient.ServerResourcesForGroupVersion(schema.GroupVersion{Group: LicensingQuerySourceResourceGroup, Version: LicensingQuerySourceResourceVersion}.String())

	if err != nil {
		return false, err
	}

	for _, resource := range resources.APIResources {

		if resource.Name == LicensingQuerySourceResourceName {
			return true, nil
		}
	}

	return false, nil
}
