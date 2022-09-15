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
//   IBM Corporation - implementation
//

package che

import (
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/deploy/tls"
	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IsTrustedBundleConfigMap detects whether given config map is the config map with additional CA certificates to be trusted by Che
func IsTrustedBundleConfigMap(cl client.Client, watchNamespace string, obj client.Object) (bool, ctrl.Request) {
	if obj.GetNamespace() == "" {
		// ignore cluster scope objects
		return false, ctrl.Request{}
	}

	checluster, num, _ := util.FindCheClusterCRInNamespace(cl, watchNamespace)
	if num != 1 {
		if num > 1 {
			logrus.Warn("More than one checluster Custom Resource found.")
		}
		return false, ctrl.Request{}
	}

	if checluster.Namespace != obj.GetNamespace() {
		// ignore object in another namespace
		return false, ctrl.Request{}
	}

	// Check if config map is the config map from CR
	if checluster.Spec.Server.ServerTrustStoreConfigMapName != obj.GetName() {
		// No, it is not form CR

		// Check for component
		if value, exists := obj.GetLabels()[deploy.KubernetesComponentLabelKey]; !exists || value != tls.CheCACertsConfigMapLabelValue {
			// Labels do not match
			return false, ctrl.Request{}
		}

		// Check for part-of
		if value, exists := obj.GetLabels()[deploy.KubernetesPartOfLabelKey]; !exists || value != deploy.CheEclipseOrg {
			// ignore not matched labels
			return false, ctrl.Request{}
		}
	}

	return true, ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: checluster.Namespace,
			Name:      checluster.Name,
		},
	}
}

// isEclipseCheRelatedObj indicates if there is an object in a che namespace with the labels:
// - 'app.kubernetes.io/part-of=che.eclipse.org'
// - 'app.kubernetes.io/instance=che'
func IsEclipseCheRelatedObj(cl client.Client, watchNamespace string, obj client.Object) (bool, ctrl.Request) {
	if obj.GetNamespace() == "" {
		// ignore cluster scope objects
		return false, ctrl.Request{}
	}

	checluster, num, _ := util.FindCheClusterCRInNamespace(cl, watchNamespace)
	if num != 1 {
		if num > 1 {
			logrus.Warn("More than one checluster Custom Resource found.")
		}
		return false, ctrl.Request{}
	}

	if checluster.Namespace != obj.GetNamespace() {
		// ignore object in another namespace
		return false, ctrl.Request{}
	}

	// Check for part-of label
	if value, exists := obj.GetLabels()[deploy.KubernetesPartOfLabelKey]; !exists || value != deploy.CheEclipseOrg {
		return false, ctrl.Request{}
	}

	return true, ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: checluster.Namespace,
			Name:      checluster.Name,
		},
	}
}

// IsWaziLicenseRelatedObj indicates if there is an object in a che namespace with the Wazi License:
func IsWaziLicenseRelatedObj(cl client.Client, watchNamespace string, obj client.Object) (bool, ctrl.Request) {
	if obj.GetNamespace() == "" {
		// ignore cluster scope objects
		return false, ctrl.Request{}
	}

	// Wazi License Related Object
	wazilicense, num, _ := util.FindWaziLicenseCRInNamespace(cl, watchNamespace)
	if num != 1 {
		if num > 1 {
			logrus.Warn("More than one wazilicense Custom Resource found.")
		}
		return false, ctrl.Request{}
	}

	if wazilicense.Namespace != obj.GetNamespace() {
		// ignore object in another namespace
		return false, ctrl.Request{}
	}

	checluster, num, err := util.FindCheClusterCRInNamespace(cl, watchNamespace)

	if num != 1 || err != nil {
		if num > 1 {
			logrus.Warn("More than one checluster Custom Resource found.")
		}
		return false, ctrl.Request{}
	}

	return true, ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: watchNamespace,
			Name:      checluster.Name,
		},
	}
}
