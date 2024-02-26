package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/extraroutes"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAddExtraRoutes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c/add_extraroutes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "router": {
        "routes": [
            { "destination" : "10.0.3.0/24", "nexthop" : "10.0.0.13" },
            { "destination" : "10.0.4.0/24", "nexthop" : "10.0.0.14" }
        ]
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "router": {
        "name": "name",
        "id": "8604a0de-7f6b-409a-a47c-a1cc7bc77b2e",
        "routes": [
            { "destination" : "10.0.1.0/24", "nexthop" : "10.0.0.11" },
            { "destination" : "10.0.2.0/24", "nexthop" : "10.0.0.12" },
            { "destination" : "10.0.3.0/24", "nexthop" : "10.0.0.13" },
            { "destination" : "10.0.4.0/24", "nexthop" : "10.0.0.14" }
        ]
    }
}
		`)
	})

	r := []routers.Route{
		{
			DestinationCIDR: "10.0.3.0/24",
			NextHop:         "10.0.0.13",
		},
		{
			DestinationCIDR: "10.0.4.0/24",
			NextHop:         "10.0.0.14",
		},
	}
	options := extraroutes.Opts{Routes: &r}

	n, err := extraroutes.Add(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, n.Routes, []routers.Route{
		{
			DestinationCIDR: "10.0.1.0/24",
			NextHop:         "10.0.0.11",
		},
		{
			DestinationCIDR: "10.0.2.0/24",
			NextHop:         "10.0.0.12",
		},
		{
			DestinationCIDR: "10.0.3.0/24",
			NextHop:         "10.0.0.13",
		},
		{
			DestinationCIDR: "10.0.4.0/24",
			NextHop:         "10.0.0.14",
		},
	})
}

func TestRemoveExtraRoutes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c/remove_extraroutes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "router": {
        "routes": [
            { "destination" : "10.0.3.0/24", "nexthop" : "10.0.0.13" },
            { "destination" : "10.0.4.0/24", "nexthop" : "10.0.0.14" }
        ]
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "router": {
        "name": "name",
        "id": "8604a0de-7f6b-409a-a47c-a1cc7bc77b2e",
        "routes": [
            { "destination" : "10.0.1.0/24", "nexthop" : "10.0.0.11" },
            { "destination" : "10.0.2.0/24", "nexthop" : "10.0.0.12" }
        ]
    }
}
		`)
	})

	r := []routers.Route{
		{
			DestinationCIDR: "10.0.3.0/24",
			NextHop:         "10.0.0.13",
		},
		{
			DestinationCIDR: "10.0.4.0/24",
			NextHop:         "10.0.0.14",
		},
	}
	options := extraroutes.Opts{Routes: &r}

	n, err := extraroutes.Remove(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, n.Routes, []routers.Route{
		{
			DestinationCIDR: "10.0.1.0/24",
			NextHop:         "10.0.0.11",
		},
		{
			DestinationCIDR: "10.0.2.0/24",
			NextHop:         "10.0.0.12",
		},
	})
}
