//
// Copyright (c) 2019-2021 Red Hat, Inc.
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
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WaziAppendFinalizer(deployContext *DeployContext, finalizer string) error {
	if !util.ContainsString(deployContext.WaziLicense.ObjectMeta.Finalizers, finalizer) {

		err := ReloadWaziLicenseCR(deployContext)
		if err != nil {
			return err
		}

		for {
			deployContext.WaziLicense.ObjectMeta.Finalizers = append(deployContext.WaziLicense.ObjectMeta.Finalizers, finalizer)
			err = deployContext.ClusterAPI.Client.Update(context.TODO(), deployContext.WaziLicense)
			if err == nil {
				logrus.Infof("Added finalizer: %s", finalizer)
				return nil
			} else if !errors.IsConflict(err) {
				return err
			}

			err = ReloadWaziLicenseCR(deployContext)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func WaziDeleteFinalizer(deployContext *DeployContext, finalizer string) error {
	if util.ContainsString(deployContext.WaziLicense.ObjectMeta.Finalizers, finalizer) {
		for {
			deployContext.WaziLicense.ObjectMeta.Finalizers = util.DoRemoveString(deployContext.WaziLicense.ObjectMeta.Finalizers, finalizer)
			err := deployContext.ClusterAPI.Client.Update(context.TODO(), deployContext.WaziLicense)
			if err == nil {
				logrus.Infof("Deleted finalizer: %s", finalizer)
				return nil
			} else if !errors.IsConflict(err) {
				return err
			}

			err = ReloadWaziLicenseCR(deployContext)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func WaziDeleteObjectWithFinalizer(deployContext *DeployContext, key client.ObjectKey, objectMeta client.Object, finalizer string) error {
	_, err := WaziDelete(deployContext, key, objectMeta)
	if err != nil {
		// failed to delete, shouldn't us prevent from removing finalizer
		logrus.Error(err)
	}

	return WaziDeleteFinalizer(deployContext, finalizer)
}

func WaziGetFinalizerName(prefix string) string {
	finalizer := prefix + ".finalizers.che.eclipse.org"
	diff := len(finalizer) - 63
	if diff > 0 {
		return finalizer[:len(finalizer)-diff]
	}
	return finalizer
}
