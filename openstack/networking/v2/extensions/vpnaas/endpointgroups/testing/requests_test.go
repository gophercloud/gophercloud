package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/vpn/endpoint-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "endpoint_group": {
        "endpoints": [
            "10.2.0.0/24",
            "10.3.0.0/24"
        ],
        "type": "cidr",
        "name": "peers"
    }
}     `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "endpoint_group": {
        "description": "",
        "tenant_id": "4ad57e7ce0b24fca8f12b9834d91079d",
        "project_id": "4ad57e7ce0b24fca8f12b9834d91079d",
        "endpoints": [
            "10.2.0.0/24",
            "10.3.0.0/24"
        ],
        "type": "cidr",
        "id": "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
        "name": "peers"
    }
}
    `)
	})

	options := endpointgroups.CreateOpts{
		Name: "peers",
		Type: endpointgroups.TypeCIDR,
		Endpoints: []string{
			"10.2.0.0/24",
			"10.3.0.0/24",
		},
	}
	actual, err := endpointgroups.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	expected := endpointgroups.EndpointGroup{
		Name:        "peers",
		TenantID:    "4ad57e7ce0b24fca8f12b9834d91079d",
		ProjectID:   "4ad57e7ce0b24fca8f12b9834d91079d",
		ID:          "6ecd9cf3-ca64-46c7-863f-f2eb1b9e838a",
		Description: "",
		Endpoints: []string{
			"10.2.0.0/24",
			"10.3.0.0/24",
		},
		Type: "cidr",
	}
	th.AssertDeepEquals(t, expected, *actual)
}
