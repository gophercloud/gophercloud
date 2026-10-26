package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func keyPairPageRecord(index int) map[string]any {
	return map[string]any{"keypair": map[string]any{"name": strconv.Itoa(index)}}
}

func TestKeyPairPagePagination(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				th.TestMethod(t, r, "GET")
				th.TestHeader(t, r, "OpenStack-API-Version", "compute 2.35")
				if requests > 2 {
					http.Error(w, "pagination did not terminate", http.StatusInternalServerError)
					return
				}
				index := 1
				if r.URL.Query().Get("page") == "2" {
					index = 2
					th.CheckEquals(t, "kept", r.URL.Query().Get("filter"))
				}
				items := []any{}
				if index <= 2 {
					items = append(items, keyPairPageRecord(index))
				}
				body := map[string]any{"keypairs": items, "unrelated": []any{"not a resource"}}
				if index == 1 {
					body["keypairs_links"] = []map[string]string{{"rel": "self", "href": fakeServer.Endpoint()}, {"rel": "next", "href": fakeServer.Endpoint() + r.URL.Path[1:] + "?page=2&filter=kept"}}
				}

				w.Header().Set("Content-Type", "application/json")
				th.AssertNoErr(t, json.NewEncoder(w).Encode(body))
			})
			c := client.ServiceClient(fakeServer)
			c.Type = "compute"
			c.Microversion = "2.35"
			pager := keypairs.List(c, nil)
			ids := []string{}
			if allPages {
				page, err := pager.AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := keypairs.ExtractKeyPairs(page)
				th.AssertNoErr(t, err)
				for _, item := range items {
					ids = append(ids, fmt.Sprint(item.Name))
				}
			} else {
				pages := 0
				err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					items, err := keypairs.ExtractKeyPairs(page)
					th.CheckEquals(t, 1, len(items))
					for _, item := range items {
						ids = append(ids, fmt.Sprint(item.Name))
					}
					return true, err
				})
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 2, pages)
			}
			th.CheckDeepEquals(t, []string{"1", "2"}, ids)
			th.CheckEquals(t, 2, requests)
		})
	}
}

func TestKeyPairPagePaginationEmpty(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNoContent} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				if requests > 1 {
					http.Error(w, "empty page did not terminate pagination", http.StatusInternalServerError)
					return
				}
				if status == http.StatusNoContent {
					w.WriteHeader(status)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"keypairs":[]}`)
			})
			c := client.ServiceClient(fakeServer)
			err := keypairs.List(c, nil).EachPage(context.Background(), func(context.Context, pagination.Page) (bool, error) {
				t.Error("empty page reached callback")
				return false, nil
			})
			th.AssertNoErr(t, err)
			th.CheckEquals(t, 1, requests)
			if status == http.StatusOK {
				requests = 0
				page, err := keypairs.List(c, nil).AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := keypairs.ExtractKeyPairs(page)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 0, len(items))
				th.CheckEquals(t, 1, requests)
			}
		})
	}
}

func TestKeyPairPageNextPageURL(t *testing.T) {
	u, err := url.Parse("https://example.com/resources?offset=7&limit=5&filter=kept")
	th.AssertNoErr(t, err)
	result := pagination.PageResult{URL: *u, Result: gophercloud.Result{StatusCode: 200, Body: map[string]any{"keypairs": []any{keyPairPageRecord(1), keyPairPageRecord(2)}}}}
	result.Body.(map[string]any)["keypairs_links"] = []map[string]string{{"rel": "next", "href": "https://example.com/next"}}
	page := keypairs.KeyPairPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err := page.NextPageURL("https://example.com/")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "https://example.com/next", next)

	result.URL = *u
	result.Body = map[string]any{"keypairs": "invalid"}
	result.Body.(map[string]any)["keypairs_links"] = "invalid"
	page = keypairs.KeyPairPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	_, err = page.NextPageURL("")
	if err == nil {
		t.Fatal("expected malformed response error")
	}
	result.StatusCode = http.StatusNoContent
	result.Body = nil
	page = keypairs.KeyPairPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err = page.NextPageURL("")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "", next)
}

func TestKeyPairPageLegacyPagination(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				th.TestHeader(t, r, "OpenStack-API-Version", "compute 2.34")
				if requests > 1 {
					http.Error(w, "unexpected continuation", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				th.AssertNoErr(t, json.NewEncoder(w).Encode(map[string]any{"keypairs": []any{keyPairPageRecord(1)}}))
			})
			c := client.ServiceClient(fakeServer)
			c.Type = "compute"
			c.Microversion = "2.34"
			pager := keypairs.List(c, nil)
			if allPages {
				page, err := pager.AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := keypairs.ExtractKeyPairs(page)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, 1, len(items))
				th.CheckEquals(t, "1", fmt.Sprint(items[0].Name))
			} else {
				pages := 0
				err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					items, err := keypairs.ExtractKeyPairs(page)
					th.CheckEquals(t, 1, len(items))
					return true, err
				})
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 1, pages)
			}
			th.CheckEquals(t, 1, requests)
		})
	}
}
