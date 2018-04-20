/*
Package nodes provides information and interaction with the nodes through
the OpenStack Clustering service.

Lists all nodes, and creates, shows information for, updates, deletes a node.

Example to Delete Node

	nodeID := "6dc6d336e3fc4c0a951b5698cd1236ee"
	err := nodes.Delete(serviceClient, nodeID).ExtractErr()
	if err != nil {
		panic(err)
	}

*/
package nodes
