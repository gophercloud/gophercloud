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
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/attachments"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func attachmentPageRecord(index int) map[string]any { return map[string]any{"id": strconv.Itoa(index)} }

func TestAttachmentPagePagination(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				th.TestMethod(t, r, "GET")
				th.TestHeader(t, r, "OpenStack-API-Version", "volume 3.27")
				if requests > 3 {
					http.Error(w, "pagination did not terminate", http.StatusInternalServerError)
					return
				}
				index := 1
				if value := r.URL.Query().Get("marker"); value != "" {
					position, err := strconv.Atoi(value)
					th.AssertNoErr(t, err)
					index = position + 1
				}
				items := []any{}
				if index <= 2 {
					items = append(items, attachmentPageRecord(index))
				}
				body := map[string]any{"attachments": items, "unrelated": []any{"not a resource"}}

				w.Header().Set("Content-Type", "application/json")
				th.AssertNoErr(t, json.NewEncoder(w).Encode(body))
			})
			c := client.ServiceClient(fakeServer)
			c.Type = "volumev3"
			c.Microversion = "3.27"
			pager := attachments.List(c, nil)
			ids := []string{}
			if allPages {
				page, err := pager.AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := attachments.ExtractAttachments(page)
				th.AssertNoErr(t, err)
				for _, item := range items {
					ids = append(ids, fmt.Sprint(item.ID))
				}
			} else {
				pages := 0
				err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					items, err := attachments.ExtractAttachments(page)
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

func TestAttachmentPagePaginationEmpty(t *testing.T) {
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
				fmt.Fprint(w, `{"attachments":[]}`)
			})
			c := client.ServiceClient(fakeServer)
			err := attachments.List(c, nil).EachPage(context.Background(), func(context.Context, pagination.Page) (bool, error) {
				t.Error("empty page reached callback")
				return false, nil
			})
			th.AssertNoErr(t, err)
			th.CheckEquals(t, 1, requests)
			if status == http.StatusOK {
				requests = 0
				page, err := attachments.List(c, nil).AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := attachments.ExtractAttachments(page)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 0, len(items))
				th.CheckEquals(t, 1, requests)
			}
		})
	}
}

func TestAttachmentPageNextPageURL(t *testing.T) {
	u, err := url.Parse("https://example.com/resources?offset=7&limit=5&filter=kept")
	th.AssertNoErr(t, err)
	result := pagination.PageResult{URL: *u, Result: gophercloud.Result{StatusCode: 200, Body: map[string]any{"attachments": []any{attachmentPageRecord(1), attachmentPageRecord(2)}}}}

	page := attachments.AttachmentPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err := page.NextPageURL("https://example.com/")
	th.AssertNoErr(t, err)
	nextURL, err := url.Parse(next)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "5", nextURL.Query().Get("limit"))
	th.CheckEquals(t, "kept", nextURL.Query().Get("filter"))
	th.CheckEquals(t, "2", nextURL.Query().Get("marker"))
	th.CheckEquals(t, "", nextURL.Query().Get("offset"))

	result.URL = *u
	result.Body = map[string]any{"attachments": "invalid"}

	page = attachments.AttachmentPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	_, err = page.NextPageURL("")
	if err == nil {
		t.Fatal("expected malformed response error")
	}
	result.StatusCode = http.StatusNoContent
	result.Body = nil
	page = attachments.AttachmentPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err = page.NextPageURL("")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "", next)
}
