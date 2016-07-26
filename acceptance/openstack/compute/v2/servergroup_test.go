// +build acceptance compute servergroups

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestServergroupsList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := servergroups.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list server groups: %v", err)
	}

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	if err != nil {
		t.Fatalf("Unable to extract server groups: %v", err)
	}

	for _, serverGroup := range allServerGroups {
		printServerGroup(t, &serverGroup)
	}
}

func TestServergroupsCreate(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	serverGroup, err := createServerGroup(t, client, "anti-affinity")
	if err != nil {
		t.Fatalf("Unable to create server group: %v", err)
	}
	defer deleteServerGroup(t, client, serverGroup)

	serverGroup, err = servergroups.Get(client, serverGroup.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get server group: %v", err)
	}

	printServerGroup(t, serverGroup)
}

func TestServergroupsAffinityPolicy(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	serverGroup, err := createServerGroup(t, client, "affinity")
	if err != nil {
		t.Fatalf("Unable to create server group: %v", err)
	}
	defer deleteServerGroup(t, client, serverGroup)

	firstServer, err := createServerInServerGroup(t, client, choices, serverGroup)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	if err = waitForStatus(client, firstServer, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, firstServer)

	firstServer, err = servers.Get(client, firstServer.ID).Extract()

	secondServer, err := createServerInServerGroup(t, client, choices, serverGroup)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	if err = waitForStatus(client, secondServer, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, secondServer)

	secondServer, err = servers.Get(client, secondServer.ID).Extract()

	if firstServer.HostID != secondServer.HostID {
		t.Fatalf("%s and %s were not scheduled on the same host.", firstServer.ID, secondServer.ID)
	}
}

func createServerGroup(t *testing.T, client *gophercloud.ServiceClient, policy string) (*servergroups.ServerGroup, error) {
	sg, err := servergroups.Create(client, &servergroups.CreateOpts{
		Name:     "test",
		Policies: []string{policy},
	}).Extract()

	if err != nil {
		t.Fatalf("Unable to create server group: %v", err)
	}

	return sg, nil
}

func createServerInServerGroup(t *testing.T, client *gophercloud.ServiceClient, choices *ComputeChoices, serverGroup *servergroups.ServerGroup) (*servers.Server, error) {
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

	serverCreateOpts := servers.CreateOpts{
		Name:      name,
		FlavorRef: choices.FlavorID,
		ImageRef:  choices.ImageID,
		AdminPass: pwd,
		Networks: []servers.Network{
			servers.Network{UUID: networkID},
		},
	}

	schedulerHintsOpts := schedulerhints.CreateOptsExt{
		serverCreateOpts,
		schedulerhints.SchedulerHints{
			Group: serverGroup.ID,
		},
	}
	server, err := servers.Create(client, schedulerHintsOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	return server, err
}

func deleteServerGroup(t *testing.T, client *gophercloud.ServiceClient, serverGroup *servergroups.ServerGroup) {
	err := servergroups.Delete(client, serverGroup.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete server group %s: %v", serverGroup.ID, err)
	}

	t.Logf("Deleted server group %s", serverGroup.ID)
}

func printServerGroup(t *testing.T, serverGroup *servergroups.ServerGroup) {
	t.Logf("ID: %s", serverGroup.ID)
	t.Logf("Name: %s", serverGroup.Name)
	t.Logf("Policies: %#v", serverGroup.Policies)
	t.Logf("Members: %#v", serverGroup.Members)
	t.Logf("Metadata: %#v", serverGroup.Metadata)
}
