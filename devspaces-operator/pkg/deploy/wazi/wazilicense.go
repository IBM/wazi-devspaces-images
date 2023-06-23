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
	"context"
	"fmt"

	chev2 "github.com/eclipse-che/che-operator/api/v2"
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/eclipse-che/che-operator/pkg/common/utils"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	envVarWaziLicenseUsage = "WAZI_LICENSE_USAGE"
)

func ReloadWaziLicenseCR(ctx *chetypes.DeployContext) error {
	return ctx.ClusterAPI.Client.Get(
		context.TODO(),
		types.NamespacedName{Name: ctx.WaziLicense.Name, Namespace: ctx.WaziLicense.Namespace},
		ctx.WaziLicense)
}

func GetWaziLicenseCR(ctx *chetypes.DeployContext) (*chev2.WaziLicense, bool, error) {

	ns := ctx.WaziLicense.Namespace
	return GetWaziLicenseCRWithNamespace(ctx, ns)
}

func GetWaziLicenseCRWithNamespace(ctx *chetypes.DeployContext, ns string) (*chev2.WaziLicense, bool, error) {

	cl := ctx.ClusterAPI.Client
	return GetWaziLicenseCRInNamespace(cl, ns)
}

func GetWaziLicenseCRInNamespace(cl client.Client, ns string) (*chev2.WaziLicense, bool, error) {

	waziLicenses := &chev2.WaziLicenseList{}

	if err := cl.List(
		context.TODO(),
		waziLicenses,
		&client.ListOptions{Namespace: ns}); err != nil {
		return nil, false, err
	}

	noOfWaziLicenses := len(waziLicenses.Items)
	if noOfWaziLicenses > 1 {
		logrus.Warn(fmt.Sprintf("expected one instance of WaziLicense custom resources, but '%d' found", len(waziLicenses.Items)))
	} else if noOfWaziLicenses < 1 {
		return nil, false, nil
	}

	namespacedName := types.NamespacedName{
		Namespace: waziLicenses.Items[0].GetNamespace(),
		Name:      waziLicenses.Items[0].GetName(),
	}

	return GetWaziLicenseCRWithNamespacedName(cl, namespacedName)
}

func GetWaziLicenseCRWithNamespacedName(cl client.Client, namespacedName types.NamespacedName) (*chev2.WaziLicense, bool, error) {

	waziLicense := &chev2.WaziLicense{}

	if err := cl.Get(
		context.TODO(),
		namespacedName,
		waziLicense); err != nil {
		return nil, false, err
	}

	return waziLicense, true, nil
}

func AppendFinalizer(ctx *chetypes.DeployContext, finalizer string) error {

	if err := ReloadWaziLicenseCR(ctx); err != nil {
		return err
	}

	if !utils.Contains(ctx.WaziLicense.ObjectMeta.Finalizers, finalizer) {
		for {
			ctx.WaziLicense.ObjectMeta.Finalizers = append(ctx.WaziLicense.ObjectMeta.Finalizers, finalizer)
			err := ctx.ClusterAPI.Client.Update(context.TODO(), ctx.WaziLicense)
			if err == nil {
				logrus.Infof("Added finalizer: %s", finalizer)
				return nil
			} else if !errors.IsConflict(err) {
				return err
			}

			err = ReloadWaziLicenseCR(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteFinalizer(ctx *chetypes.DeployContext, finalizer string) error {
	if utils.Contains(ctx.WaziLicense.ObjectMeta.Finalizers, finalizer) {
		for {
			ctx.WaziLicense.ObjectMeta.Finalizers = utils.Remove(ctx.WaziLicense.ObjectMeta.Finalizers, finalizer)
			err := ctx.ClusterAPI.Client.Update(context.TODO(), ctx.WaziLicense)
			if err == nil {
				logrus.Infof("Deleted finalizer: %s", finalizer)
				return nil
			} else if !errors.IsConflict(err) {
				return err
			}

			err = ReloadWaziLicenseCR(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetWaziAnnotations(ctx *chetypes.DeployContext) map[string]string {

	switch FromLicenseUsageString(ctx.WaziLicense.Spec.License.Use) {

	case IDzEE:
		return IDzEEAnnotations
	default:
		annotations := WaziAnnotations
		updateAnnotationProductVersion(annotations, ctx)
		return annotations
	}
}

func updateAnnotationProductVersion(annotations map[string]string, ctx *chetypes.DeployContext) (bool, error) {

	cl := ctx.ClusterAPI.Client
	ns := ctx.WaziLicense.Namespace

	subscription := &operatorsv1alpha1.Subscription{}
	clusterserviceversion := &operatorsv1alpha1.ClusterServiceVersion{}

	if err := cl.Get(context.TODO(), types.NamespacedName{Name: WaziPackageName, Namespace: ns}, subscription); err != nil {
		return false, err
	}

	if csvName := subscription.Status.InstalledCSV; len(csvName) > 0 {

		if err := cl.Get(context.TODO(), types.NamespacedName{Name: csvName, Namespace: ns}, clusterserviceversion); err != nil {
			return false, err
		}
		annotations["productVersion"] = clusterserviceversion.Spec.Version.String()
	}

	return true, nil
}

func AddWaziLicenseUsageEnvVar(ctx *chetypes.DeployContext, namespace string, env *[]v1.EnvVar) (bool, error) {

	waziLicense, done, err := GetWaziLicenseCRWithNamespace(ctx, namespace)
	if !done {
		return false, err
	}

	if waziLicense != nil {
		*env = append(*env, v1.EnvVar{
			Name:  envVarWaziLicenseUsage,
			Value: FromLicenseUsageString(waziLicense.Spec.License.Use).value,
		})
		return true, nil
	}

	return false, nil
}

func IsWaziLicense(cl client.Client, watchNamespace string, obj client.Object) (bool, ctrl.Request) {
	if obj.GetNamespace() == "" {
		return false, ctrl.Request{}
	}

	wazilicense, _, err := GetWaziLicenseCRInNamespace(cl, watchNamespace)

	if err != nil {
		logrus.Warn("a wazilicense Custom Resource was not found.")
		return false, ctrl.Request{}
	}

	if wazilicense.Namespace != obj.GetNamespace() {
		return false, ctrl.Request{}
	}

	return true, ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: wazilicense.Namespace,
			Name:      wazilicense.Name,
		},
	}
}
