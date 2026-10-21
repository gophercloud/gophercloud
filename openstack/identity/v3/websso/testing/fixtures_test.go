package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const unscopedTokenID = "unscoped-federation-token-aaa111"

func federationAuthResponse() string {
	return fmt.Sprintf(`{
		"token": {
			"methods": ["mapped"],
			"user": {
				"id": "federation-user-id-001",
				"name": "federation-user",
				"domain": {"id": "Federated", "name": "Federated"}
			},
			"expires_at": %q,
			"is_domain": false
		}
	}`, time.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano))
}

func handleWebSSOTokenValidation(t *testing.T, fakeServer th.FakeServer, expectedToken string) {
	t.Helper()
	fakeServer.Mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Subject-Token", expectedToken)
		w.Header().Set("X-Subject-Token", expectedToken)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, federationAuthResponse())
	})
}
