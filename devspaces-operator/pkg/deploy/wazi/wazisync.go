// Copyright (c) 2019-2021 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//
//	Red Hat, Inc. - initial API and implementation
package wazi

import (
	"context"
	"fmt"
	"reflect"

	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Sync syncs the blueprint to the cluster in a generic (as much as Go allows) manner.
// Returns true if object is up-to-date otherwise returns false
func WaziSync(deployContext *chetypes.DeployContext, blueprint client.Object, diffOpts ...cmp.Option) (bool, error) {
	cli := getClientForObject(blueprint.GetNamespace(), deployContext)
	return WaziSyncWithClient(cli, deployContext, blueprint, diffOpts...)
}

func WaziSyncWithClient(cli client.Client, deployContext *chetypes.DeployContext, blueprint client.Object, diffOpts ...cmp.Option) (bool, error) {
	runtimeObject, ok := blueprint.(runtime.Object)
	if !ok {
		return false, fmt.Errorf("object %T is not a runtime.Object. Cannot sync it", runtimeObject)
	}

	// we will compare this object later with blueprint
	// we can't use runtimeObject.DeepCopyObject()
	actual, err := deployContext.ClusterAPI.Scheme.New(runtimeObject.GetObjectKind().GroupVersionKind())
	if err != nil {
		return false, err
	}

	key := types.NamespacedName{Name: blueprint.GetName(), Namespace: blueprint.GetNamespace()}
	exists, err := WaziGetWithClient(cli, key, actual.(client.Object))
	if err != nil {
		return false, err
	}

	// set GroupVersionKind (it might be empty)
	actual.GetObjectKind().SetGroupVersionKind(runtimeObject.GetObjectKind().GroupVersionKind())
	if !exists {
		return WaziCreateWithClient(cli, deployContext, blueprint, false)
	}

	return WaziUpdateWithClient(cli, deployContext, actual.(client.Object), blueprint, diffOpts...)
}

func WaziSyncAndAddFinalizer(
	deployContext *chetypes.DeployContext,
	blueprint metav1.Object,
	diffOpts cmp.Option,
	finalizer string) (bool, error) {

	// eclipse-che custom resource is being deleted, we shouldn't sync
	// TODO move this check before `Sync` invocation
	if deployContext.WaziLicense.ObjectMeta.DeletionTimestamp.IsZero() {
		done, err := WaziSync(deployContext, blueprint.(client.Object), diffOpts)
		if !done {
			return done, err
		}
		err = AppendFinalizer(deployContext, finalizer)
		return err == nil, err
	}
	return true, nil
}

// Gets object by key.
// Returns true if object exists otherwise returns false.
func WaziGet(deployContext *chetypes.DeployContext, key client.ObjectKey, actual client.Object) (bool, error) {
	cli := getClientForObject(key.Namespace, deployContext)
	return WaziGetWithClient(cli, key, actual)
}

// Gets namespaced scope object by name
// Returns true if object exists otherwise returns false.
func WaziGetNamespacedObject(deployContext *chetypes.DeployContext, name string, actual client.Object) (bool, error) {
	client := deployContext.ClusterAPI.Client
	key := types.NamespacedName{Name: name, Namespace: deployContext.WaziLicense.Namespace}
	return WaziGetWithClient(client, key, actual)
}

// Gets cluster scope object by name
// Returns true if object exists otherwise returns false
func WaziGetClusterObject(deployContext *chetypes.DeployContext, name string, actual client.Object) (bool, error) {
	client := deployContext.ClusterAPI.NonCachingClient
	key := types.NamespacedName{Name: name}
	return WaziGetWithClient(client, key, actual)
}

// Creates object.
// Return true if a new object is created, false if it has been already created or error occurred.
func WaziCreateIfNotExists(deployContext *chetypes.DeployContext, blueprint client.Object) (isCreated bool, err error) {
	cli := getClientForObject(blueprint.GetNamespace(), deployContext)
	return WaziCreateIfNotExistsWithClient(cli, deployContext, blueprint)
}

func WaziCreateIfNotExistsWithClient(cli client.Client, deployContext *chetypes.DeployContext, blueprint client.Object) (isCreated bool, err error) {
	key := types.NamespacedName{Name: blueprint.GetName(), Namespace: blueprint.GetNamespace()}
	actual := blueprint.DeepCopyObject().(client.Object)
	exists, err := WaziGetWithClient(cli, key, actual)
	if exists {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return WaziCreateWithClient(cli, deployContext, blueprint, false)
}

// Creates object.
// Return true if a new object is created otherwise returns false.
func WaziCreate(deployContext *chetypes.DeployContext, blueprint client.Object) (bool, error) {
	client := getClientForObject(blueprint.GetNamespace(), deployContext)
	return WaziCreateWithClient(client, deployContext, blueprint, false)
}

// Deletes object.
// Returns true if object deleted or not found otherwise returns false.
func WaziDelete(deployContext *chetypes.DeployContext, key client.ObjectKey, objectMeta client.Object) (bool, error) {
	client := getClientForObject(key.Namespace, deployContext)
	return WaziDeleteByKeyWithClient(client, key, objectMeta)
}

func WaziDeleteNamespacedObject(deployContext *chetypes.DeployContext, name string, objectMeta client.Object) (bool, error) {
	client := deployContext.ClusterAPI.Client
	key := types.NamespacedName{Name: name, Namespace: deployContext.WaziLicense.Namespace}
	return WaziDeleteByKeyWithClient(client, key, objectMeta)
}

func WaziDeleteClusterObject(deployContext *chetypes.DeployContext, name string, objectMeta client.Object) (bool, error) {
	client := deployContext.ClusterAPI.NonCachingClient
	key := types.NamespacedName{Name: name}
	return WaziDeleteByKeyWithClient(client, key, objectMeta)
}

// Updates object.
// Returns true if object is up to date otherwiser return false
func WaziUpdateWithClient(client client.Client, deployContext *chetypes.DeployContext, actual client.Object, blueprint client.Object, diffOpts ...cmp.Option) (bool, error) {
	actualMeta, ok := actual.(metav1.Object)
	if !ok {
		return false, fmt.Errorf("object %T is not a metav1.Object. Cannot sync it", actualMeta)
	}

	diff := cmp.Diff(actual, blueprint, diffOpts...)
	if len(diff) > 0 {
		// don't print difference if there are no diffOpts mainly to avoid huge output
		if len(diffOpts) != 0 {
			fmt.Printf("Difference:\n%s", diff)
		}

		targetLabels := map[string]string{}
		targetAnnos := map[string]string{}

		for k, v := range actualMeta.GetAnnotations() {
			targetAnnos[k] = v
		}
		for k, v := range actualMeta.GetLabels() {
			targetLabels[k] = v
		}

		for k, v := range blueprint.GetAnnotations() {
			targetAnnos[k] = v
		}
		for k, v := range blueprint.GetLabels() {
			targetLabels[k] = v
		}

		blueprint.SetAnnotations(targetAnnos)
		blueprint.SetLabels(targetLabels)

		if isUpdateUsingDeleteCreate(actual.GetObjectKind().GroupVersionKind().Kind) {
			done, err := WaziDeleteWithClient(client, actual)
			if !done {
				return false, err
			}
			return WaziCreateWithClient(client, deployContext, blueprint, false)
		} else {
			logrus.Infof("Updating existing object: %s, name: %s", WaziGetObjectType(actualMeta), actualMeta.GetName())
			err := setOwnerReferenceIfNeeded(deployContext, blueprint)
			if err != nil {
				return false, err
			}

			// to be able to update, we need to set the resource version of the object that we know of
			blueprint.(metav1.Object).SetResourceVersion(actualMeta.GetResourceVersion())
			err = client.Update(context.TODO(), blueprint)
			return false, err
		}
	}
	return true, nil
}

func WaziCreateWithClient(client client.Client, deployContext *chetypes.DeployContext, blueprint client.Object, returnTrueIfAlreadyExists bool) (bool, error) {
	logrus.Infof("Creating a new object: %s, name: %s", WaziGetObjectType(blueprint), blueprint.GetName())

	err := setOwnerReferenceIfNeeded(deployContext, blueprint)
	if err != nil {
		return false, err
	}

	err = client.Create(context.TODO(), blueprint)
	if err == nil {
		return true, nil
	} else if errors.IsAlreadyExists(err) {
		return returnTrueIfAlreadyExists, nil
	} else {
		return false, err
	}
}

func WaziDeleteByKeyWithClient(cli client.Client, key client.ObjectKey, objectMeta client.Object) (bool, error) {
	runtimeObject, ok := objectMeta.(runtime.Object)
	if !ok {
		return false, fmt.Errorf("object %T is not a runtime.Object. Cannot sync it", runtimeObject)
	}

	actual := runtimeObject.DeepCopyObject().(client.Object)
	exists, err := WaziGetWithClient(cli, key, actual)
	if !exists {
		return true, nil
	} else if err != nil {
		return false, err
	}

	return WaziDeleteWithClient(cli, actual)
}

func WaziDeleteWithClient(client client.Client, actual client.Object) (bool, error) {
	logrus.Infof("Deleting object: %s, name: %s", WaziGetObjectType(actual), actual.GetName())

	err := client.Delete(context.TODO(), actual)
	if err == nil || errors.IsNotFound(err) {
		return true, nil
	} else {
		return false, err
	}
}

func WaziGetWithClient(client client.Client, key client.ObjectKey, object client.Object) (bool, error) {
	err := client.Get(context.TODO(), key, object)
	if err == nil {
		return true, nil
	} else if errors.IsNotFound(err) {
		return false, nil
	} else {
		return false, err
	}
}

func isUpdateUsingDeleteCreate(kind string) bool {
	return kind == "Service" || kind == "Ingress" || kind == "Route" || kind == "Job" || kind == "Secret"
}

func setOwnerReferenceIfNeeded(deployContext *chetypes.DeployContext, blueprint metav1.Object) error {
	if shouldSetOwnerReferenceForObject(deployContext, blueprint) {
		return controllerutil.SetControllerReference(deployContext.WaziLicense, blueprint, deployContext.ClusterAPI.Scheme)
	}

	return nil
}

func shouldSetOwnerReferenceForObject(deployContext *chetypes.DeployContext, blueprint metav1.Object) bool {
	// empty workspace (cluster scope object) or object in another namespace
	return blueprint.GetNamespace() == deployContext.WaziLicense.Namespace
}

func getClientForObject(objectNamespace string, deployContext *chetypes.DeployContext) client.Client {
	// empty namespace (cluster scope object) or object in another namespace
	if deployContext.WaziLicense.Namespace == objectNamespace {
		return deployContext.ClusterAPI.Client
	}
	return deployContext.ClusterAPI.NonCachingClient
}

func WaziGetObjectType(obj interface{}) string {
	objType := reflect.TypeOf(obj).String()
	if reflect.TypeOf(obj).Kind().String() == "ptr" {
		objType = objType[1:]
	}

	return objType
}