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
	"context"

	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"

	licensingoperatorv1 "github.com/IBM/ibm-licensing-operator/api/v1"
	odlmv1alpha1 "github.com/IBM/operand-deployment-lifecycle-manager/api/v1alpha1"
)

// UpdateWaziCRSpec - updates Wazi CR "spec" by field
func UpdateWaziCRSpec(deployContext *DeployContext, field string, value string) error {
	err := deployContext.ClusterAPI.Client.Update(context.TODO(), deployContext.WaziLicense)
	if err == nil {
		logrus.Infof("Custom resource spec %s updated with %s: %s", deployContext.WaziLicense.Name, field, value)
		return nil
	}
	return err
}

func ReloadWaziLicenseCR(deployContext *DeployContext) error {
	return deployContext.ClusterAPI.Client.Get(
		context.TODO(),
		types.NamespacedName{Name: deployContext.WaziLicense.Name, Namespace: deployContext.WaziLicense.Namespace},
		deployContext.WaziLicense)
}

func ReconcileOperandRequestFinalizer(deployContext *DeployContext) (err error) {
	if deployContext.WaziLicense.ObjectMeta.DeletionTimestamp.IsZero() {
		return WaziAppendFinalizer(deployContext, util.WaziLicensingOperandRequestFinalizerName)
	} else {
		return WaziDeleteObjectWithFinalizer(deployContext, types.NamespacedName{Name: util.WaziOperandRequestName}, &odlmv1alpha1.OperandRequest{}, util.WaziLicensingOperandRequestFinalizerName)
	}
}

func ReconcileLicensingQuerySourceFinalizer(deployContext *DeployContext) (err error) {
	if deployContext.WaziLicense.ObjectMeta.DeletionTimestamp.IsZero() {
		return WaziAppendFinalizer(deployContext, util.WaziLicensingQuerySourceFinalizerName)
	} else {
		return WaziDeleteObjectWithFinalizer(deployContext, types.NamespacedName{Name: util.WaziLicensingQuerySourceName}, &licensingoperatorv1.IBMLicensingQuerySource{}, util.WaziLicensingQuerySourceFinalizerName)
	}
}
