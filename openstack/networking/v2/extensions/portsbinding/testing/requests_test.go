package testing

import (
	"context"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListSuccessfully(t, fakeServer)

	type PortWithExt struct {
		ports.Port
		portsbinding.PortsBindingExt
	}
	var actual []PortWithExt

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
			PortsBindingExt: portsbinding.PortsBindingExt{
				VNICType: "normal",
				HostID:   "devstack",
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

	HandleGet(t, fakeServer)

	var s struct {
		ports.Port
		portsbinding.PortsBindingExt
	}

	err := ports.Get(context.TODO(), fake.ServiceClient(fakeServer), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "ACTIVE", s.Status)
	th.AssertEquals(t, "", s.Name)
	th.AssertTrue(t, s.AdminStateUp)
	th.AssertEquals(t, "a87cc70a-3e15-4acf-8205-9b711a3531b7", s.NetworkID)
	th.AssertEquals(t, "7e02058126cc4950b75f9970368ba177", s.TenantID)
	th.AssertEquals(t, "network:router_interface", s.DeviceOwner)
	th.AssertEquals(t, "fa:16:3e:23:fd:d7", s.MACAddress)
	th.AssertDeepEquals(t, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.1"},
	}, s.FixedIPs)
	th.AssertEquals(t, "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", s.ID)
	th.AssertDeepEquals(t, []string{}, s.SecurityGroups)
	th.AssertEquals(t, "5e3898d7-11be-483e-9732-b2f5eccd2b2e", s.DeviceID)

	th.AssertEquals(t, "devstack", s.HostID)
	th.AssertEquals(t, "normal", s.VNICType)
	th.AssertEquals(t, "ovs", s.VIFType)
	th.AssertDeepEquals(t, map[string]any{"port_filter": true, "ovs_hybrid_plug": true}, s.VIFDetails)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreate(t, fakeServer)

	var s struct {
		ports.Port
		portsbinding.PortsBindingExt
	}

	asu := true
	portCreateOpts := ports.CreateOpts{
		Name:         "private-port",
		AdminStateUp: &asu,
		NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
		},
		SecurityGroups: &[]string{"foo"},
	}

	createOpts := portsbinding.CreateOptsExt{
		CreateOptsBuilder: portCreateOpts,
		HostID:            "HOST1",
		VNICType:          "normal",
	}

	err := ports.Create(context.TODO(), fake.ServiceClient(fakeServer), createOpts).ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "DOWN", s.Status)
	th.AssertEquals(t, "private-port", s.Name)
	th.AssertTrue(t, s.AdminStateUp)
	th.AssertEquals(t, "a87cc70a-3e15-4acf-8205-9b711a3531b7", s.NetworkID)
	th.AssertEquals(t, "d6700c0c9ffa4f1cb322cd4a1f3906fa", s.TenantID)
	th.AssertEquals(t, "", s.DeviceOwner)
	th.AssertEquals(t, "fa:16:3e:c9:cb:f0", s.MACAddress)
	th.AssertDeepEquals(t, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	}, s.FixedIPs)
	th.AssertEquals(t, "65c0ee9f-d634-4522-8954-51021b570b0d", s.ID)
	th.AssertDeepEquals(t, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"}, s.SecurityGroups)
	th.AssertEquals(t, "HOST1", s.HostID)
	th.AssertEquals(t, "normal", s.VNICType)
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := ports.Create(context.TODO(), fake.ServiceClient(fakeServer), portsbinding.CreateOptsExt{CreateOptsBuilder: ports.CreateOpts{}})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdate(t, fakeServer)

	var s struct {
		ports.Port
		portsbinding.PortsBindingExt
	}

	name := "new_port_name"
	portUpdateOpts := ports.UpdateOpts{
		Name: &name,
		FixedIPs: []ports.IP{
			{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
		},
		SecurityGroups: &[]string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"},
	}

	hostID := "HOST1"
	updateOpts := portsbinding.UpdateOptsExt{
		UpdateOptsBuilder: portUpdateOpts,
		HostID:            &hostID,
		VNICType:          "normal",
	}

	err := ports.Update(context.TODO(), fake.ServiceClient(fakeServer), "65c0ee9f-d634-4522-8954-51021b570b0d", updateOpts).ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "new_port_name", s.Name)
	th.AssertDeepEquals(t, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	}, s.FixedIPs)
	th.AssertDeepEquals(t, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"}, s.SecurityGroups)
	th.AssertEquals(t, "HOST1", s.HostID)
	th.AssertEquals(t, "normal", s.VNICType)
}
