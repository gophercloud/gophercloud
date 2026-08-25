package config

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	tokenstesting "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens/testing"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestNewProviderClientV3(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		w.Header().Set("X-Subject-Token", "token")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, tokenstesting.TokenOutput)
	})

	client, err := NewProviderClientV3(
		context.Background(),
		fakeServer.Endpoint()+"v3/",
		&tokens.AuthOptions{
			UserID:   "user-id",
			Password: "password",
		},
	)
	th.AssertNoErr(t, err)
	if client.Token() != "token" {
		t.Fatalf("unexpected token: %q", client.Token())
	}
}
