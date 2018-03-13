/*
Package nodes provides information and interaction with the nodes through
the OpenStack Compute service.

Lists all nodes, and creates, shows information for, updates, deletes a node.

Example to List Nodes

	listOpts := nodes.ListOpts{
		Name: "testnode",
	}

	allNodes, err := images.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allNodes, err := images.ExtractNodes(allPages)
	if err != nil {
		panic(err)
	}

	for _, node := range allNodes {
		fmt.Printf("%+v\n", node)
	}

Example to Create a Node

	createOpts := nodes.CreateOpts{
    Name:       "testnode"
    ClusterID:  "b7b870ee-d3c5-4a93-b9d7-846c53b2c2dc"
    Role:       "master"
	}

	node, err := nodes.Create(computeClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package nodes
