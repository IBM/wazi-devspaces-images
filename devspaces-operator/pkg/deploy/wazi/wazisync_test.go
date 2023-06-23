// Copyright (c) 2019-2021 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//
//	IBM Corporation - initial API and implementation
package wazi

import (
	"context"

	chev2 "github.com/eclipse-che/che-operator/api/v2"
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/eclipse-che/che-operator/pkg/common/utils"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"testing"
)

var (
	testWaziObj = &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "eclipse-che",
		},
	}
	testWaziObjLabeled = &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-secret",
			Namespace:   "eclipse-che",
			Labels:      map[string]string{"a": "b"},
			Annotations: map[string]string{"d": "c"},
		},
		Data: map[string][]byte{"x": []byte("y")},
	}
	testWaziKey = client.ObjectKey{Name: "test-secret", Namespace: "eclipse-che"}
)

func TestWaziGet(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	err := cli.Create(context.TODO(), testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	actual := &corev1.Secret{}
	exists, err := WaziGet(deployContext, testWaziKey, actual)
	if !exists || err != nil {
		t.Fatalf("Failed to get object: %v", err)
	}
}

func TestWaziGetNamespacedObject(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	err := cli.Create(context.TODO(), testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	actual := &corev1.Secret{}
	done, err := WaziGetNamespacedObject(deployContext, testWaziKey.Name, actual)
	if err != nil {
		t.Fatalf("Failed to get namespace object: %v", err)
	}
	if !done {
		t.Fatalf("Namespaced object was not found")
	}
}

func TestWaziCreate(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	done, err := WaziCreate(deployContext, testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	if !done {
		t.Fatalf("Object has not been created")
	}

	actual := &corev1.Secret{}
	err = cli.Get(context.TODO(), testWaziKey, actual)
	if err != nil && !errors.IsNotFound(err) {
		t.Fatalf("Failed to get object: %v", err)
	}
}

func TestWaziCreateIfNotExistsShouldReturnTrueIfObjectCreated(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	done, err := WaziCreateIfNotExists(deployContext, testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	if !done {
		t.Fatalf("Object has not been created")
	}

	actual := &corev1.Secret{}
	err = cli.Get(context.TODO(), testWaziKey, actual)
	if err != nil && !errors.IsNotFound(err) {
		t.Fatalf("Failed to get object: %v", err)
	}
}

func TestWaziCreateIfNotExistsShouldReturnFalseIfObjectExist(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	err := cli.Create(context.TODO(), testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	isCreated, err := WaziCreateIfNotExists(deployContext, testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	if isCreated {
		t.Fatalf("Object has been created")
	}
}

func TestWaziUpdate(t *testing.T) { //
	cli, deployContext := initWaziDeployContext()

	err := cli.Create(context.TODO(), testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	actual := &corev1.Secret{}
	err = cli.Get(context.TODO(), testWaziKey, actual)
	if err != nil && !errors.IsNotFound(err) {
		t.Fatalf("Failed to get object: %v", err)
	}

	_, err = WaziUpdateWithClient(cli, deployContext, actual, testWaziObjLabeled.DeepCopy(), cmp.Options{})
	if err != nil {
		t.Fatalf("Failed to update object: %v", err)
	}

	err = cli.Get(context.TODO(), testWaziKey, actual)
	if err != nil && !errors.IsNotFound(err) {
		t.Fatalf("Failed to get object: %v", err)
	}

	if actual.Labels["a"] != "b" {
		t.Fatalf("Object hasn't been updated")
	}
}

func TestWaziSyncAndAddFinalizer(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	cli.Create(context.TODO(), deployContext.WaziLicense)

	// Sync object
	done, err := WaziSyncAndAddFinalizer(deployContext, testWaziObj.DeepCopy(), cmp.Options{}, "test-finalizer")
	if !done || err != nil {
		t.Fatalf("Error syncing object: %v", err)
	}

	actual := &corev1.Secret{}
	err = cli.Get(context.TODO(), testWaziKey, actual)
	if err != nil {
		t.Fatalf("Failed to get object: %v", err)
	}

	if !utils.Contains(deployContext.WaziLicense.Finalizers, "test-finalizer") {
		t.Fatalf("Failed to add finalizer")
	}
}

func TestWaziShouldDeleteExistedObject(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	err := cli.Create(context.TODO(), testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	done, err := WaziDelete(deployContext, testWaziKey, testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to delete object: %v", err)
	}

	if !done {
		t.Fatalf("Object hasn't been deleted")
	}

	actualObj := &corev1.Secret{}
	err = cli.Get(context.TODO(), testWaziKey, actualObj)
	if err != nil && !errors.IsNotFound(err) {
		t.Fatalf("Failed to get object: %v", err)
	}

	if err == nil {
		t.Fatalf("Object hasn't been deleted")
	}
}

func TestWaziShouldNotDeleteObject(t *testing.T) {
	_, deployContext := initWaziDeployContext()

	done, err := WaziDelete(deployContext, testWaziKey, testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to delete object: %v", err)
	}

	if !done {
		t.Fatalf("Object has not been deleted")
	}
}

func TestWaziShouldDeleteNamespacedObject(t *testing.T) {
	cli, deployContext := initWaziDeployContext()

	err := cli.Create(context.TODO(), testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to create object: %v", err)
	}

	done, err := WaziDeleteNamespacedObject(deployContext, testWaziKey.Name, testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to delete object: %v", err)
	}

	if !done {
		t.Fatalf("Object hasn't been deleted")
	}

	actualObj := &corev1.Secret{}
	err = cli.Get(context.TODO(), testWaziKey, actualObj)
	if err != nil && !errors.IsNotFound(err) {
		t.Fatalf("Failed to get object: %v", err)
	}

	if err == nil {
		t.Fatalf("Object hasn't been deleted")
	}
}

func TestWaziShouldNotDeleteNamespacedObject(t *testing.T) {
	_, deployContext := initWaziDeployContext()

	done, err := WaziDeleteNamespacedObject(deployContext, testWaziKey.Name, testWaziObj.DeepCopy())
	if err != nil {
		t.Fatalf("Failed to delete object: %v", err)
	}

	if !done {
		t.Fatalf("Object has not been deleted")
	}
}

func initWaziDeployContext() (client.Client, *chetypes.DeployContext) {
	waziLicense := &chev2.WaziLicense{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "eclipse-che",
			Name:      "eclipse-che",
		},
	}

	chev2.SchemeBuilder.AddToScheme(scheme.Scheme)

	cli := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithRuntimeObjects(waziLicense).Build()
	deployContext := &chetypes.DeployContext{
		WaziLicense: waziLicense,
		ClusterAPI: chetypes.ClusterAPI{
			Client:           cli,
			NonCachingClient: cli,
			Scheme:           scheme.Scheme,
		},
	}

	return cli, deployContext
}
