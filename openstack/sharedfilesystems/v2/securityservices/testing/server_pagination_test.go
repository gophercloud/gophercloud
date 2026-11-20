package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

type securityServiceServerPaginationOpts struct {
	limit string
}

func (opts securityServiceServerPaginationOpts) ToSecurityServiceListQuery() (string, error) {
	if opts.limit == "" {
		return "", nil
	}
	return "?limit=" + opts.limit, nil
}

// Manila retains pagination parameters as resource filters for admins, but
// removes unknown filters for non-admins before applying offset and limit.
// Preserve the single-page contract, including server-side caps and the admin
// limitation when a caller supplies a pagination parameter.
// https://github.com/openstack/manila/blob/stable/2026.1/manila/api/v2/security_service.py
// https://github.com/openstack/manila/blob/stable/2025.1/manila/api/common.py
func TestSecurityServiceListSinglePageServerFiltering(t *testing.T) {
	for _, admin := range []bool{false, true} {
		for _, limit := range []string{"", "1"} {
			for _, allPages := range []bool{false, true} {
				t.Run(fmt.Sprintf("Admin=%t/Limit=%s/AllPages=%t", admin, limit, allPages), func(t *testing.T) {
					fakeServer := th.SetupHTTP()
					defer fakeServer.Teardown()
					requests := 0
					fakeServer.Mux.HandleFunc("/security-services/detail", func(w http.ResponseWriter, r *http.Request) {
						requests++
						th.TestMethod(t, r, "GET")
						if requests > 1 {
							http.Error(w, "unexpected pagination request", http.StatusInternalServerError)
							return
						}
						query := r.URL.Query()
						services := []map[string]string{{"id": "1"}, {"id": "2"}}
						filtered := []map[string]string{}
						for _, service := range services {
							matches := true
							for key, values := range query {
								switch key {
								case "all_tenants", "status":
									continue
								case "name", "id", "type", "user", "server", "dns_ip", "domain":
								default:
									if !admin {
										continue
									}
								}
								value, ok := service[key]
								if !ok || value != values[0] {
									matches = false
								}
							}
							if matches {
								filtered = append(filtered, service)
							}
						}
						offset := 0
						if value := query.Get("offset"); value != "" {
							var err error
							offset, err = strconv.Atoi(value)
							th.AssertNoErr(t, err)
						}
						th.AssertEquals(t, true, offset >= 0)
						pageSize := 1
						if value := query.Get("limit"); value != "" {
							requested, err := strconv.Atoi(value)
							th.AssertNoErr(t, err)
							th.AssertEquals(t, true, requested >= 0)
							if requested > 0 {
								pageSize = min(pageSize, requested)
							}
						}
						start := min(offset, len(filtered))
						end := min(start+pageSize, len(filtered))
						w.Header().Set("Content-Type", "application/json")
						th.AssertNoErr(t, json.NewEncoder(w).Encode(map[string]any{"security_services": filtered[start:end]}))
					})

					pager := securityservices.List(client.ServiceClient(fakeServer), securityServiceServerPaginationOpts{limit: limit})
					actual := []string{}
					collect := func(page pagination.Page) error {
						items, err := securityservices.ExtractSecurityServices(page)
						if err != nil {
							return err
						}
						for _, item := range items {
							actual = append(actual, item.ID)
						}
						return nil
					}
					if allPages {
						page, err := pager.AllPages(context.Background())
						th.AssertNoErr(t, err)
						th.AssertNoErr(t, collect(page))
					} else {
						err := pager.EachPage(context.Background(), func(_ context.Context, page pagination.Page) (bool, error) {
							return true, collect(page)
						})
						th.AssertNoErr(t, err)
					}
					// The mock server caps responses at one item. An explicit limit
					// also filters out every service for admins.
					expected := []string{"1"}
					if admin && limit != "" {
						expected = []string{}
					}
					th.CheckDeepEquals(t, expected, actual)
					th.CheckEquals(t, 1, requests)
				})
			}
		}
	}
}
