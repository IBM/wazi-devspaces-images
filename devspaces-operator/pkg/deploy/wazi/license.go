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
	licensequerysrc "github.com/IBM/ibm-licensing-operator/api/v1"
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/eclipse-che/che-operator/pkg/common/constants"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	waziLicensingResourceGroup = "operator.ibm.com"
	waziLicensingVersion       = "v1"
	waziLicensingResourceName  = "ibmlicensingquerysources"
	waziLicensingFinalizer     = "ibmlicensingquerysources" + constants.FinalizerSuffix
	waziLicensingKind          = "IBMLicensingQuerySource"
	waziLicensingAPIVersion    = "operator.ibm.com/v1"
	waziLicensingName          = "ibm-wazi-code-licensing-query-source"
	waziLicensinQuery          = "count(container_processes{namespace=~'.+-devspaces\\\\W\\\\w+',container=~'wazi'})"
)

var licenseDiffOpts = cmp.Options{
	cmpopts.IgnoreFields(licensequerysrc.IBMLicensingQuerySource{}, "TypeMeta", "ObjectMeta"),
}

type LicenseReconciler struct {
	deploy.Reconcilable
}

func NewLicenseReconciler() *LicenseReconciler {
	return &LicenseReconciler{}
}

func (d *LicenseReconciler) Reconcile(ctx *chetypes.DeployContext) (reconcile.Result, bool, error) {

	if d.isLicenseQuerySource(ctx) {

		licensingSpec := d.getLicensingSpec(ctx)

		if _, err := WaziSync(ctx, licensingSpec, licenseDiffOpts); err != nil {
			return reconcile.Result{}, false, err
		}

		err := AppendFinalizer(ctx, waziLicensingFinalizer)
		if err != nil {
			return reconcile.Result{}, false, err
		}

		return reconcile.Result{}, true, nil
	}

	return reconcile.Result{}, false, nil
}

func (d *LicenseReconciler) Finalize(ctx *chetypes.DeployContext) bool {
	done := true
	if err := DeleteFinalizer(ctx, waziLicensingFinalizer); err != nil {
		done = false
		logrus.Errorf("Error deleting finalizer: %v", err)
	}
	return done
}

func (d *LicenseReconciler) isLicenseQuerySource(ctx *chetypes.DeployContext) bool {

	resources, err := ctx.ClusterAPI.DiscoveryClient.ServerResourcesForGroupVersion(
		schema.GroupVersion{
			Group:   waziLicensingResourceGroup,
			Version: waziLicensingVersion}.String())

	if err != nil {
		return false
	}

	for _, resource := range resources.APIResources {

		if resource.Name == waziLicensingResourceName {
			return true
		}
	}
	return false
}

func (d *LicenseReconciler) getLicensingSpec(ctx *chetypes.DeployContext) *licensequerysrc.IBMLicensingQuerySource {

	ns := ctx.WaziLicense.Namespace
	annotations := GetWaziAnnotations(ctx)

	licensingSpec := &licensequerysrc.IBMLicensingQuerySource{
		TypeMeta: metav1.TypeMeta{
			Kind:       waziLicensingKind,
			APIVersion: waziLicensingAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      waziLicensingName,
			Namespace: ns,
		},
		Spec: licensequerysrc.IBMLicensingQuerySourceSpec{
			Query:       waziLicensinQuery,
			Annotations: annotations,
		},
	}

	return licensingSpec
}
