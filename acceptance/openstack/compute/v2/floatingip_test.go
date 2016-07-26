// +build acceptance compute servers

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestFloatingIPsList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := floatingips.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve floating IPs: %v", err)
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		t.Fatalf("Unable to extract floating IPs: %v", err)
	}

	for _, floatingIP := range allFloatingIPs {
		printFloatingIP(t, &floatingIP)
	}
}

func TestFloatingIPsCreate(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	floatingIP, err := createFloatingIP(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create floating IP: %v", err)
	}
	defer deleteFloatingIP(t, client, floatingIP)

	printFloatingIP(t, floatingIP)
}

func TestFloatingIPsAssociate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, server)

	floatingIP, err := createFloatingIP(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create floating IP: %v", err)
	}
	defer deleteFloatingIP(t, client, floatingIP)

	printFloatingIP(t, floatingIP)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: floatingIP.IP,
	}

	t.Logf("Attempting to associate floating IP %s to instance %s", floatingIP.IP, server.ID)
	err = floatingips.AssociateInstance(client, server.ID, associateOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to associate floating IP %s with server %s: %v", floatingIP.IP, server.ID, err)
	}
	defer disassociateFloatingIP(t, client, floatingIP, server)
	t.Logf("Floating IP %s is associated with Fixed IP %s", floatingIP.IP, floatingIP.FixedIP)

	newFloatingIP, err := floatingips.Get(client, floatingIP.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get floating IP %s: %v", floatingIP.ID, err)
	}

	printFloatingIP(t, newFloatingIP)
}

func TestFloatingIPsFixedIPAssociate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, server)

	newServer, err := servers.Get(client, server.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get server %s: %v", server.ID, err)
	}

	floatingIP, err := createFloatingIP(t, client, choices)
	if err != nil {
		t.Fatalf("Unable to create floating IP: %v", err)
	}
	defer deleteFloatingIP(t, client, floatingIP)

	printFloatingIP(t, floatingIP)

	var fixedIP string
	for _, networkAddresses := range newServer.Addresses[choices.NetworkName].([]interface{}) {
		address := networkAddresses.(map[string]interface{})
		if address["OS-EXT-IPS:type"] == "fixed" {
			if address["version"].(float64) == 4 {
				fixedIP = address["addr"].(string)
			}
		}
	}

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: floatingIP.IP,
		FixedIP:    fixedIP,
	}

	t.Logf("Attempting to associate floating IP %s to instance %s", floatingIP.IP, newServer.ID)
	err = floatingips.AssociateInstance(client, newServer.ID, associateOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to associate floating IP %s with server %s: %v", floatingIP.IP, newServer.ID, err)
	}
	defer disassociateFloatingIP(t, client, floatingIP, newServer)
	t.Logf("Floating IP %s is associated with Fixed IP %s", floatingIP.IP, floatingIP.FixedIP)

	newFloatingIP, err := floatingips.Get(client, floatingIP.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get floating IP %s: %v", floatingIP.ID, err)
	}

	printFloatingIP(t, newFloatingIP)
}

func createFloatingIP(t *testing.T, client *gophercloud.ServiceClient, choices *ComputeChoices) (*floatingips.FloatingIP, error) {
	createOpts := floatingips.CreateOpts{
		Pool: choices.FloatingIPPoolName,
	}
	floatingIP, err := floatingips.Create(client, createOpts).Extract()
	if err != nil {
		return floatingIP, err
	}

	t.Logf("Created floating IP: %s", floatingIP.ID)
	return floatingIP, nil
}

func deleteFloatingIP(t *testing.T, client *gophercloud.ServiceClient, floatingIP *floatingips.FloatingIP) {
	err := floatingips.Delete(client, floatingIP.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete floating IP %s: %v", floatingIP.ID, err)
	}

	t.Logf("Deleted floating IP: %s", floatingIP.ID)
}

func disassociateFloatingIP(t *testing.T, client *gophercloud.ServiceClient, floatingIP *floatingips.FloatingIP, server *servers.Server) {
	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: floatingIP.IP,
	}

	err := floatingips.DisassociateInstance(client, server.ID, disassociateOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to disassociate floating IP %s from server %s: %v", floatingIP.IP, server.ID, err)
	}

	t.Logf("Disassociated floating IP %s from server %s", floatingIP.IP, server.ID)
}

func printFloatingIP(t *testing.T, floatingIP *floatingips.FloatingIP) {
	t.Logf("ID: %s", floatingIP.ID)
	t.Logf("Fixed IP: %s", floatingIP.FixedIP)
	t.Logf("Instance ID: %s", floatingIP.InstanceID)
	t.Logf("IP: %s", floatingIP.IP)
	t.Logf("Pool: %s", floatingIP.Pool)
}
