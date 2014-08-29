package openstack

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

func TestNewClientV3(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	const ID = "0123456789"

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": {
					"values": [
						{
							"status": "stable",
							"id": "v3.0",
							"links": [
								{ "href": "%s", "rel": "self" }
							]
						},
						{
							"status": "stable",
							"id": "v2.0",
							"links": [
								{ "href": "%s", "rel": "self" }
							]
						}
					]
				}
			}
		`, testhelper.Endpoint()+"v3/", testhelper.Endpoint()+"v2.0/")
	})

	testhelper.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	options := gophercloud.AuthOptions{
		UserID:           "me",
		Password:         "secret",
		IdentityEndpoint: testhelper.Endpoint(),
	}
	client, err := NewClient(options)

	if err != nil {
		t.Fatalf("Unexpected error from NewClient: %s", err)
	}

	if client.TokenID != ID {
		t.Errorf("Expected token ID to be [%s], but was [%s]", ID, client.TokenID)
	}
}
