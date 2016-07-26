// +build acceptance compute servers

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestServersList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := servers.List(client, servers.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve servers: %v", err)
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		t.Fatalf("Unable to extract servers: %v", err)
	}

	for _, server := range allServers {
		printServer(t, &server)
	}
}

func TestServersCreateDestroy(t *testing.T) {
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
		t.Errorf("Unable to retrieve server: %v", err)
	}
	printServer(t, newServer)

	allAddressPages, err := servers.ListAddresses(client, server.ID).AllPages()
	if err != nil {
		t.Errorf("Unable to list server addresses: %v", err)
	}

	allAddresses, err := servers.ExtractAddresses(allAddressPages)
	if err != nil {
		t.Errorf("Unable to extract server addresses: %v", err)
	}

	for network, address := range allAddresses {
		t.Logf("Addresses on %s: %+v", network, address)
	}

	allNetworkAddressPages, err := servers.ListAddressesByNetwork(client, server.ID, choices.NetworkName).AllPages()
	if err != nil {
		t.Errorf("Unable to list server addresses: %v", err)
	}

	allNetworkAddresses, err := servers.ExtractNetworkAddresses(allNetworkAddressPages)
	if err != nil {
		t.Errorf("Unable to extract server addresses: %v", err)
	}

	t.Logf("Addresses on %s:", choices.NetworkName)
	for _, address := range allNetworkAddresses {
		t.Logf("%+v", address)
	}
}

func TestServersUpdate(t *testing.T) {
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
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, client, server)

	alternateName := tools.RandomString("ACPTTEST", 16)
	for alternateName == server.Name {
		alternateName = tools.RandomString("ACPTTEST", 16)
	}

	t.Logf("Attempting to rename the server to %s.", alternateName)

	updateOpts := servers.UpdateOpts{
		Name: alternateName,
	}

	updated, err := servers.Update(client, server.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to rename server: %v", err)
	}

	if updated.ID != server.ID {
		t.Errorf("Updated server ID [%s] didn't match original server ID [%s]!", updated.ID, server.ID)
	}

	err = tools.WaitFor(func() (bool, error) {
		latest, err := servers.Get(client, updated.ID).Extract()
		if err != nil {
			return false, err
		}

		return latest.Name == alternateName, nil
	})
}

func TestServersMetadata(t *testing.T) {
	t.Parallel()

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, client, server)

	metadata, err := servers.UpdateMetadata(client, server.ID, servers.MetadataOpts{
		"foo":  "bar",
		"this": "that",
	}).Extract()
	if err != nil {
		t.Fatalf("Unable to update metadata: %v", err)
	}
	t.Logf("UpdateMetadata result: %+v\n", metadata)

	err = servers.DeleteMetadatum(client, server.ID, "foo").ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete metadatum: %v", err)
	}

	metadata, err = servers.CreateMetadatum(client, server.ID, servers.MetadatumOpts{
		"foo": "baz",
	}).Extract()
	if err != nil {
		t.Fatalf("Unable to create metadatum: %v", err)
	}
	t.Logf("CreateMetadatum result: %+v\n", metadata)

	metadata, err = servers.Metadatum(client, server.ID, "foo").Extract()
	if err != nil {
		t.Fatalf("Unable to get metadatum: %v", err)
	}
	t.Logf("Metadatum result: %+v\n", metadata)
	th.AssertEquals(t, "baz", metadata["foo"])

	metadata, err = servers.Metadata(client, server.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get metadata: %v", err)
	}
	t.Logf("Metadata result: %+v\n", metadata)

	metadata, err = servers.ResetMetadata(client, server.ID, servers.MetadataOpts{}).Extract()
	if err != nil {
		t.Fatalf("Unable to reset metadata: %v", err)
	}
	t.Logf("ResetMetadata result: %+v\n", metadata)
	th.AssertDeepEquals(t, map[string]string{}, metadata)
}

func TestServersActionChangeAdminPassword(t *testing.T) {
	t.Parallel()

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
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, client, server)

	randomPassword := tools.MakeNewPassword(server.AdminPass)
	res := servers.ChangeAdminPassword(client, server.ID, randomPassword)
	if res.Err != nil {
		t.Fatal(res.Err)
	}

	if err = waitForStatus(client, server, "PASSWORD"); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestServersActionReboot(t *testing.T) {
	t.Parallel()

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
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, client, server)

	rebootOpts := &servers.RebootOpts{
		Type: servers.SoftReboot,
	}

	t.Logf("Attempting reboot of server %s", server.ID)
	res := servers.Reboot(client, server.ID, rebootOpts)
	if res.Err != nil {
		t.Fatalf("Unable to reboot server: %v", res.Err)
	}

	if err = waitForStatus(client, server, "REBOOT"); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestServersActionRebuild(t *testing.T) {
	t.Parallel()

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
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, client, server)

	t.Logf("Attempting to rebuild server %s", server.ID)

	rebuildOpts := servers.RebuildOpts{
		Name:      tools.RandomString("ACPTTEST", 16),
		AdminPass: tools.MakeNewPassword(server.AdminPass),
		ImageID:   choices.ImageID,
	}

	rebuilt, err := servers.Rebuild(client, server.ID, rebuildOpts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	if rebuilt.ID != server.ID {
		t.Errorf("Expected rebuilt server ID of [%s]; got [%s]", server.ID, rebuilt.ID)
	}

	if err = waitForStatus(client, rebuilt, "REBUILD"); err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, rebuilt, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestServersActionResizeConfirm(t *testing.T) {
	t.Parallel()

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, client, server)

	t.Logf("Attempting to resize server %s", server.ID)
	resizeServer(t, client, server, choices)

	t.Logf("Attempting to confirm resize for server %s", server.ID)
	if res := servers.ConfirmResize(client, server.ID); res.Err != nil {
		t.Fatal(res.Err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func TestServersActionResizeRevert(t *testing.T) {
	t.Parallel()

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	server, err := createServer(t, client, choices)
	if err != nil {
		t.Fatal(err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
	defer deleteServer(t, client, server)

	t.Logf("Attempting to resize server %s", server.ID)
	resizeServer(t, client, server, choices)

	t.Logf("Attempting to revert resize for server %s", server.ID)
	if res := servers.RevertResize(client, server.ID); res.Err != nil {
		t.Fatal(res.Err)
	}

	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal(err)
	}
}

func createServer(t *testing.T, client *gophercloud.ServiceClient, choices *ComputeChoices) (*servers.Server, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	networkID, err := getNetworkIDFromTenantNetworks(t, client, choices.NetworkName)
	if err != nil {
		t.Fatalf("Failed to obtain network ID: %v", err)
	}

	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create server: %s", name)

	pwd := tools.MakeNewPassword("")

	server, err := servers.Create(client, servers.CreateOpts{
		Name:      name,
		FlavorRef: choices.FlavorID,
		ImageRef:  choices.ImageID,
		AdminPass: pwd,
		Networks: []servers.Network{
			servers.Network{UUID: networkID},
		},
		Personality: servers.Personality{
			&servers.File{
				Path:     "/etc/test",
				Contents: []byte("hello world"),
			},
		},
	}).Extract()
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	th.AssertEquals(t, pwd, server.AdminPass)

	return server, err
}

func resizeServer(t *testing.T, client *gophercloud.ServiceClient, server *servers.Server, choices *ComputeChoices) {
	opts := &servers.ResizeOpts{
		FlavorRef: choices.FlavorIDResize,
	}
	if res := servers.Resize(client, server.ID, opts); res.Err != nil {
		t.Fatal(res.Err)
	}

	if err := waitForStatus(client, server, "VERIFY_RESIZE"); err != nil {
		t.Fatal(err)
	}
}

func deleteServer(t *testing.T, client *gophercloud.ServiceClient, server *servers.Server) {
	err := servers.Delete(client, server.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete server %s: %s", server.ID, err)
	}

	t.Logf("Deleted server: %s", server.ID)
}

func printServer(t *testing.T, server *servers.Server) {
	t.Logf("ID: %s", server.ID)
	t.Logf("TenantID: %s", server.TenantID)
	t.Logf("UserID: %s", server.UserID)
	t.Logf("Name: %s", server.Name)
	t.Logf("Updated: %s", server.Updated)
	t.Logf("Created: %s", server.Created)
	t.Logf("HostID: %s", server.HostID)
	t.Logf("Status: %s", server.Status)
	t.Logf("Progress: %d", server.Progress)
	t.Logf("AccessIPv4: %s", server.AccessIPv4)
	t.Logf("AccessIPv6: %s", server.AccessIPv6)
	t.Logf("Image: %s", server.Image)
	t.Logf("Flavor: %s", server.Flavor)
	t.Logf("Addresses: %#v", server.Addresses)
	t.Logf("Metadata: %#v", server.Metadata)
	t.Logf("Links: %#v", server.Links)
	t.Logf("KeyName: %s", server.KeyName)
	t.Logf("AdminPass: %s", server.AdminPass)
	t.Logf("SecurityGroups: %#v", server.SecurityGroups)
}
