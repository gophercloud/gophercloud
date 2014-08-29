package v3

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"github.com/rackspace/gophercloud/testhelper"
)

func TestNewClient(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	provider := &gophercloud.ProviderClient{
		IdentityEndpoint: testhelper.Endpoint() + "v3/",
	}
	client := NewClient(provider)

	expected := testhelper.Endpoint() + "v3/"
	if client.Endpoint != expected {
		t.Errorf("Expected endpoint to be %s, but was %s", expected, client.Endpoint)
	}
}

func TestAuthentication(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()
	const ID = "aaaa1111"

	testhelper.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	provider := &gophercloud.ProviderClient{
		IdentityEndpoint: testhelper.Endpoint() + "v3/",
	}
	client := NewClient(provider)

	token, err := client.Authenticate(gophercloud.AuthOptions{UserID: "me", Password: "swordfish"})
	if err != nil {
		t.Errorf("Unexpected error from authentication: %v", err)
	}

	if token.ID != ID {
		t.Errorf("Expected token ID [%s], but got [%s]", ID, token.ID)
	}

	expectedExpiration, _ := time.Parse(tokens.RFC3339Milli, "2013-02-02T18:30:59.000000Z")
	if token.ExpiresAt != expectedExpiration {
		t.Errorf("Expected token expiration [%v], but got [%v]", expectedExpiration, token.ExpiresAt)
	}
}
