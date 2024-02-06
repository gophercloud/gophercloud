package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput provides a single page of ServiceCatalog results.
const ListOutput = `
{
    "catalog": [
        {
            "endpoints": [
                {
                    "id": "39dc322ce86c4111b4f06c2eeae0841b",
                    "interface": "public",
                    "region": "RegionOne",
                    "url": "http://localhost:5000"
                },
                {
                    "id": "ec642f27474842e78bf059f6c48f4e99",
                    "interface": "internal",
                    "region": "RegionOne",
                    "url": "http://localhost:5000"
                },
                {
                    "id": "c609fc430175452290b62a4242e8a7e8",
                    "interface": "admin",
                    "region": "RegionOne",
                    "url": "http://localhost:5000"
                }
            ],
            "id": "4363ae44bdf34a3981fde3b823cb9aa2",
            "type": "identity",
            "name": "keystone"
        }
    ],
    "links": {
        "self": "https://example.com/identity/v3/catalog",
        "previous": null,
        "next": null
    }
}
`

// ExpectedCatalogSlice is the slice of domains expected to be returned from ListOutput.
var ExpectedCatalogSlice = []tokens.CatalogEntry{{
	ID:   "4363ae44bdf34a3981fde3b823cb9aa2",
	Name: "keystone",
	Type: "identity",
	Endpoints: []tokens.Endpoint{
		{
			ID:        "39dc322ce86c4111b4f06c2eeae0841b",
			Interface: "public",
			Region:    "RegionOne",
			URL:       "http://localhost:5000",
		},
		{
			ID:        "ec642f27474842e78bf059f6c48f4e99",
			Interface: "internal",
			Region:    "RegionOne",
			URL:       "http://localhost:5000",
		},
		{
			ID:        "c609fc430175452290b62a4242e8a7e8",
			Region:    "RegionOne",
			Interface: "admin",
			URL:       "http://localhost:5000",
		},
	},
}}

// HandleListCatalogSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that responds with a list of two domains.
func HandleListCatalogSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/auth/catalog", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, ListOutput)
	})
}
