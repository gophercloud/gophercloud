package stacks

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

// CreateExpected represents the expected object from a Create request.
var CreateExpected = &CreatedStack{
	ID: "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
	Links: []gophercloud.Link{
		gophercloud.Link{
			Href: "http://168.28.170.117:8004/v1/98606384f58drad0bhdb7d02779549ac/stacks/stackcreated/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
			Rel:  "self",
		},
	},
}

// CreateOutput represents the response body from a Create request.
const CreateOutput = `
{
  "stack": {
    "id": "16ef0584-4458-41eb-87c8-0dc8d5f66c87",
    "links": [
    {
      "href": "http://168.28.170.117:8004/v1/98606384f58drad0bhdb7d02779549ac/stacks/stackcreated/16ef0584-4458-41eb-87c8-0dc8d5f66c87",
      "rel": "self"
    }
    ]
  }
}`

// HandleCreateSuccessfully creates an HTTP handler at `/stacks` on the test handler mux
// that responds with a `Create` response.
func HandleCreateSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, output)
	})
}
