//go:build acceptance || networking || agents
// +build acceptance networking agents

package agents

import (
	"testing"
	"time"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	networking "github.com/bizflycloud/gophercloud/acceptance/openstack/networking/v2"
	spk "github.com/bizflycloud/gophercloud/acceptance/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/networking/v2/extensions/agents"
	"github.com/bizflycloud/gophercloud/openstack/networking/v2/extensions/bgp/speakers"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestAgentsRUD(t *testing.T) {
	t.Skip("TestAgentsRUD needs to be re-worked to work with both ML2/OVS and OVN")
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

func TestBGPAgentRUD(t *testing.T) {
	timeout := 120 * time.Second
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// List BGP Agents
	listOpts := &agents.ListOpts{
		AgentType: "BGP Dynamic Routing Agent",
	}
	allPages, err := agents.List(client, listOpts).AllPages()
	th.AssertNoErr(t, err)

	allAgents, err := agents.ExtractAgents(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP agents")
	tools.PrintResource(t, allAgents)

	// Create a BGP Speaker
	bgpSpeaker, err := spk.CreateBGPSpeaker(t, client)
	th.AssertNoErr(t, err)
	pages, err := agents.ListDRAgentHostingBGPSpeakers(client, bgpSpeaker.ID).AllPages()
	th.AssertNoErr(t, err)
	bgpAgents, err := agents.ExtractAgents(pages)
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, len(bgpAgents), 1)

	// List the BGP Agents that accommodate the BGP Speaker
	err = tools.WaitForTimeout(
		func() (bool, error) {
			flag := true
			for _, agt := range bgpAgents {
				t.Logf("BGP Speaker %s has been scheduled to agent %s", bgpSpeaker.ID, agt.ID)
				bgpAgent, err := agents.Get(client, agt.ID).Extract()
				th.AssertNoErr(t, err)
				numOfSpeakers := int(bgpAgent.Configurations["bgp_speakers"].(float64))
				flag = flag && (numOfSpeakers == 1)
			}
			return flag, nil
		}, timeout)
	th.AssertNoErr(t, err)

	// Remove the BGP Speaker from the first agent
	err = agents.RemoveBGPSpeaker(client, bgpAgents[0].ID, bgpSpeaker.ID).ExtractErr()
	th.AssertNoErr(t, err)
	t.Logf("BGP Speaker %s has been removed from agent %s", bgpSpeaker.ID, bgpAgents[0].ID)
	err = tools.WaitForTimeout(
		func() (bool, error) {
			bgpAgent, err := agents.Get(client, bgpAgents[0].ID).Extract()
			th.AssertNoErr(t, err)
			agentConf := bgpAgent.Configurations
			numOfSpeakers := int(agentConf["bgp_speakers"].(float64))
			t.Logf("Agent %s has %d speakers", bgpAgent.ID, numOfSpeakers)
			return numOfSpeakers == 0, nil
		}, timeout)
	th.AssertNoErr(t, err)

	// Remove all BGP Speakers from the agent
	pages, err = agents.ListBGPSpeakers(client, bgpAgents[0].ID).AllPages()
	th.AssertNoErr(t, err)
	allSpeakers, err := agents.ExtractBGPSpeakers(pages)
	th.AssertNoErr(t, err)
	for _, speaker := range allSpeakers {
		th.AssertNoErr(t, agents.RemoveBGPSpeaker(client, bgpAgents[0].ID, speaker.ID).ExtractErr())
	}

	// Schedule a BGP Speaker to an agent
	opts := agents.ScheduleBGPSpeakerOpts{
		SpeakerID: bgpSpeaker.ID,
	}
	err = agents.ScheduleBGPSpeaker(client, bgpAgents[0].ID, opts).ExtractErr()
	th.AssertNoErr(t, err)
	t.Logf("Successfully scheduled speaker %s to agent %s", bgpSpeaker.ID, bgpAgents[0].ID)

	err = tools.WaitForTimeout(
		func() (bool, error) {
			bgpAgent, err := agents.Get(client, bgpAgents[0].ID).Extract()
			th.AssertNoErr(t, err)
			agentConf := bgpAgent.Configurations
			numOfSpeakers := int(agentConf["bgp_speakers"].(float64))
			t.Logf("Agent %s has %d speakers", bgpAgent.ID, numOfSpeakers)
			return 1 == numOfSpeakers, nil
		}, timeout)
	th.AssertNoErr(t, err)

	// Delete the BGP Speaker
	speakers.Delete(client, bgpSpeaker.ID).ExtractErr()
	th.AssertNoErr(t, err)
	t.Logf("Successfully deleted the BGP Speaker, %s", bgpSpeaker.ID)
	err = tools.WaitForTimeout(
		func() (bool, error) {
			bgpAgent, err := agents.Get(client, bgpAgents[0].ID).Extract()
			th.AssertNoErr(t, err)
			agentConf := bgpAgent.Configurations
			numOfSpeakers := int(agentConf["bgp_speakers"].(float64))
			t.Logf("Agent %s has %d speakers", bgpAgent.ID, numOfSpeakers)
			return 0 == numOfSpeakers, nil
		}, timeout)
	th.AssertNoErr(t, err)
}
