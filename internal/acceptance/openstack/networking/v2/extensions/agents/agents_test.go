//go:build acceptance || networking || agents

package agents

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	spk "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/agents"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAgentsCRUD(t *testing.T) {
	t.Skip("TestAgentsCRUD needs to be re-worked to work with both ML2/OVS and OVN")
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := agents.List(client, agents.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAgents, err := agents.ExtractAgents(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved Networking V2 agents")
	tools.PrintResource(t, allAgents)

	// List DHCP agents
	listOpts := &agents.ListOpts{
		AgentType: "DHCP agent",
	}
	allPages, err = agents.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAgents, err = agents.ExtractAgents(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved Networking V2 DHCP agents")
	tools.PrintResource(t, allAgents)

	// List DHCP agent networks
	for _, agent := range allAgents {
		t.Logf("Retrieving DHCP networks from the agent: %s", agent.ID)
		networks, err := agents.ListDHCPNetworks(context.TODO(), client, agent.ID).Extract()
		th.AssertNoErr(t, err)
		for _, network := range networks {
			t.Logf("Retrieved %q network, assigned to a %q DHCP agent", network.ID, agent.ID)
		}
	}

	// Get a single agent
	agent, err := agents.Get(context.TODO(), client, allAgents[0].ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, agent)

	// Update an agent
	description := "updated agent"
	updateOpts := &agents.UpdateOpts{
		Description: &description,
	}
	agent, err = agents.Update(context.TODO(), client, allAgents[0].ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, agent.Description, description)

	// Restore original description
	agent, err = agents.Update(context.TODO(), client, allAgents[0].ID, &agents.UpdateOpts{Description: &allAgents[0].Description}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, agent.Description, allAgents[0].Description)

	// Assign a new network to a DHCP agent
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	opts := &agents.ScheduleDHCPNetworkOpts{
		NetworkID: network.ID,
	}
	err = agents.ScheduleDHCPNetwork(context.TODO(), client, allAgents[0].ID, opts).ExtractErr()
	th.AssertNoErr(t, err)

	err = agents.RemoveDHCPNetwork(context.TODO(), client, allAgents[0].ID, network.ID).ExtractErr()
	th.AssertNoErr(t, err)

	// skip this part
	t.Skip("Skip DHCP agent deletion")

	// Delete a DHCP agent
	err = agents.Delete(context.TODO(), client, allAgents[0].ID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestBGPAgentCRUD(t *testing.T) {
	timeout := 15 * time.Minute
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "bgp")

	// List BGP Agents
	listOpts := &agents.ListOpts{
		AgentType: "BGP Dynamic Routing Agent",
	}
	allPages, err := agents.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allAgents, err := agents.ExtractAgents(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP agents")
	tools.PrintResource(t, allAgents)

	// Create a BGP Speaker
	bgpSpeaker, err := spk.CreateBGPSpeaker(t, client)
	th.AssertNoErr(t, err)

	// List BGP Speaker-Agent associations
	pages, err := agents.ListDRAgentHostingBGPSpeakers(client, bgpSpeaker.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	bgpAgentsForSpeaker, err := agents.ExtractAgents(pages)
	th.AssertNoErr(t, err)

	// If there are no associations, we can assume the static scheduler is in
	// effect and we must manually associate/disassociate the speaker from the
	// agent.
	//
	// https://docs.openstack.org/neutron-dynamic-routing/latest/admin/agent-scheduler.html
	doManualAssignment := len(bgpAgentsForSpeaker) == 0
	var agentID string

	if doManualAssignment {
		// If using manual assignment, schedule a BGP Speaker to an agent
		agentID = allAgents[0].ID
		opts := agents.ScheduleBGPSpeakerOpts{
			SpeakerID: bgpSpeaker.ID,
		}
		err = agents.ScheduleBGPSpeaker(context.TODO(), client, agentID, opts).ExtractErr()
		th.AssertNoErr(t, err)
		t.Logf("Successfully scheduled speaker %s to agent %s", bgpSpeaker.ID, agentID)
	} else {
		// If using automatic assignment, pick the first agent that the speaker
		// was assigned to (it may be assigned to many, depending on how many
		// nodes there are)
		agentID = bgpAgentsForSpeaker[0].ID
	}

	// Wait for the association to complete.
	pages, err = agents.ListDRAgentHostingBGPSpeakers(client, bgpSpeaker.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	bgpAgentsForSpeaker, err = agents.ExtractAgents(pages)
	th.AssertNoErr(t, err)
	err = tools.WaitForTimeout(
		func(ctx context.Context) (bool, error) {
			bgpAgent, err := agents.Get(ctx, client, agentID).Extract()
			th.AssertNoErr(t, err)
			agentConf := bgpAgent.Configurations
			numOfSpeakers := int(agentConf["bgp_speakers"].(float64))
			t.Logf("Agent %s has %d speaker(s)", agentID, numOfSpeakers)
			return numOfSpeakers == 1, nil
		}, timeout)
	th.AssertNoErr(t, err)

	// Disassociate the BGP Speaker from the agent.
	err = agents.RemoveBGPSpeaker(context.TODO(), client, bgpAgentsForSpeaker[0].ID, bgpSpeaker.ID).ExtractErr()
	th.AssertNoErr(t, err)
	t.Logf("BGP Speaker %s has been removed from agent %s", bgpSpeaker.ID, bgpAgentsForSpeaker[0].ID)

	// Only validate the disassociation if we know the static scheduler is in
	// effect as it'll simply be recreated if we're using the chance scheduler
	// and running in a single node deployment.
	if doManualAssignment {
		err = tools.WaitForTimeout(
			func(ctx context.Context) (bool, error) {
				bgpAgent, err := agents.Get(ctx, client, bgpAgentsForSpeaker[0].ID).Extract()
				th.AssertNoErr(t, err)
				agentConf := bgpAgent.Configurations
				numOfSpeakers := int(agentConf["bgp_speakers"].(float64))
				t.Logf("Agent %s has %d speaker(s)", bgpAgent.ID, numOfSpeakers)
				return numOfSpeakers == 0, nil
			}, timeout)
		th.AssertNoErr(t, err)
	}

	// Delete the BGP Speaker
	err = speakers.Delete(context.TODO(), client, bgpSpeaker.ID).ExtractErr()
	th.AssertNoErr(t, err)
	t.Logf("Successfully deleted the BGP Speaker, %s", bgpSpeaker.ID)
	err = tools.WaitForTimeout(
		func(ctx context.Context) (bool, error) {
			bgpAgent, err := agents.Get(ctx, client, agentID).Extract()
			th.AssertNoErr(t, err)
			agentConf := bgpAgent.Configurations
			numOfSpeakers := int(agentConf["bgp_speakers"].(float64))
			t.Logf("Agent %s has %d speaker(s)", bgpAgent.ID, numOfSpeakers)
			return numOfSpeakers == 0, nil
		}, timeout)
	th.AssertNoErr(t, err)
}
