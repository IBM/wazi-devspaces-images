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
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

func WaziGetLabels(component string) map[string]string {

	cheFlavor := os.Getenv("CHE_FLAVOR")
	if len(cheFlavor) == 0 {
		err := errors.New("Environment variable for CHE_FLAVOR not set.")
		logrus.Fatalf(err.Error())
		return map[string]string{}
	}

	return map[string]string{
		KubernetesNameLabelKey:      cheFlavor,
		KubernetesInstanceLabelKey:  cheFlavor,
		KubernetesPartOfLabelKey:    CheEclipseOrg,
		KubernetesComponentLabelKey: component,
		KubernetesManagedByLabelKey: cheFlavor + "-operator",
	}
}
