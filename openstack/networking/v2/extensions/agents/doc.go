/*
Package agents provides the ability to retrieve and manage Agents through the Neutron API.

Example of Listing Agents

	listOpts := agents.ListOpts{
		AgentType: "Open vSwitch agent",
	}

	allPages, err := agents.List(networkClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allAgents, err := agents.ExtractAgents(allPages)
	if err != nil {
		panic(err)
	}

	for _, agent := range allAgents {
		fmt.Printf("%+v\n", agent)
	}

Example to Get an Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	agent, err := agents.Get(networkClient, agentID).Extract()
	if err != nil {
		panic(err)
	}

Example to Update an Agent

	adminStateUp := true
	description := "agent description"
	updateOpts := &agents.UpdateOpts{
		Description:  &description,
		AdminStateUp: &adminStateUp,
	}
	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	agent, err := agents.Update(networkClient, agentID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete an Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	err := agents.Delete(networkClient, agentID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to List Networks hosted by a DHCP Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	networks, err := agents.ListDHCPNetworks(networkClient, agentID).Extract()
	if err != nil {
		panic(err)
	}

	for _, network := range networks {
		fmt.Printf("%+v\n", network)
	}

Example to Schedule a network to a DHCP Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	opts := &agents.ScheduleDHCPNetworkOpts{
		NetworkID: "1ae075ca-708b-4e66-b4a7-b7698632f05f",
	}
	err := agents.ScheduleDHCPNetwork(networkClient, agentID, opts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Remove a network from a DHCP Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	networkID := "1ae075ca-708b-4e66-b4a7-b7698632f05f"
	err := agents.RemoveDHCPNetwork(networkClient, agentID, networkID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to List bgp speakers by dragent

	spks, err := agents.ListBGPSpeakers(c, agentID).Extract()
	if err != nil {
		log.Panic(nil)
	}
	for _, spk := range spks {
		log.Printf("agent: %s, speaker: %s", agentID, spk.ID)
	}

Example to Schedule bgp speaker to dragent

	err := agents.ScheduleBGPSpeaker(c, agentID, speakerID).ExtractErr()
	if err != nil {
		log.Panic(err)
	}

Example to Remove bgp speaker from dragent

        err := agents.RemoveBGPSpeaker(c, agentID, speakerID).ExtractErr()
        if err != nil {
                log.Panic(err)
        }

Example to list dragents hosting specific bgp speaker

	pages, err := agents.ListDRAgentHostingBGPSpeakers(client, os.Args[1]).AllPages()
        if err != nil {
                log.Panic(err)
        }
        allAgents, err := agents.ExtractAgents(pages)
        if err != nil {
                log.Panic(err)
        }
        for _, a := range allAgents {
                log.Printf("%+v", a)
        }
*/

package agents
