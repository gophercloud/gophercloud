package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/conductors"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Ironic returns the next page as a top-level "next" URL.
func TestConductorListFollowsNextLink(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/conductors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		switch r.URL.Query().Get("marker") {
		case "":
			fmt.Fprintf(w, `{"conductors": [{"hostname": "first"}], "next": "%s/conductors?limit=1&marker=first"}`, fakeServer.Server.URL)
		case "first":
			fmt.Fprint(w, `{"conductors": [{"hostname": "second"}]}`)
		default:
			t.Fatalf("unexpected marker %q", r.URL.Query().Get("marker"))
		}
	})

	pages, items := 0, 0
	err := conductors.List(client.ServiceClient(fakeServer), conductors.ListOpts{Limit: 1}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		actual, err := conductors.ExtractConductors(page)
		items += len(actual)
		return true, err
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.AssertEquals(t, 2, items)
}
