package testing

import (
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// HandleDiagnosticGetSuccessfully sets up the test server to respond to a diagnostic Get request.
func HandleDiagnosticGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/servers/1234asdf/diagnostics", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"cpu0_time":173,"memory":524288}`))
		th.AssertNoErr(t, err)
	})
}
