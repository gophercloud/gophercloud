package testing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	count := 0

	err := ports.List(fake.ServiceClient(), ports.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := ports.ExtractPorts(page)
		if err != nil {
			t.Errorf("Failed to extract subnets: %v", err)
			return false, nil
		}

		expected := []ports.Port{
			{
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
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListWithExtensions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	type portWithExt struct {
		ports.Port
		portsecurity.PortSecurityExt
	}

	var allPorts []portWithExt

	allPages, err := ports.List(fake.ServiceClient(), ports.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	err = ports.ExtractPortsInto(allPages, &allPorts)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, allPorts[0].Status, "ACTIVE")
	th.AssertEquals(t, allPorts[0].PortSecurityEnabled, false)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponse)
	})

	n, err := ports.Get(context.TODO(), fake.ServiceClient(), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertEquals(t, n.Name, "")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "7e02058126cc4950b75f9970368ba177")
	th.AssertEquals(t, n.DeviceOwner, "network:router_interface")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:23:fd:d7")
	th.AssertDeepEquals(t, n.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.1"},
	})
	th.AssertEquals(t, n.ID, "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2")
	th.AssertDeepEquals(t, n.SecurityGroups, []string{})
	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertEquals(t, n.DeviceID, "5e3898d7-11be-483e-9732-b2f5eccd2b2e")
	th.AssertEquals(t, n.CreatedAt, time.Date(2019, time.June, 30, 4, 15, 37, 0, time.UTC))
	th.AssertEquals(t, n.UpdatedAt, time.Date(2019, time.June, 30, 5, 18, 49, 0, time.UTC))
}

func TestGetWithExtensions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponse)
	})

	var portWithExtensions struct {
		ports.Port
		portsecurity.PortSecurityExt
	}

	err := ports.Get(context.TODO(), fake.ServiceClient(), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").ExtractInto(&portWithExtensions)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, portWithExtensions.Status, "ACTIVE")
	th.AssertEquals(t, portWithExtensions.PortSecurityEnabled, false)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateResponse)
	})

	asu := true
	options := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		SecurityGroups: &[]string{"foo"},
		AllowedAddressPairs: []ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}
	n, err := ports.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "DOWN")
	th.AssertEquals(t, n.Name, "private-port")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, n.DeviceOwner, "")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, n.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, n.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertDeepEquals(t, n.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
	th.AssertDeepEquals(t, n.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
}

func TestCreateOmitSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateOmitSecurityGroupsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateOmitSecurityGroupsResponse)
	})

	asu := true
	options := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		AllowedAddressPairs: []ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}
	n, err := ports.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "DOWN")
	th.AssertEquals(t, n.Name, "private-port")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, n.DeviceOwner, "")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, n.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, n.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertDeepEquals(t, n.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
	th.AssertDeepEquals(t, n.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
}

func TestCreateWithNoSecurityGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateWithNoSecurityGroupsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateWithNoSecurityGroupsResponse)
	})

	asu := true
	options := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		SecurityGroups: &[]string{},
		AllowedAddressPairs: []ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}
	n, err := ports.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "DOWN")
	th.AssertEquals(t, n.Name, "private-port")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, n.DeviceOwner, "")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, n.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, n.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertDeepEquals(t, n.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
}

func TestCreateWithPropagateUplinkStatus(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreatePropagateUplinkStatusRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreatePropagateUplinkStatusResponse)
	})

	asu := true
	propagateUplinkStatus := true
	options := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		PropagateUplinkStatus: &propagateUplinkStatus,
	}
	n, err := ports.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "DOWN")
	th.AssertEquals(t, n.Name, "private-port")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, n.DeviceOwner, "")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, n.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, n.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, n.PropagateUplinkStatus, propagateUplinkStatus)
}

func TestCreateWithValueSpecs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateValueSpecRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateValueSpecResponse)
	})

	asu := true
	options := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		SecurityGroups: &[]string{"foo"},
		AllowedAddressPairs: []ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
		ValueSpecs: &map[string]string{
			"test": "value",
		},
	}
	n, err := ports.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "DOWN")
	th.AssertEquals(t, n.Name, "private-port")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, n.DeviceOwner, "")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, n.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, n.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertDeepEquals(t, n.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
	th.AssertDeepEquals(t, n.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
}

func TestCreateWithInvalidValueSpecs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateValueSpecRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateValueSpecResponse)
	})

	asu := true
	options := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		SecurityGroups: &[]string{"foo"},
		AllowedAddressPairs: []ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
		ValueSpecs: &map[string]string{
			// This is a forbidden key
			"shared": "value",
		},
	}

	// We expect an error here since we used a fobidden key in the value specs.
	_, err := ports.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertErr(t, err)

	options.ValueSpecs = &map[string]string{
		// Try to overwrite an existing field
		"name": "overwrite",
	}

	// We expect an error here since the value specs would overwrite an existing field.
	_, err = ports.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertErr(t, err)
}

func TestRequiredCreateOpts(t *testing.T) {
	res := ports.Create(context.TODO(), fake.ServiceClient(), ports.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestCreatePortSecurity(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreatePortSecurityRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreatePortSecurityResponse)
	})

	var portWithExt struct {
		ports.Port
		portsecurity.PortSecurityExt
	}

	asu := true
	iFalse := false
	portCreateOpts := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		SecurityGroups: &[]string{"foo"},
		AllowedAddressPairs: []ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}
	createOpts := portsecurity.PortCreateOptsExt{
		CreateOptsBuilder:   portCreateOpts,
		PortSecurityEnabled: &iFalse,
	}

	err := ports.Create(context.TODO(), fake.ServiceClient(), createOpts).ExtractInto(&portWithExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, portWithExt.Status, "DOWN")
	th.AssertEquals(t, portWithExt.PortSecurityEnabled, false)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})

	name := "new_port_name"
	options := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
		SecurityGroups: &[]string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"},
		AllowedAddressPairs: &[]ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}

	s, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "new_port_name")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	})
	th.AssertDeepEquals(t, s.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
	th.AssertDeepEquals(t, s.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
}

func TestUpdateOmitSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateOmitSecurityGroupsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateOmitSecurityGroupsResponse)
	})

	name := "new_port_name"
	options := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
		AllowedAddressPairs: &[]ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}

	s, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "new_port_name")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	})
	th.AssertDeepEquals(t, s.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
	th.AssertDeepEquals(t, s.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
}

func TestUpdatePropagateUplinkStatus(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdatePropagateUplinkStatusRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdatePropagateUplinkStatusResponse)
	})

	propagateUplinkStatus := true
	options := ports.UpdateOpts{
		PropagateUplinkStatus: &propagateUplinkStatus,
	}

	s, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, s.PropagateUplinkStatus, propagateUplinkStatus)
}

func TestUpdateValueSpecs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateValueSpecsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateValueSpecsResponse)
	})

	options := ports.UpdateOpts{
		ValueSpecs: &map[string]string{
			"test": "update",
		},
	}

	_, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)
}

func TestUpdatePortSecurity(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdatePortSecurityRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdatePortSecurityResponse)
	})

	var portWithExt struct {
		ports.Port
		portsecurity.PortSecurityExt
	}

	iFalse := false
	portUpdateOpts := ports.UpdateOpts{}
	updateOpts := portsecurity.PortUpdateOptsExt{
		UpdateOptsBuilder:   portUpdateOpts,
		PortSecurityEnabled: &iFalse,
	}

	err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&portWithExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, portWithExt.Status, "DOWN")
	th.AssertEquals(t, portWithExt.Name, "private-port")
	th.AssertEquals(t, portWithExt.PortSecurityEnabled, false)
}

func TestUpdateRevision(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeaderUnset(t, r, "If-Match")
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})
	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0e", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "If-Match", "revision_number=42")
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})

	name := "new_port_name"
	options := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
		SecurityGroups: &[]string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"},
		AllowedAddressPairs: &[]ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}
	_, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	revisionNumber := 42
	options.RevisionNumber = &revisionNumber
	_, err = ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0e", options).Extract()
	th.AssertNoErr(t, err)
}

func TestRemoveSecurityGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RemoveSecurityGroupRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RemoveSecurityGroupResponse)
	})

	name := "new_port_name"
	options := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
		SecurityGroups: &[]string{},
		AllowedAddressPairs: &[]ports.AddressPair{
			{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
		},
	}

	s, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "new_port_name")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	})
	th.AssertDeepEquals(t, s.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
	th.AssertDeepEquals(t, s.SecurityGroups, []string(nil))
}

func TestRemoveAllowedAddressPairs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, RemoveAllowedAddressPairsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, RemoveAllowedAddressPairsResponse)
	})

	name := "new_port_name"
	options := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
		SecurityGroups:      &[]string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"},
		AllowedAddressPairs: &[]ports.AddressPair{},
	}

	s, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "new_port_name")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	})
	th.AssertDeepEquals(t, s.AllowedAddressPairs, []ports.AddressPair(nil))
	th.AssertDeepEquals(t, s.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
}

func TestDontUpdateAllowedAddressPairs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, DontUpdateAllowedAddressPairsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, DontUpdateAllowedAddressPairsResponse)
	})

	name := "new_port_name"
	options := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
		SecurityGroups: &[]string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"},
	}

	s, err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "new_port_name")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	})
	th.AssertDeepEquals(t, s.AllowedAddressPairs, []ports.AddressPair{
		{IPAddress: "10.0.0.4", MACAddress: "fa:16:3e:c9:cb:f0"},
	})
	th.AssertDeepEquals(t, s.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := ports.Delete(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertNoErr(t, res.Err)
}

func TestGetWithExtraDHCPOpts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetWithExtraDHCPOptsResponse)
	})

	var s struct {
		ports.Port
		extradhcpopts.ExtraDHCPOptsExt
	}

	err := ports.Get(context.TODO(), fake.ServiceClient(), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, "ACTIVE")
	th.AssertEquals(t, s.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, s.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, s.AdminStateUp, true)
	th.AssertEquals(t, s.Name, "port-with-extra-dhcp-opts")
	th.AssertEquals(t, s.DeviceOwner, "")
	th.AssertEquals(t, s.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.4"},
	})
	th.AssertEquals(t, s.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, s.DeviceID, "")

	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].OptName, "option1")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].OptValue, "value1")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].IPVersion, 4)
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[1].OptName, "option2")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[1].OptValue, "value2")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[1].IPVersion, 4)
}

func TestCreateWithExtraDHCPOpts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateWithExtraDHCPOptsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, CreateWithExtraDHCPOptsResponse)
	})

	adminStateUp := true
	portCreateOpts := ports.CreateOpts{
		Name:         "port-with-extra-dhcp-opts",
		AdminStateUp: &adminStateUp,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
	}

	createOpts := extradhcpopts.CreateOptsExt{
		CreateOptsBuilder: portCreateOpts,
		ExtraDHCPOpts: []extradhcpopts.CreateExtraDHCPOpt{
			{
				OptName:  "option1",
				OptValue: "value1",
			},
		},
	}

	var s struct {
		ports.Port
		extradhcpopts.ExtraDHCPOptsExt
	}

	err := ports.Create(context.TODO(), fake.ServiceClient(), createOpts).ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, "DOWN")
	th.AssertEquals(t, s.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, s.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, s.AdminStateUp, true)
	th.AssertEquals(t, s.Name, "port-with-extra-dhcp-opts")
	th.AssertEquals(t, s.DeviceOwner, "")
	th.AssertEquals(t, s.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, s.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, s.DeviceID, "")

	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].OptName, "option1")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].OptValue, "value1")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].IPVersion, 4)
}

func TestUpdateWithExtraDHCPOpts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateWithExtraDHCPOptsRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateWithExtraDHCPOptsResponse)
	})

	name := "updated-port-with-dhcp-opts"
	portUpdateOpts := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
	}

	edoValue2 := "value2"
	updateOpts := extradhcpopts.UpdateOptsExt{
		UpdateOptsBuilder: portUpdateOpts,
		ExtraDHCPOpts: []extradhcpopts.UpdateExtraDHCPOpt{
			{
				OptName: "option1",
			},
			{
				OptName:  "option2",
				OptValue: &edoValue2,
			},
		},
	}

	var s struct {
		ports.Port
		extradhcpopts.ExtraDHCPOptsExt
	}

	err := ports.Update(context.TODO(), fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, "DOWN")
	th.AssertEquals(t, s.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, s.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, s.AdminStateUp, true)
	th.AssertEquals(t, s.Name, "updated-port-with-dhcp-opts")
	th.AssertEquals(t, s.DeviceOwner, "")
	th.AssertEquals(t, s.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	})
	th.AssertEquals(t, s.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, s.DeviceID, "")

	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].OptName, "option2")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].OptValue, "value2")
	th.AssertDeepEquals(t, s.ExtraDHCPOpts[0].IPVersion, 4)
}

func TestPortsListOpts(t *testing.T) {
	for _, tt := range []struct {
		listOpts ports.ListOpts
		params   []struct{ key, value string }
	}{
		{
			listOpts: ports.ListOpts{
				FixedIPs: []ports.FixedIPOpts{{IPAddress: "1.2.3.4", SubnetID: "42"}},
			},
			params: []struct {
				key   string
				value string
			}{
				{"fixed_ips", "ip_address=1.2.3.4"},
				{"fixed_ips", "subnet_id=42"},
			},
		},
	} {
		v := url.Values{}
		for _, param := range tt.params {
			v.Add(param.key, param.value)
		}
		expected := "?" + v.Encode()

		actual, _ := tt.listOpts.ToPortListQuery()
		th.AssertEquals(t, expected, actual)
	}
}
