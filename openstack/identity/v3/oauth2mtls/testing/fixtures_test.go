package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const (
	tokenID       = "gAAAAABl-mTLS-token-abc123"
	testClientID  = "6c3145f4-313d-4910-b3a8-9dfc72da9e75"
	tokenResponse = `{
		"access_token": "gAAAAABl-mTLS-token-abc123",
		"token_type": "Bearer",
		"expires_in": 3600
	}`
	validateTokenResponse = `{
		"token": {
			"catalog": [{
				"type": "compute",
				"name": "nova",
				"endpoints": [{
					"interface": "public",
					"region": "iad1",
					"url": "http://127.0.0.1:8774/v2.1"
				}]
			}],
			"expires_at": "2030-06-15T18:00:00.000000Z"
		}
	}`
)

func handleTokenSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/OS-OAUTH2/token", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeader(t, r, "Content-Type", "application/x-www-form-urlencoded")
		th.TestHeader(t, r, "Accept", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		th.AssertEquals(t, "client_credentials", r.PostForm.Get("grant_type"))
		th.AssertEquals(t, testClientID, r.PostForm.Get("client_id"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, tokenResponse)
	})
}

func handleValidateTokenSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		th.TestHeader(t, r, "X-Subject-Token", tokenID)
		w.Header().Set("X-Subject-Token", tokenID)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, validateTokenResponse)
	})
}
