package accounts

import (
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

const tokenId = "abcabcabcabc"

var metadata = map[string]string{"gophercloud-test": "accounts"}

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenId},
		Endpoint: testhelper.Endpoint(),
	}
}

func TestUpdateAccount(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		testhelper.TestHeader(t, r, "X-Account-Meta-Gophercloud-Test", "accounts")
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	err := Update(client, UpdateOpts{Metadata: metadata})
	if err != nil {
		t.Fatalf("Unable to update account: %v", err)
	}
}

func TestGetAccount(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "HEAD")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenId)
		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()
	_, err := Get(client, GetOpts{})
	if err != nil {
		t.Fatalf("Unable to get account metadata: %v", err)
	}
}
