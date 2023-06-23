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
	"github.com/eclipse-che/che-operator/pkg/common/chetypes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconcilable interface {
	Reconcile(ctx *chetypes.DeployContext) (result reconcile.Result, done bool, err error)
	Finalize(ctx *chetypes.DeployContext) (done bool)
}

type WaziReconcileManager struct {
	reconcilers      []Reconcilable
	failedReconciler Reconcilable
}

func NewReconcileManager() *WaziReconcileManager {
	return &WaziReconcileManager{
		reconcilers:      make([]Reconcilable, 0),
		failedReconciler: nil,
	}
}

func (manager *WaziReconcileManager) RegisterReconciler(reconciler Reconcilable) {
	manager.reconcilers = append(manager.reconcilers, reconciler)
}

func (manager *WaziReconcileManager) ReconcileAll(ctx *chetypes.DeployContext) (reconcile.Result, bool, error) {
	for _, reconciler := range manager.reconcilers {
		result, done, err := reconciler.Reconcile(ctx)
		if err != nil {
			manager.failedReconciler = reconciler
		} else if manager.failedReconciler == reconciler {
			manager.failedReconciler = nil
		}

		if !done {
			return result, done, err
		}
	}

	return reconcile.Result{}, true, nil
}

func (manager *WaziReconcileManager) FinalizeAll(ctx *chetypes.DeployContext) (done bool) {
	done = true
	for _, reconciler := range manager.reconcilers {
		if completed := reconciler.Finalize(ctx); !completed {
			done = false
		}
	}
	return done
}
