package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

// PortWithDHCPOpts represents a port with extra DHCP options configuration.
type PortWithDHCPOpts struct {
	ports.Port
	extradhcpopts.DHCPOptsExt
}

// CreatePortWithDHCPOpts will create a port with DHCP options on the specified subnet.
// An error will be returned if the port could not be created.
func CreatePortWithDHCPOpts(t *testing.T, client *gophercloud.ServiceClient, networkID, subnetID string) (*PortWithDHCPOpts, error) {
	portName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create port: %s", portName)

	portCreateOpts := ports.CreateOpts{
		NetworkID:    networkID,
		Name:         portName,
		AdminStateUp: gophercloud.Enabled,
		FixedIPs:     []ports.IP{ports.IP{SubnetID: subnetID}},
	}
	createOpts := extradhcpopts.CreateOptsExt{
		CreateOptsBuilder: portCreateOpts,
		DHCPOpts: []extradhcpopts.DHCPOpts{
			{
				DHCPOptName:  "test_option_1",
				DHCPOptValue: "test_value_1",
			},
		},
	}
	port := &PortWithDHCPOpts{}

	err := ports.Create(client, createOpts).ExtractInto(port)
	if err != nil {
		return nil, err
	}

	if err := v2.WaitForPortToCreate(client, port.ID, 60); err != nil {
		return nil, err
	}

	err = ports.Get(client, port.ID).ExtractInto(port)
	if err != nil {
		return port, err
	}

	t.Logf("Successfully created port: %s", portName)

	return port, nil
}
