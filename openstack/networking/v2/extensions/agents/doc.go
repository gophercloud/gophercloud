/*
Package agents provides the ability to retrieve and manage Agents through the Neutron API.

Example of Listing Agents

	listOpts := agents.ListOpts{
		AgentType: "Open vSwitch agent",
	}

	allPages, err := agents.List(networkClient, listOpts).AllPages(context.TODO())
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
	agent, err := agents.Get(context.TODO(), networkClient, agentID).Extract()
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
	agent, err := agents.Update(context.TODO(), networkClient, agentID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete an Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	err := agents.Delete(context.TODO(), networkClient, agentID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to List Networks hosted by a DHCP Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	networks, err := agents.ListDHCPNetworks(context.TODO(), networkClient, agentID).Extract()
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
	err := agents.ScheduleDHCPNetwork(context.TODO(), networkClient, agentID, opts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Remove a network from a DHCP Agent

	agentID := "76af7b1f-d61b-4526-94f7-d2e14e2698df"
	networkID := "1ae075ca-708b-4e66-b4a7-b7698632f05f"
	err := agents.RemoveDHCPNetwork(context.TODO(), networkClient, agentID, networkID).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to List BGP speakers by dragent

	pages, err := agents.ListBGPSpeakers(c, agentID).AllPages(context.TODO())
	if err != nil {
		log.Panicf("%v", err)
	}
	allSpeakers, err := agents.ExtractBGPSpeakers(pages)
	if err != nil {
		log.Panicf("%v", err)
	}
	for _, s := range allSpeakers {
		log.Printf("%v", s)
	}

Example to Schedule bgp speaker to dragent

	var opts agents.ScheduleBGPSpeakerOpts
	opts.SpeakerID = speakerID
	err := agents.ScheduleBGPSpeaker(c, agentID, opts).ExtractErr()
	if err != nil {
		log.Panic(err)
	}

Example to Remove bgp speaker from dragent

	err := agents.RemoveBGPSpeaker(c, agentID, speakerID).ExtractErr()
	if err != nil {
		log.Panic(err)
	}

Example to list dragents hosting specific bgp speaker

	pages, err := agents.ListDRAgentHostingBGPSpeakers(client, speakerID).AllPages(context.TODO())
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

Example to list routers scheduled to L3 agent

        routers, err := agents.ListL3Routers(neutron, "655967f5-d6f3-4732-88f5-617b0ff5c356").Extract()
        if err != nil {
            log.Panic(err)
        }

        for _, r := range routers {
            log.Printf("%+v", r)
        }

Example to remove router from L3 agent

	agentID := "0e1095ae-6f36-40f3-8322-8e1c9a5e68ca"
	routerID := "e6fa0457-efc2-491d-ac12-17ab60417efd"
        err = agents.RemoveL3Router(neutron, agentID, routerID).ExtractErr()
        if err != nil {
            log.Panic(err)
        }

Example to schedule router to L3 agent

	agentID := "0e1095ae-6f36-40f3-8322-8e1c9a5e68ca"
	routerID := "e6fa0457-efc2-491d-ac12-17ab60417efd"
	err = agents.ScheduleL3Router(neutron, agentID, agents.ScheduleL3RouterOpts{RouterID: routerID}).ExtractErr()
        if err != nil {
            log.Panic(err)
        }


*/

package agents
