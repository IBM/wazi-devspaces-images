//
// Copyright (c) 2019-2021 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//
package deploy

import (
	"os"

	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var crWaziDiffOpts = cmp.Options{
	cmpopts.IgnoreFields(rbac.ClusterRole{}, "TypeMeta", "ObjectMeta"),
	cmpopts.IgnoreFields(rbac.PolicyRule{}, "ResourceNames", "NonResourceURLs"),
}

func WaziSyncClusterRoleToCluster(
	deployContext *DeployContext,
	name string,
	policyRule []rbac.PolicyRule) (bool, error) {

	crWaziSpec := getWaziClusterRoleSpec(deployContext, name, policyRule)
	return WaziSync(deployContext, crWaziSpec, crWaziDiffOpts)
}

func getWaziClusterRoleSpec(deployContext *DeployContext, name string, policyRule []rbac.PolicyRule) *rbac.ClusterRole {
	labels := WaziGetLabels(os.Getenv("CHE_FLAVOR"))
	clusterRole := &rbac.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: rbac.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
			Annotations: map[string]string{
				CheEclipseOrgNamespace: deployContext.WaziLicense.Namespace,
			},
		},
		Rules: policyRule,
	}

	return clusterRole
}

func GetWaziLicenseClusterRoleName(namespace string) string {
	return namespace + util.WaziLicenseClusterRoleName
}
