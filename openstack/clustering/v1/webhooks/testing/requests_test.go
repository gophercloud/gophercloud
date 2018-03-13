package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestWebhooks(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
	    "services": [
	        {
						"binary": "senlin-engine",
						"disabled_reason": null,
						"host": "host1",
						"id": "f93f83f6-762b-41b6-b757-80507834d394",
						"state": "up",
						"status": "enabled",
						"topic": "senlin-engine",
						"updated_at": "2017-04-24T07:43:12",
					},
	    ]
			}`)
	})

	// TODO: Need to implement
}
