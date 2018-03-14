// +build acceptance compute servers

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestAttachDetachInterface(t *testing.T) {
	clients.RequireLong(t)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	server, err := CreateServer(t, client)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	defer DeleteServer(t, client, server)

	iface, err := AttachInterface(t, client, server.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer DetachInterface(t, client, server.ID, iface.PortID)

	tools.PrintResource(t, iface)

	server, err = servers.Get(client, server.ID).Extract()
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	for _, networkAddresses := range server.Addresses[choices.NetworkName].([]interface{}) {
		address := networkAddresses.(map[string]interface{})
		if address["OS-EXT-IPS:type"] == "fixed" {
			fixedIP := address["addr"].(string)

			for _, v := range iface.FixedIPs {
				if fixedIP == v.IPAddress {
					found = true
				}
			}
		}
	}

	th.AssertEquals(t, found, true)
}
