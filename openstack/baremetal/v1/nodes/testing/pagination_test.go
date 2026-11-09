package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Ironic returns the next page as a top-level "next" URL.
func TestNodeListFollowsNextLink(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		switch r.URL.Query().Get("marker") {
		case "":
			fmt.Fprintf(w, `{"nodes": [{"uuid": "first"}], "next": "%s/nodes?limit=1&marker=first"}`, fakeServer.Server.URL)
		case "first":
			fmt.Fprint(w, `{"nodes": [{"uuid": "second"}]}`)
		default:
			t.Fatalf("unexpected marker %q", r.URL.Query().Get("marker"))
		}
	})

	pages, items := 0, 0
	err := nodes.List(client.ServiceClient(fakeServer), nodes.ListOpts{Limit: 1}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		actual, err := nodes.ExtractNodes(page)
		items += len(actual)
		return true, err
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.AssertEquals(t, 2, items)
}
