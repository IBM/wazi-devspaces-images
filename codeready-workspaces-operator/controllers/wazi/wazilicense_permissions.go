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
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *WaziLicenseReconciler) reconcileWaziLicensePermissions(deployContext *deploy.DeployContext) (bool, error) {
	done, err := r.delegateWaziLicensePermissions(deployContext)
	if !done {
		return false, err
	}

	return true, nil
}

func (r *WaziLicenseReconciler) delegateWaziLicensePermissions(deployContext *deploy.DeployContext) (bool, error) {
	waziLicenseClusterRoleName := deploy.GetWaziLicenseClusterRoleName(deployContext.WaziLicense.Namespace)

	done, err := deploy.WaziSyncClusterRoleToCluster(deployContext, waziLicenseClusterRoleName, getWaziLicensePolicies())
	if !done {
		return false, err
	}

	err = deploy.WaziAppendFinalizer(deployContext, util.WaziLicenseClusterPermissionsFinalizerName)
	return err == nil, err
}

func (r *WaziLicenseReconciler) removeWaziLicensePermissions(deployContext *deploy.DeployContext) (bool, error) {
	waziLicenseClusterRoleName := deploy.GetWaziLicenseClusterRoleName(deployContext.WaziLicense.Namespace)

	done, err := deploy.WaziDelete(deployContext, types.NamespacedName{Name: waziLicenseClusterRoleName}, &rbac.ClusterRole{})
	if !done {
		return false, err
	}

	err = deploy.WaziDeleteFinalizer(deployContext, util.WaziLicenseClusterPermissionsFinalizerName)
	return err == nil, err
}

func (r *WaziLicenseReconciler) reconcileWaziLicensePermissionsFinalizers(deployContext *deploy.DeployContext) (bool, error) {
	if !deployContext.WaziLicense.ObjectMeta.DeletionTimestamp.IsZero() {

		return r.removeWaziLicensePermissions(deployContext)
	}

	return true, nil
}

func getWaziLicensePolicies() []rbac.PolicyRule {
	k8sPolicies := []rbac.PolicyRule{
		{
			APIGroups: []string{"route.openshift.io"},
			Resources: []string{"routes"},
			Verbs:     []string{"list", "get", "patch"},
		},
		{
			APIGroups: []string{"route.openshift.io"},
			Resources: []string{"routes/custom-host"},
			Verbs:     []string{"create"},
		},
	}
	return k8sPolicies
}
