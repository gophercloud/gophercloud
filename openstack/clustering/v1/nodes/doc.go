/*
Package nodes provides information and interaction with the nodes through
the OpenStack Clustering service.

Create, Lists, Update, Delete a node.

Example to List Nodes

	listOpts := nodes.ListOpts{
		Name: "testnode",
	}

	allNodes, err := images.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allNodes, err := nodes.ExtractNodes(allPages)
	if err != nil {
		panic(err)
	}

	for _, node := range allNodes {
		fmt.Printf("%+v\n", node)
	}

*/
package nodes
