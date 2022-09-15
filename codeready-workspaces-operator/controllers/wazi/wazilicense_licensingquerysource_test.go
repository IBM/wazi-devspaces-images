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

package wazi

import (
	"context"
	"reflect"
	"testing"

	ibmlicensingv1 "github.com/IBM/ibm-licensing-operator/api/v1"
	orgv1 "github.com/eclipse-che/che-operator/api/v1"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/eclipse-che/che-operator/pkg/util"
	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestReconcileWaziLicensingQuerySource(t *testing.T) {
	type testCase struct {
		name                         string
		testLicenseQuerySourceExists bool
	}

	testCases := []testCase{
		{
			name:                         "Condition: Check for query and annotations if IBMLicensingQuerySource is absent",
			testLicenseQuerySourceExists: false,
		},
		{
			name:                         "Condition: Check for query and annotations if IBMLicensingQuerySource is existent",
			testLicenseQuerySourceExists: true,
		},
	}
	for _, testCases := range testCases {
		objs, dc, scheme := createAPIObjects()
		if testCases.testLicenseQuerySourceExists {
			testIBMLicenseQuerySource := initWaziWithIBMLicensingQuerySource(util.WaziLicensingQuerySourceName)
			objs = append(objs, testIBMLicenseQuerySource)
		}
		cl := fake.NewFakeClient(objs...)
		r := NewReconciler(cl, cl, dc, &scheme, "")
		clusterAPI := deploy.ClusterAPI{
			Client:           r.client,
			NonCachingClient: r.client,
			Scheme:           r.Scheme,
			DiscoveryClient:  r.discoveryClient,
		}

		waziCR := &orgv1.WaziLicense{
			ObjectMeta: v1.ObjectMeta{
				Name:      waziLicenseName,
				Namespace: namespace,
			},
		}
		deployContext := deploy.DeployContext{
			WaziLicense: waziCR,
			ClusterAPI:  clusterAPI,
		}
		if _, err := r.ReconcileWaziLicensingQuerySource(&deployContext); err != nil {
			t.Fatal(err)
		}
		returnedIBMLicenseQueryCr := ibmlicensingv1.IBMLicensingQuerySource{}
		if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: "eclipse-che", Name: util.WaziLicensingQuerySourceName}, &returnedIBMLicenseQueryCr); err != nil {
			t.Fatalf("Could not fetch IBMLicensingQuerySource from cluster API: err=%v", err)
		}

		if !reflect.DeepEqual(returnedIBMLicenseQueryCr.Spec, getWaziLicensingQuerySourceSpec(&deployContext, util.WaziLicensingQuerySourceName).Spec) {
			t.Fatalf("Returned IBMLicensingQuerySourceSpec and IBMLicensingQuerySourceSpec expected from API server for the given IBMLicensingQuerySource are different (-want +got): %v", cmp.Diff(returnedIBMLicenseQueryCr.Spec, getWaziLicensingQuerySourceSpec(&deployContext, util.WaziLicensingQuerySourceName).Spec))
		}

		lqs, err := getWaziLicensingQuerySource(&deployContext, util.WaziLicensingQuerySourceName)
		if err != nil {
			t.Fatalf("error when getting wazi licesing query source: %v", err)
		}
		if lqs == nil {
			t.Fatalf("Wazi licesing query source is nil")
		}

	}
}
