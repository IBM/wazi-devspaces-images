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
package util

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	orgv1 "github.com/eclipse-che/che-operator/api/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LicenseUsage struct {
	value string
}

var (
	Wazi  = LicenseUsage{"wazi"}
	IDzEE = LicenseUsage{"idzee"}
)

func (t LicenseUsage) String() string {
	return t.value
}

func FromLicenseUsageString(s string) LicenseUsage {

	if strings.EqualFold(s, IDzEE.value) {
		return IDzEE
	}
	return Wazi
}

func GetWaziLicenseUsage(waziLicense *orgv1.WaziLicense) string {
	if waziLicense.Spec.License.Use != "" {
		return (FromLicenseUsageString(waziLicense.Spec.License.Use)).value
	}
	return Wazi.value
}

func AddWaziLicenseEnv(cl client.Client, namespace string, env *[]corev1.EnvVar) {
	waziLicense, _, _ := FindWaziLicenseCRInNamespace(cl, namespace)
	if waziLicense != nil {
		*env = append(*env, corev1.EnvVar{
			Name:  WaziLicenseEnvUsage,
			Value: GetWaziLicenseUsage(waziLicense),
		})
	}
}

// Finds WaziLicense custom resource in a given namespace.
// If namespace is empty then WaziLicense will be found in any namespace.
func FindWaziLicenseCRInNamespace(cl client.Client, namespace string) (*orgv1.WaziLicense, int, error) {
	waziLicenses := &orgv1.WaziLicenseList{}
	listOptions := &client.ListOptions{Namespace: namespace}
	if err := cl.List(context.TODO(), waziLicenses, listOptions); err != nil {
		return nil, -1, err
	}
	if len(waziLicenses.Items) != 1 {
		return nil, len(waziLicenses.Items), fmt.Errorf("Expected one instance of WaziLicense custom resource, but '%d' found.", len(waziLicenses.Items))
	}
	waziLicense := &orgv1.WaziLicense{}
	namespacedName := types.NamespacedName{Namespace: waziLicenses.Items[0].GetNamespace(), Name: waziLicenses.Items[0].GetName()}
	err := cl.Get(context.TODO(), namespacedName, waziLicense)
	if err != nil {
		return nil, -1, err
	}
	return waziLicense, 1, nil
}

func FindWaziPodInNamespace(cl client.Client, namespace string) (*corev1.Pod, int, error) {

	waziPods := &corev1.PodList{}

	cheFlavor := os.Getenv("CHE_FLAVOR")
	if len(cheFlavor) == 0 {
		err := errors.New("Environment variable for CHE_FLAVOR not set.")
		logrus.Fatalf(err.Error())
		return nil, -1, err
	}

	selector, err := labels.Parse(fmt.Sprintf("%s=%s", "app", cheFlavor+"-operator"))
	if err != nil {
		return nil, -1, err
	}

	listOptions := &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: selector,
	}

	if err := cl.List(context.TODO(), waziPods, listOptions); err != nil {
		return nil, -1, err
	}

	if len(waziPods.Items) != 1 {
		return nil, len(waziPods.Items), fmt.Errorf("Expected one operator pod, but '%d' found.", len(waziPods.Items))
	}

	return &waziPods.Items[0], 1, nil
}

func FindWaziSubscriptionInNamespace(cl client.Client, namespace string) (*operatorsv1alpha1.Subscription, error) {

	waziSubscription := &operatorsv1alpha1.Subscription{}
	namespacedName := types.NamespacedName{Namespace: namespace, Name: WaziPackageName}
	if err := cl.Get(context.TODO(), namespacedName, waziSubscription); err != nil {
		return nil, err
	}

	return waziSubscription, nil
}

func FindWaziClusterServiceVersionInNamespace(cl client.Client, namespace string, name string) (*operatorsv1alpha1.ClusterServiceVersion, error) {

	waziCSV := &operatorsv1alpha1.ClusterServiceVersion{}
	namespacedName := types.NamespacedName{Namespace: namespace, Name: name}
	if err := cl.Get(context.TODO(), namespacedName, waziCSV); err != nil {
		return nil, err
	}

	return waziCSV, nil
}
