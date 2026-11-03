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
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/messages"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func messagePageRecord(index int) map[string]any { return map[string]any{"id": strconv.Itoa(index)} }

func TestMessagePagePagination(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				th.TestMethod(t, r, "GET")
				th.TestHeader(t, r, "OpenStack-API-Version", "sharev2 2.37")
				if requests > 3 {
					http.Error(w, "pagination did not terminate", http.StatusInternalServerError)
					return
				}
				index := 1
				if value := r.URL.Query().Get("offset"); value != "" {
					position, err := strconv.Atoi(value)
					th.AssertNoErr(t, err)
					index = position + 1
				}
				items := []any{}
				if index <= 2 {
					items = append(items, messagePageRecord(index))
				}
				body := map[string]any{"messages": items, "unrelated": []any{"not a resource"}}

				body["messages_links"] = []map[string]string{{"rel": "next", "href": fakeServer.Endpoint() + "messages?marker=ignored"}}
				w.Header().Set("Content-Type", "application/json")
				th.AssertNoErr(t, json.NewEncoder(w).Encode(body))
			})
			c := client.ServiceClient(fakeServer)
			c.Type = "sharev2"
			c.Microversion = "2.37"
			pager := messages.List(c, nil)
			ids := []string{}
			if allPages {
				page, err := pager.AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := messages.ExtractMessages(page)
				th.AssertNoErr(t, err)
				for _, item := range items {
					ids = append(ids, fmt.Sprint(item.ID))
				}
			} else {
				pages := 0
				err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					items, err := messages.ExtractMessages(page)
					th.CheckEquals(t, 1, len(items))
					for _, item := range items {
						ids = append(ids, fmt.Sprint(item.ID))
					}
					return true, err
				})
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 2, pages)
			}
			th.CheckDeepEquals(t, []string{"1", "2"}, ids)
			th.CheckEquals(t, 3, requests)
		})
	}
}

func TestMessagePagePaginationEmpty(t *testing.T) {
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
				fmt.Fprint(w, `{"messages":[]}`)
			})
			c := client.ServiceClient(fakeServer)
			err := messages.List(c, nil).EachPage(context.Background(), func(context.Context, pagination.Page) (bool, error) {
				t.Error("empty page reached callback")
				return false, nil
			})
			th.AssertNoErr(t, err)
			th.CheckEquals(t, 1, requests)
			if status == http.StatusOK {
				requests = 0
				page, err := messages.List(c, nil).AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := messages.ExtractMessages(page)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 0, len(items))
				th.CheckEquals(t, 1, requests)
			}
		})
	}
}

func TestMessagePageNextPageURL(t *testing.T) {
	u, err := url.Parse("https://example.com/resources?offset=7&limit=5&filter=kept")
	th.AssertNoErr(t, err)
	result := pagination.PageResult{URL: *u, Result: gophercloud.Result{StatusCode: 200, Body: map[string]any{"messages": []any{messagePageRecord(1), messagePageRecord(2)}}}}

	page := messages.MessagePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err := page.NextPageURL("https://example.com/")
	th.AssertNoErr(t, err)
	nextURL, err := url.Parse(next)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "5", nextURL.Query().Get("limit"))
	th.CheckEquals(t, "kept", nextURL.Query().Get("filter"))
	th.CheckEquals(t, "9", nextURL.Query().Get("offset"))
	result.RawQuery = "offset=invalid"
	page = messages.MessagePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	_, err = page.NextPageURL("")
	if err == nil {
		t.Fatal("expected invalid offset error")
	}
	result.URL = *u
	result.Body = map[string]any{"messages": "invalid"}

	page = messages.MessagePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	_, err = page.NextPageURL("")
	if err == nil {
		t.Fatal("expected malformed response error")
	}
	result.StatusCode = http.StatusNoContent
	result.Body = nil
	page = messages.MessagePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err = page.NextPageURL("")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "", next)
}
