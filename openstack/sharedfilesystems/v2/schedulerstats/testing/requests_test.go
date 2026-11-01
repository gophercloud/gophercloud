package testing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/schedulerstats"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListPoolFilters(t *testing.T) {
	for _, detail := range []bool{false, true} {
		for _, filtered := range []bool{false, true} {
			t.Run(fmt.Sprintf("detail=%t/filtered=%t", detail, filtered), func(t *testing.T) {
				fakeServer := th.SetupHTTP()
				defer fakeServer.Teardown()

				opts := schedulerstats.ListOpts{ProjectID: "project-in-endpoint"}
				expected := url.Values{}
				if filtered {
					opts.PoolName = "pool#1"
					opts.HostName = "host.example.com"
					opts.BackendName = "backend@1"
					opts.Capabilities = "capability value"
					opts.ShareType = "share type"
					expected = url.Values{
						"pool": {opts.PoolName}, "host": {opts.HostName},
						"backend": {opts.BackendName}, "capabilities": {opts.Capabilities},
						"share_type": {opts.ShareType},
					}
				}
				path := "/scheduler-stats/pools"
				if detail {
					path += "/detail"
				}
				fakeServer.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					th.TestMethod(t, r, "GET")
					th.CheckDeepEquals(t, expected, r.URL.Query())
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, PoolsListBody)
				})

				service := client.ServiceClient(fakeServer)
				var pager pagination.Pager
				if detail {
					pager = schedulerstats.ListDetail(service, schedulerstats.ListDetailOpts(opts))
				} else {
					pager = schedulerstats.List(service, opts)
				}
				pages, err := pager.AllPages(context.TODO())
				th.AssertNoErr(t, err)
				pools, err := schedulerstats.ExtractPools(pages)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 4, len(pools))
			})
		}
	}
}

func TestListPoolsDetail(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePoolsListSuccessfully(t, fakeServer)

	pages := 0
	err := schedulerstats.List(client.ServiceClient(fakeServer), schedulerstats.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := schedulerstats.ExtractPools(page)
		th.AssertNoErr(t, err)

		if len(actual) != 4 {
			t.Fatalf("Expected 4 backends, got %d", len(actual))
		}
		th.CheckDeepEquals(t, PoolFake1, actual[0])
		th.CheckDeepEquals(t, PoolFake2, actual[1])
		th.CheckDeepEquals(t, PoolFake3, actual[2])
		th.CheckDeepEquals(t, PoolFake4, actual[3])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}

	pages = 0
	err = schedulerstats.ListDetail(client.ServiceClient(fakeServer), schedulerstats.ListDetailOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := schedulerstats.ExtractPools(page)
		th.AssertNoErr(t, err)

		if len(actual) != 4 {
			t.Fatalf("Expected 4 backends, got %d", len(actual))
		}
		th.CheckDeepEquals(t, PoolDetailFake1, actual[0])
		th.CheckDeepEquals(t, PoolDetailFake2, actual[1])
		th.CheckDeepEquals(t, PoolDetailFake3, actual[2])
		th.CheckDeepEquals(t, PoolDetailFake4, actual[3])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
