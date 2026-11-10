package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRequestAuthenticationHeaderPrecedence(t *testing.T) {
	for _, header := range []string{"X-Auth-Token", "Authorization"} {
		for _, mode := range []string{"generated", "override", "omit", "override and omit"} {
			for _, retry := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%s/retry=%t", header, mode, retry), func(t *testing.T) {
					p := new(gophercloud.ProviderClient)
					p.UseTokenLock()
					p.SetToken("original")
					value := func(token string) string {
						if header == "Authorization" {
							return "Bearer " + token
						}
						return token
					}
					if header == "Authorization" {
						p.AuthenticatedHeadersFunc = func(token string) map[string]string {
							return map[string]string{header: value(token)}
						}
					}
					reauths := 0
					p.ReauthFunc = func(context.Context) error {
						reauths++
						p.SetToken("refreshed")
						return nil
					}
					opts := &gophercloud.RequestOpts{OkCodes: []int{http.StatusNoContent}}
					if mode == "override" || mode == "override and omit" {
						opts.MoreHeaders = map[string]string{header: value("caller")}
					}
					if mode == "omit" || mode == "override and omit" {
						opts.OmitHeaders = []string{header}
					}
					fakeServer := th.SetupHTTP()
					defer fakeServer.Teardown()
					requests := 0
					fakeServer.Mux.HandleFunc("/resource", func(w http.ResponseWriter, r *http.Request) {
						requests++
						expected := value("original")
						if requests > 1 {
							expected = value("refreshed")
						}
						switch mode {
						case "override":
							expected = value("caller")
						case "omit", "override and omit":
							expected = ""
						}
						th.CheckEquals(t, expected, r.Header.Get(header))
						if retry && requests == 1 {
							w.WriteHeader(http.StatusUnauthorized)
							return
						}
						w.WriteHeader(http.StatusNoContent)
					})
					_, err := p.Request(context.Background(), http.MethodGet, fakeServer.Endpoint()+"resource", opts)
					th.AssertNoErr(t, err)
					expectedRequests, expectedReauths := 1, 0
					if retry {
						expectedRequests, expectedReauths = 2, 1
					}
					th.CheckEquals(t, expectedRequests, requests)
					th.CheckEquals(t, expectedReauths, reauths)
				})
			}
		}
	}
}
