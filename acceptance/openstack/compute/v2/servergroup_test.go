// +build acceptance compute servergroups

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestServergroupsCreateDelete(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	serverGroup, err := CreateServerGroup(t, client, "anti-affinity")
	if err != nil {
		t.Fatalf("Unable to create server group: %v", err)
	}
	defer DeleteServerGroup(t, client, serverGroup)

	serverGroup, err = servergroups.Get(client, serverGroup.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get server group: %v", err)
	}

	tools.PrintResource(t, serverGroup)

	allPages, err := servergroups.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list server groups: %v", err)
	}

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	if err != nil {
		t.Fatalf("Unable to extract server groups: %v", err)
	}

	var found bool
	for _, sg := range allServerGroups {
		tools.PrintResource(t, serverGroup)

		if sg.ID == serverGroup.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestServergroupsAffinityPolicy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	serverGroup, err := CreateServerGroup(t, client, "affinity")
	if err != nil {
		t.Fatalf("Unable to create server group: %v", err)
	}
	defer DeleteServerGroup(t, client, serverGroup)

	firstServer, err := CreateServerInServerGroup(t, client, serverGroup)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	if err = WaitForComputeStatus(client, firstServer, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer DeleteServer(t, client, firstServer)

	firstServer, err = servers.Get(client, firstServer.ID).Extract()

	secondServer, err := CreateServerInServerGroup(t, client, serverGroup)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}

	if err = WaitForComputeStatus(client, secondServer, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer DeleteServer(t, client, secondServer)

	secondServer, err = servers.Get(client, secondServer.ID).Extract()

	th.AssertEquals(t, firstServer.HostID, secondServer.HostID)
}
