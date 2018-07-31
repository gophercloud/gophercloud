package testing

import (
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

// HandleServerRescueSuccessfully sets up the test server to respond to a server Rescue request.
func HandleServerRescueSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/servers/1234asdf/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `{ "rescue": { "adminPass": "1234567890" } }`)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "adminPass": "1234567890" }`))
	})
}
