/*
Package nodes provides information and interaction with the nodes through
the OpenStack Clustering service.

Lists all nodes, and creates, shows information for, updates, deletes a node.

Example to Update Nodes

	nodeID := "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1"
	node, err := nodes.Update(serviceClient, nodeID, nodes.UpdateOpts{Name: "new-node-name"}).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", node)

*/
package nodes
