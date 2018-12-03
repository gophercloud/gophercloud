package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, NetworksVLANTransparentListResult)
	})

	count := 0

	iTrue := true
	networkListOpts := networks.ListOpts{}
	listOpts := vlantransparent.ListOptsExt{
		ListOptsBuilder: networkListOpts,
		VLANTransparent: &iTrue,
	}

	networks.List(fake.ServiceClient(), listOpts).EachPage(func(page pagination.Page) (bool, error) {
		count++

		type networkVLANTransparentExt struct {
			networks.Network
			vlantransparent.TransparentExt
		}
		var networkWithVLANTransparentExt []networkVLANTransparentExt

		err := networks.ExtractNetworksInto(page, &networkWithVLANTransparentExt)
		if err != nil {
			t.Errorf("Failed to extract networks: %v", err)
			return false, nil
		}

		networksCount := len(networkWithVLANTransparentExt)
		if networksCount != 2 {
			t.Fatalf("Expected 2 networks, got %d", networksCount)
		}

		th.AssertEquals(t, "db193ab3-96e3-4cb3-8fc5-05f4296d0324", networkWithVLANTransparentExt[0].ID)
		th.AssertEquals(t, "private", networkWithVLANTransparentExt[0].Name)
		th.AssertEquals(t, true, networkWithVLANTransparentExt[0].AdminStateUp)
		th.AssertEquals(t, "ACTIVE", networkWithVLANTransparentExt[0].Status)
		th.AssertDeepEquals(t, []string{"08eae331-0402-425a-923c-34f7cfe39c1b"}, networkWithVLANTransparentExt[0].Subnets)
		th.AssertEquals(t, "26a7980765d0414dbc1fc1f88cdb7e6e", networkWithVLANTransparentExt[0].TenantID)
		th.AssertEquals(t, false, networkWithVLANTransparentExt[0].Shared)
		th.AssertEquals(t, true, networkWithVLANTransparentExt[0].TransparentExt.VLANTransparent)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
