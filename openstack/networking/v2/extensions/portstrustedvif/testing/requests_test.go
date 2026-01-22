package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portstrustedvif"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	type PortWithExt struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}
	var actual []PortWithExt

	iTrue := true
	expected := []PortWithExt{
		{
			Port: ports.Port{
				Status:       "ACTIVE",
				Name:         "",
				AdminStateUp: true,
				NetworkID:    "70c1db1f-b701-45bd-96e0-a313ee3430b3",
				TenantID:     "",
				DeviceOwner:  "network:router_gateway",
				MACAddress:   "fa:16:3e:58:42:ed",
				FixedIPs: []ports.IP{
					{
						SubnetID:  "008ba151-0b8c-4a67-98b5-0d2b87666062",
						IPAddress: "172.24.4.2",
					},
				},
				ID:             "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8b",
				SecurityGroups: []string{},
				DeviceID:       "9ae135f4-b6e0-4dad-9e91-3c223e385824",
				CreatedAt:      time.Date(2019, time.June, 30, 4, 15, 37, 0, time.UTC),
				UpdatedAt:      time.Date(2019, time.June, 30, 5, 18, 49, 0, time.UTC),
			},
			PortTrustedVIFExt: portstrustedvif.PortTrustedVIFExt{
				PortTrustedVIF: &iTrue,
			},
		},
	}

	allPages, err := ports.List(fake.ServiceClient(fakeServer), ports.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	err = ports.ExtractPortsInto(allPages, &actual)
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, expected, actual)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/ports/46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetPortTrustedVIFStatusResponse)
	})

	var s struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}

	err := ports.Get(context.TODO(), fake.ServiceClient(fakeServer), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, "DOWN")
	th.AssertEquals(t, s.Name, "private-port")
	th.AssertEquals(t, s.AdminStateUp, true)
	th.AssertEquals(t, s.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, s.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, s.DeviceOwner, "")
	th.AssertEquals(t, s.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.1"},
	})
	th.AssertEquals(t, s.ID, "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2")
	th.AssertEquals(t, s.DeviceID, "")

	if s.PortTrustedVIF == nil {
		t.Fatalf("Expected s.PortTrustedVIF to be not nil")
	}
	th.AssertEquals(t, *s.PortTrustedVIF, false)
}

func TestGetUnset(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/ports/46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetPortTrustedVIFUnsetResponse)
	})

	var s struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}

	err := ports.Get(context.TODO(), fake.ServiceClient(fakeServer), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, "DOWN")
	th.AssertEquals(t, s.Name, "private-port")
	th.AssertEquals(t, s.AdminStateUp, true)
	th.AssertEquals(t, s.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, s.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, s.DeviceOwner, "")
	th.AssertEquals(t, s.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.1"},
	})
	th.AssertEquals(t, s.ID, "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2")
	th.AssertEquals(t, s.DeviceID, "")

	if s.PortTrustedVIF != nil {
		t.Fatalf("Expected s.PortTrustedVIF to be nil")
	}
}

func TestCreateWithPortTrustedVIF(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreatePortTrustedVIFRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreatePortTrustedVIFResponse)
	})

	var portWithExt struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}

	asu := true
	portTrustedVIFStatus := true
	options := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
	}
	createOpts := portstrustedvif.PortCreateOptsExt{
		CreateOptsBuilder: options,
		PortTrustedVIF:    &portTrustedVIFStatus,
	}
	err := ports.Create(context.TODO(), fake.ServiceClient(fakeServer), createOpts).ExtractInto(&portWithExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, portWithExt.Status, "DOWN")
	th.AssertEquals(t, portWithExt.Name, "private-port")
	th.AssertEquals(t, portWithExt.AdminStateUp, true)
	th.AssertEquals(t, portWithExt.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, portWithExt.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, portWithExt.DeviceOwner, "")
	th.AssertEquals(t, portWithExt.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, portWithExt.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, portWithExt.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, *portWithExt.PortTrustedVIF, true)
}

func TestUpdatePortTrustedVIF(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdatePortTrustedVIFRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdatePortTrustedVIFResponse)
	})

	var portWithExt struct {
		ports.Port
		portstrustedvif.PortTrustedVIFExt
	}

	portTrustedVIFStatus := false
	portUpdateOpts := ports.UpdateOpts{}
	updateOpts := portstrustedvif.PortUpdateOptsExt{
		UpdateOptsBuilder: portUpdateOpts,
		PortTrustedVIF:    &portTrustedVIFStatus,
	}

	err := ports.Update(context.TODO(), fake.ServiceClient(fakeServer), "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&portWithExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, portWithExt.Name, "private-port")
	th.AssertDeepEquals(t, *portWithExt.PortTrustedVIF, false)
}
