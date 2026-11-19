package testing

import (
	"context"
	"errors"
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
		for _, test := range []struct {
			name     string
			opts     schedulerstats.ListOptsBuilder
			expected url.Values
		}{
			{
				name:     "nil",
				expected: url.Values{},
			},
			{
				name:     "empty",
				opts:     schedulerstats.ListOpts{},
				expected: url.Values{},
			},
			{
				name: "filters",
				opts: schedulerstats.ListOpts{
					Pool:      "^pool#1$",
					Host:      `^host\.example\.com$`,
					Backend:   "backend@1",
					ShareType: "share type+1",
				},
				expected: url.Values{
					"pool":       {"^pool#1$"},
					"host":       {`^host\.example\.com$`},
					"backend":    {"backend@1"},
					"share_type": {"share type+1"},
				},
			},
		} {
			t.Run(fmt.Sprintf("detail=%t/%s", detail, test.name), func(t *testing.T) {
				fakeServer := th.SetupHTTP()
				defer fakeServer.Teardown()

				path := "/scheduler-stats/pools"
				if detail {
					path += "/detail"
				}
				fakeServer.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					th.TestMethod(t, r, "GET")
					th.CheckDeepEquals(t, test.expected, r.URL.Query())
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, PoolsListBody)
				})

				service := client.ServiceClient(fakeServer)
				var pager pagination.Pager
				if detail {
					pager = schedulerstats.ListDetail(service, test.opts)
				} else {
					pager = schedulerstats.List(service, test.opts)
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

type errorListOpts struct {
	err error
}

func (opts errorListOpts) ToPoolsListQuery() (string, error) {
	return "", opts.err
}

func TestListPoolQueryError(t *testing.T) {
	for _, detail := range []bool{false, true} {
		t.Run(fmt.Sprintf("detail=%t", detail), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				t.Error("Unexpected HTTP request after query builder error")
			})

			expectedErr := errors.New("cannot build pool query")
			opts := errorListOpts{err: expectedErr}
			service := client.ServiceClient(fakeServer)
			var pager pagination.Pager
			if detail {
				pager = schedulerstats.ListDetail(service, opts)
			} else {
				pager = schedulerstats.List(service, opts)
			}
			_, err := pager.AllPages(context.TODO())
			th.CheckEquals(t, expectedErr, err)
		})
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
	err = schedulerstats.ListDetail(client.ServiceClient(fakeServer), schedulerstats.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
