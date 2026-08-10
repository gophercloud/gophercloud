package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const (
	UnscopedTokenID        = "unscoped-federation-token-aaa111"
	FederationAuthResponse = `{
		"token": {
			"methods": ["mapped"],
			"user": {
				"id": "federation-user-id-001",
				"name": "federation-user",
				"domain": {"id": "Federated", "name": "Federated"}
			},
			"expires_at": "2035-06-03T02:19:49Z",
			"is_domain": false
		}
	}`
)

func HandleWebSSOTokenValidation(t *testing.T, fakeServer th.FakeServer, expectedToken string) {
	t.Helper()
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Subject-Token", expectedToken)
		w.Header().Set("X-Subject-Token", expectedToken)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, FederationAuthResponse)
	})
}
