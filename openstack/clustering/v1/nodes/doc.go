/*
Package nodes provides information and interaction with the nodes through
the OpenStack Compute service.

Lists all nodes, and creates, shows information for, updates, deletes a node.

Example to Get Node

	nodeName := "node123"
	node, err := nodes.Get(clusteringClient, nodeName).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", node)

*/
package nodes
