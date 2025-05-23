package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func MockListResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-availability-zone", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
        {
            "availability_zones": [
                {
                    "name": "nova",
                    "created_at": "2015-09-18T09:50:55.000000",
                    "updated_at": null,
                    "id": "388c983d-258e-4a0e-b1ba-10da37d766db"
                }
            ]
        }`)
	})
}
