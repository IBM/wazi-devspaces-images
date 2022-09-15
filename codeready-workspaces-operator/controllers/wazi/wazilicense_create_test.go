package wazi

import (
	"context"
	"reflect"
	"testing"

	orgv1 "github.com/eclipse-che/che-operator/api/v1"
	"github.com/eclipse-che/che-operator/pkg/deploy"
	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	testObjectMeta = v1.ObjectMeta{
		Name:      "sample-licenseSpecCreate-name",
		Namespace: "sample-licenseSpecCreate-namespace",
	}
)

func TestGenerateandSaveFields(t *testing.T) {
	type testCase struct {
		name                  string
		testCR                *orgv1.WaziLicense
		expectedType          string
		expectedLicenseAccept bool
	}
	testCases := []testCase{
		{
			name:                  "Check WaziLicenseSpec when type is `Wazi` and accept value `False`",
			testCR:                returnGenerateAndSaveFieldsWaziCR("wazi", false),
			expectedType:          "wazi",
			expectedLicenseAccept: false,
		},
		{
			name:                  "Check WaziLicenseSpec when type is `idzee` and accept value `False` ",
			testCR:                returnGenerateAndSaveFieldsWaziCR("idzee", false),
			expectedType:          "idzee",
			expectedLicenseAccept: false,
		},
		{
			name:                  "Check WaziLicenseSpec when type is not specified and accept value `False`",
			testCR:                returnGenerateAndSaveFieldsWaziCR("", false),
			expectedType:          "wazi",
			expectedLicenseAccept: false,
		},
		{
			name:                  "Check WaziLicenseSpec when type is not specidied and accept value `True`",
			testCR:                returnGenerateAndSaveFieldsWaziCR("", true),
			expectedType:          "wazi",
			expectedLicenseAccept: true,
		},
	}
	for _, testCase := range testCases {
		objs, dc, scheme := createAPIObjects()
		objs = append(objs, testCase.testCR)
		cl := fake.NewFakeClient(objs...)
		r := NewReconciler(cl, cl, dc, &scheme, "")
		clusterAPI := deploy.ClusterAPI{
			Client:           r.client,
			NonCachingClient: r.client,
			Scheme:           r.Scheme,
		}

		deployContext := deploy.DeployContext{
			WaziLicense: testCase.testCR,
			ClusterAPI:  clusterAPI,
		}
		if err := r.GenerateAndSaveFields(&deployContext); err != nil {
			t.Fatal("GenerateAndSaveFields: err =", err)
		}
		returnedWaziCR := orgv1.WaziLicense{}
		if err := deployContext.ClusterAPI.Client.Get(context.TODO(), types.NamespacedName{Namespace: deployContext.WaziLicense.Namespace, Name: deployContext.WaziLicense.Name}, &returnedWaziCR); err != nil {
			t.Fatalf("Unable to retrieve waziLicense from the cluster API: err=%v", err)
		}
		if !reflect.DeepEqual(returnedWaziCR.Spec.License.Accept, testCase.expectedLicenseAccept) {
			t.Errorf("Expected WaziLicense.Spec.License.Accept and got WaziLicense.Spec.License.Accept different (-want +got): %v", cmp.Diff(testCase.expectedLicenseAccept, returnedWaziCR.Spec.License.Accept))
		}
		if !reflect.DeepEqual(returnedWaziCR.Spec.License.Use, testCase.expectedType) {
			t.Errorf("Expected WaziLicense.Spec.License.Use and and returned WaziLicense.Spec.License.Accept are different (-want +got): %v", cmp.Diff(testCase.expectedType, returnedWaziCR.Spec.License.Use))
		}
	}
}

func returnGenerateAndSaveFieldsWaziCR(licenseUse string, license bool) *orgv1.WaziLicense {
	if licenseUse != "" {
		return &orgv1.WaziLicense{
			ObjectMeta: testObjectMeta,
			Spec: orgv1.WaziLicenseSpec{
				License: orgv1.WaziLicenseSpecLicense{
					Use: licenseUse,
				},
			},
		}
	}
	if license {
		return &orgv1.WaziLicense{
			ObjectMeta: testObjectMeta,
			Spec: orgv1.WaziLicenseSpec{
				License: orgv1.WaziLicenseSpecLicense{
					Accept: true,
				},
			},
		}
	}
	return &orgv1.WaziLicense{
		ObjectMeta: testObjectMeta,
	}
}
