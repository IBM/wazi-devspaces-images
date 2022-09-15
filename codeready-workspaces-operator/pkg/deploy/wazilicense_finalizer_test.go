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
	"github.com/eclipse-che/che-operator/pkg/util"

	"testing"
)

const (
	wazi_finalizer = "some.finalizer"
)

func TestWaziAppendFinalizer(t *testing.T) {
	_, deployContext := initWaziDeployContext()

	err := WaziAppendFinalizer(deployContext, wazi_finalizer)
	if err != nil {
		t.Fatalf("Failed to append finalizer: %v", err)
	}

	if !util.ContainsString(deployContext.WaziLicense.ObjectMeta.Finalizers, wazi_finalizer) {
		t.Fatalf("Failed to append finalizer: %v", err)
	}

	// shouldn't add finalizer twice
	err = WaziAppendFinalizer(deployContext, wazi_finalizer)
	if err != nil {
		t.Fatalf("Failed to append finalizer: %v", err)
	}

	if len(deployContext.WaziLicense.ObjectMeta.Finalizers) != 1 {
		t.Fatalf("Finalizer shouldn't be added twice")
	}
}

func TestWaziDeleteFinalizer(t *testing.T) {
	_, deployContext := initWaziDeployContext()
	err := WaziAppendFinalizer(deployContext, wazi_finalizer)

	if err != nil {
		t.Fatalf("Failed to append finalizer when trying to delete: %v", err)
	}

	err = WaziDeleteFinalizer(deployContext, wazi_finalizer)
	if err != nil {
		t.Fatalf("Failed to append finalizer: %v", err)
	}

	if util.ContainsString(deployContext.WaziLicense.ObjectMeta.Finalizers, wazi_finalizer) {
		t.Fatalf("Failed to delete finalizer: %v", err)
	}
}

func TestWaziDeleteObjectWithFinalizer(t *testing.T) {
	_, deployContext := initWaziDeployContext()

	err := WaziDeleteObjectWithFinalizer(deployContext, testWaziKey, testWaziObj.DeepCopy(), wazi_finalizer)

	if err != nil {
		t.Fatalf("Failed to delete finalizer when deleting object with finalizer: %v", err)
	}

	if util.ContainsString(deployContext.WaziLicense.ObjectMeta.Finalizers, wazi_finalizer) {
		t.Fatalf("Failed to delete finalizer when deleting object with finalizer: %v", err)
	}
}

func TestWaziGetFinalizerNameShouldReturnStringLess64Chars(t *testing.T) {
	expected := "7890123456789012345678901234567891234567.finalizers.che.eclipse"
	prefix := "7890123456789012345678901234567891234567"

	actual := WaziGetFinalizerName(prefix)
	if expected != actual {
		t.Fatalf("Incorrect finalizer name: %s", actual)
	}
}
