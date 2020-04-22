// +build acceptance networking agents

package agents

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/agents"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestAgentsRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := agents.List(client, agents.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	allAgents, err := agents.ExtractAgents(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved Networking V2 agents")
	tools.PrintResource(t, allAgents)

	// List DHCP agents
	listOpts := &agents.ListOpts{
		AgentType: "DHCP agent",
	}
	allPages, err = agents.List(client, listOpts).AllPages()
	th.AssertNoErr(t, err)

	allAgents, err = agents.ExtractAgents(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved Networking V2 DHCP agents")
	tools.PrintResource(t, allAgents)

	// List DHCP agent networks
	for _, agent := range allAgents {
		t.Logf("Retrieving DHCP networks from the agent: %s", agent.ID)
		networks, err := agents.ListDHCPNetworks(client, agent.ID).Extract()
		th.AssertNoErr(t, err)
		for _, network := range networks {
			t.Logf("Retrieved %q network, assigned to a %q DHCP agent", network.ID, agent.ID)
		}
	}

	// Get a single agent
	agent, err := agents.Get(client, allAgents[0].ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, agent)

	// Update an agent
	description := "updated agent"
	updateOpts := &agents.UpdateOpts{
		Description: &description,
	}
	agent, err = agents.Update(client, allAgents[0].ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, agent.Description, description)

	// Restore original description
	agent, err = agents.Update(client, allAgents[0].ID, &agents.UpdateOpts{Description: &allAgents[0].Description}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, agent.Description, allAgents[0].Description)

	// skip this part
	// t.Skip("Skip DHCP agent network scheduling")

	// Assign a new network to a DHCP agent
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	opts := &agents.ScheduleDHCPNetworkOpts{
		NetworkID: network.ID,
	}
	err = agents.ScheduleDHCPNetwork(client, allAgents[0].ID, opts).ExtractErr()
	th.AssertNoErr(t, err)

	err = agents.RemoveDHCPNetwork(client, allAgents[0].ID, network.ID).ExtractErr()
	th.AssertNoErr(t, err)

	// skip this part
	t.Skip("Skip DHCP agent deletion")

	// Delete a DHCP agent
	err = agents.Delete(client, allAgents[0].ID).ExtractErr()
	th.AssertNoErr(t, err)
}
