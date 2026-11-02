package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Ironic returns the next page as a top-level "next" URL.
func TestPortListFollowsNextLink(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/ports", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		switch r.URL.Query().Get("marker") {
		case "":
			fmt.Fprintf(w, `{"ports": [{"uuid": "first"}], "next": "%s/ports?limit=1&marker=first"}`, fakeServer.Server.URL)
		case "first":
			fmt.Fprint(w, `{"ports": [{"uuid": "second"}]}`)
		default:
			t.Fatalf("unexpected marker %q", r.URL.Query().Get("marker"))
		}
	})

	pages, items := 0, 0
	err := ports.List(client.ServiceClient(fakeServer), ports.ListOpts{Limit: 1}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		actual, err := ports.ExtractPorts(page)
		items += len(actual)
		return true, err
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.AssertEquals(t, 2, items)
}
