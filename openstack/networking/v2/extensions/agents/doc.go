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
*/
package agents
