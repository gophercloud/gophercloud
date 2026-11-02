package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/drivers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Ironic returns the next page as a top-level "next" URL.
func TestDriverListFollowsNextLink(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/drivers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		switch r.URL.Query().Get("marker") {
		case "":
			fmt.Fprintf(w, `{"drivers": [{"name": "first"}], "next": "%s/drivers?limit=1&marker=first"}`, fakeServer.Server.URL)
		case "first":
			fmt.Fprint(w, `{"drivers": [{"name": "second"}]}`)
		default:
			t.Fatalf("unexpected marker %q", r.URL.Query().Get("marker"))
		}
	})

	pages, items := 0, 0
	err := drivers.ListDrivers(client.ServiceClient(fakeServer), drivers.ListDriversOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		actual, err := drivers.ExtractDrivers(page)
		items += len(actual)
		return true, err
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, pages)
	th.AssertEquals(t, 2, items)
}
