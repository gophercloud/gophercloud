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
	qosrules "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/rules"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func bandwidthLimitRulePageRecord(index int) map[string]any {
	return map[string]any{"id": strconv.Itoa(index)}
}

func TestBandwidthLimitRulePagePagination(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				th.TestMethod(t, r, "GET")

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
					items = append(items, bandwidthLimitRulePageRecord(index))
				}
				body := map[string]any{"bandwidth_limit_rules": items, "unrelated": []any{"not a resource"}}
				if index == 1 {
					body["bandwidth_limit_rules_links"] = []map[string]string{{"rel": "self", "href": fakeServer.Endpoint()}, {"rel": "next", "href": fakeServer.Endpoint() + r.URL.Path[1:] + "?page=2&filter=kept"}}
				}

				w.Header().Set("Content-Type", "application/json")
				th.AssertNoErr(t, json.NewEncoder(w).Encode(body))
			})
			c := client.ServiceClient(fakeServer)

			pager := qosrules.ListBandwidthLimitRules(c, "policy", nil)
			ids := []string{}
			if allPages {
				page, err := pager.AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := qosrules.ExtractBandwidthLimitRules(page)
				th.AssertNoErr(t, err)
				for _, item := range items {
					ids = append(ids, fmt.Sprint(item.ID))
				}
			} else {
				pages := 0
				err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					items, err := qosrules.ExtractBandwidthLimitRules(page)
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
			th.CheckEquals(t, 2, requests)
		})
	}
}

func TestBandwidthLimitRulePagePaginationEmpty(t *testing.T) {
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
				fmt.Fprint(w, `{"bandwidth_limit_rules":[]}`)
			})
			c := client.ServiceClient(fakeServer)
			err := qosrules.ListBandwidthLimitRules(c, "policy", nil).EachPage(context.Background(), func(context.Context, pagination.Page) (bool, error) {
				t.Error("empty page reached callback")
				return false, nil
			})
			th.AssertNoErr(t, err)
			th.CheckEquals(t, 1, requests)
			if status == http.StatusOK {
				requests = 0
				page, err := qosrules.ListBandwidthLimitRules(c, "policy", nil).AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := qosrules.ExtractBandwidthLimitRules(page)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 0, len(items))
				th.CheckEquals(t, 1, requests)
			}
		})
	}
}

func TestBandwidthLimitRulePageNextPageURL(t *testing.T) {
	u, err := url.Parse("https://example.com/resources?offset=7&limit=5&filter=kept")
	th.AssertNoErr(t, err)
	result := pagination.PageResult{URL: *u, Result: gophercloud.Result{StatusCode: 200, Body: map[string]any{"bandwidth_limit_rules": []any{bandwidthLimitRulePageRecord(1), bandwidthLimitRulePageRecord(2)}}}}
	result.Body.(map[string]any)["bandwidth_limit_rules_links"] = []map[string]string{{"rel": "next", "href": "https://example.com/next"}}
	page := qosrules.BandwidthLimitRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err := page.NextPageURL("https://example.com/")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "https://example.com/next", next)

	result.URL = *u
	result.Body = map[string]any{"bandwidth_limit_rules": "invalid"}
	result.Body.(map[string]any)["bandwidth_limit_rules_links"] = "invalid"
	page = qosrules.BandwidthLimitRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	_, err = page.NextPageURL("")
	if err == nil {
		t.Fatal("expected malformed response error")
	}
	result.StatusCode = http.StatusNoContent
	result.Body = nil
	page = qosrules.BandwidthLimitRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err = page.NextPageURL("")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "", next)
}

func dSCPMarkingRulePageRecord(index int) map[string]any {
	return map[string]any{"id": strconv.Itoa(index)}
}

func TestDSCPMarkingRulePagePagination(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				th.TestMethod(t, r, "GET")

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
					items = append(items, dSCPMarkingRulePageRecord(index))
				}
				body := map[string]any{"dscp_marking_rules": items, "unrelated": []any{"not a resource"}}
				if index == 1 {
					body["dscp_marking_rules_links"] = []map[string]string{{"rel": "self", "href": fakeServer.Endpoint()}, {"rel": "next", "href": fakeServer.Endpoint() + r.URL.Path[1:] + "?page=2&filter=kept"}}
				}

				w.Header().Set("Content-Type", "application/json")
				th.AssertNoErr(t, json.NewEncoder(w).Encode(body))
			})
			c := client.ServiceClient(fakeServer)

			pager := qosrules.ListDSCPMarkingRules(c, "policy", nil)
			ids := []string{}
			if allPages {
				page, err := pager.AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := qosrules.ExtractDSCPMarkingRules(page)
				th.AssertNoErr(t, err)
				for _, item := range items {
					ids = append(ids, fmt.Sprint(item.ID))
				}
			} else {
				pages := 0
				err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					items, err := qosrules.ExtractDSCPMarkingRules(page)
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
			th.CheckEquals(t, 2, requests)
		})
	}
}

func TestDSCPMarkingRulePagePaginationEmpty(t *testing.T) {
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
				fmt.Fprint(w, `{"dscp_marking_rules":[]}`)
			})
			c := client.ServiceClient(fakeServer)
			err := qosrules.ListDSCPMarkingRules(c, "policy", nil).EachPage(context.Background(), func(context.Context, pagination.Page) (bool, error) {
				t.Error("empty page reached callback")
				return false, nil
			})
			th.AssertNoErr(t, err)
			th.CheckEquals(t, 1, requests)
			if status == http.StatusOK {
				requests = 0
				page, err := qosrules.ListDSCPMarkingRules(c, "policy", nil).AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := qosrules.ExtractDSCPMarkingRules(page)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 0, len(items))
				th.CheckEquals(t, 1, requests)
			}
		})
	}
}

func TestDSCPMarkingRulePageNextPageURL(t *testing.T) {
	u, err := url.Parse("https://example.com/resources?offset=7&limit=5&filter=kept")
	th.AssertNoErr(t, err)
	result := pagination.PageResult{URL: *u, Result: gophercloud.Result{StatusCode: 200, Body: map[string]any{"dscp_marking_rules": []any{dSCPMarkingRulePageRecord(1), dSCPMarkingRulePageRecord(2)}}}}
	result.Body.(map[string]any)["dscp_marking_rules_links"] = []map[string]string{{"rel": "next", "href": "https://example.com/next"}}
	page := qosrules.DSCPMarkingRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err := page.NextPageURL("https://example.com/")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "https://example.com/next", next)

	result.URL = *u
	result.Body = map[string]any{"dscp_marking_rules": "invalid"}
	result.Body.(map[string]any)["dscp_marking_rules_links"] = "invalid"
	page = qosrules.DSCPMarkingRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	_, err = page.NextPageURL("")
	if err == nil {
		t.Fatal("expected malformed response error")
	}
	result.StatusCode = http.StatusNoContent
	result.Body = nil
	page = qosrules.DSCPMarkingRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err = page.NextPageURL("")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "", next)
}

func minimumBandwidthRulePageRecord(index int) map[string]any {
	return map[string]any{"id": strconv.Itoa(index)}
}

func TestMinimumBandwidthRulePagePagination(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			requests := 0
			fakeServer.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				th.TestMethod(t, r, "GET")

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
					items = append(items, minimumBandwidthRulePageRecord(index))
				}
				body := map[string]any{"minimum_bandwidth_rules": items, "unrelated": []any{"not a resource"}}
				if index == 1 {
					body["minimum_bandwidth_rules_links"] = []map[string]string{{"rel": "self", "href": fakeServer.Endpoint()}, {"rel": "next", "href": fakeServer.Endpoint() + r.URL.Path[1:] + "?page=2&filter=kept"}}
				}

				w.Header().Set("Content-Type", "application/json")
				th.AssertNoErr(t, json.NewEncoder(w).Encode(body))
			})
			c := client.ServiceClient(fakeServer)

			pager := qosrules.ListMinimumBandwidthRules(c, "policy", nil)
			ids := []string{}
			if allPages {
				page, err := pager.AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := qosrules.ExtractMinimumBandwidthRules(page)
				th.AssertNoErr(t, err)
				for _, item := range items {
					ids = append(ids, fmt.Sprint(item.ID))
				}
			} else {
				pages := 0
				err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					items, err := qosrules.ExtractMinimumBandwidthRules(page)
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
			th.CheckEquals(t, 2, requests)
		})
	}
}

func TestMinimumBandwidthRulePagePaginationEmpty(t *testing.T) {
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
				fmt.Fprint(w, `{"minimum_bandwidth_rules":[]}`)
			})
			c := client.ServiceClient(fakeServer)
			err := qosrules.ListMinimumBandwidthRules(c, "policy", nil).EachPage(context.Background(), func(context.Context, pagination.Page) (bool, error) {
				t.Error("empty page reached callback")
				return false, nil
			})
			th.AssertNoErr(t, err)
			th.CheckEquals(t, 1, requests)
			if status == http.StatusOK {
				requests = 0
				page, err := qosrules.ListMinimumBandwidthRules(c, "policy", nil).AllPages(context.Background())
				th.AssertNoErr(t, err)
				items, err := qosrules.ExtractMinimumBandwidthRules(page)
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 0, len(items))
				th.CheckEquals(t, 1, requests)
			}
		})
	}
}

func TestMinimumBandwidthRulePageNextPageURL(t *testing.T) {
	u, err := url.Parse("https://example.com/resources?offset=7&limit=5&filter=kept")
	th.AssertNoErr(t, err)
	result := pagination.PageResult{URL: *u, Result: gophercloud.Result{StatusCode: 200, Body: map[string]any{"minimum_bandwidth_rules": []any{minimumBandwidthRulePageRecord(1), minimumBandwidthRulePageRecord(2)}}}}
	result.Body.(map[string]any)["minimum_bandwidth_rules_links"] = []map[string]string{{"rel": "next", "href": "https://example.com/next"}}
	page := qosrules.MinimumBandwidthRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err := page.NextPageURL("https://example.com/")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "https://example.com/next", next)

	result.URL = *u
	result.Body = map[string]any{"minimum_bandwidth_rules": "invalid"}
	result.Body.(map[string]any)["minimum_bandwidth_rules_links"] = "invalid"
	page = qosrules.MinimumBandwidthRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	_, err = page.NextPageURL("")
	if err == nil {
		t.Fatal("expected malformed response error")
	}
	result.StatusCode = http.StatusNoContent
	result.Body = nil
	page = qosrules.MinimumBandwidthRulePage{LinkedPageBase: pagination.LinkedPageBase{PageResult: result}}
	next, err = page.NextPageURL("")
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "", next)
}
