package accounts

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var metadata = map[string]string{"gophercloud-test": "accounts"}

func TestUpdateAccount(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		testhelper.TestHeader(t, r, "X-Account-Meta-Gophercloud-Test", "accounts")
		w.WriteHeader(http.StatusNoContent)
	})

	err := Update(fake.ServiceClient(), UpdateOpts{Metadata: metadata})
	if err != nil {
		t.Fatalf("Unable to update account: %v", err)
	}
}

func TestGetAccount(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "HEAD")
		testhelper.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := Get(fake.ServiceClient(), GetOpts{})
	if err != nil {
		t.Fatalf("Unable to get account metadata: %v", err)
	}
}

func TestExtractAccountMetadata(t *testing.T) {
	getResult := &http.Response{}

	expected := map[string]string{}

	actual := ExtractMetadata(getResult)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
	}
}
