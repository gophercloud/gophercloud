package v3

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

func TestAuthentication(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", "aaaa1111")

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	provider := &gophercloud.ProviderClient{
		IdentityEndpoint: testhelper.Endpoint(),
	}
	client := NewClient(provider)

	expected := testhelper.Endpoint() + "v3/"
	if client.Endpoint != expected {
		t.Errorf("Expected endpoint to be %s, but was %s", expected, client.Endpoint)
	}
}
