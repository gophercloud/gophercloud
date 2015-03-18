package lbpools

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List returns all load balancer pools that are associated with RackConnect.
func List(c *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	createPage := func(r pagination.PageResult) pagination.Page {
		return PoolPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(c, url, createPage)
}

// Get retrieves a specific load balancer pool (that is associated with RackConnect)
// based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = c.Request("GET", getURL(c, id), gophercloud.RequestOpts{
		JSONResponse: &res.Body,
		OkCodes:      []int{200},
	})
	return res
}

// ListNodes returns all load balancer pool nodes that are associated with RackConnect
// for the given LB pool ID.
func ListNodes(c *gophercloud.ServiceClient, id string) pagination.Pager {
	url := listNodesURL(c, id)
	createPage := func(r pagination.PageResult) pagination.Page {
		return NodePage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(c, url, createPage)
}

// CreateNode adds the cloud server with the given serverID to the load balancer
// pool with the given poolID.
func CreateNode(c *gophercloud.ServiceClient, poolID, serverID string) CreateNodeResult {
	var res CreateNodeResult
	reqBody := map[string]interface{}{
		"cloud_server": map[string]string{
			"id": serverID,
		},
	}
	_, res.Err = c.Request("POST", createNodeURL(c, poolID), gophercloud.RequestOpts{
		JSONBody:     &reqBody,
		JSONResponse: &res.Body,
		OkCodes:      []int{201},
	})
	return res
}

// ListNodesDetails returns all load balancer pool nodes that are associated with RackConnect
// for the given LB pool ID with all their details.
func ListNodesDetails(c *gophercloud.ServiceClient, id string) pagination.Pager {
	url := listNodesDetailsURL(c, id)
	createPage := func(r pagination.PageResult) pagination.Page {
		return NodeDetailsPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(c, url, createPage)
}

// GetNode retrieves a specific LB pool node (that is associated with RackConnect)
// based on its unique ID and the LB pool's unique ID.
func GetNode(c *gophercloud.ServiceClient, poolID, nodeID string) GetNodeResult {
	var res GetNodeResult
	_, res.Err = c.Request("GET", nodeURL(c, poolID, nodeID), gophercloud.RequestOpts{
		JSONResponse: &res.Body,
		OkCodes:      []int{200},
	})
	return res
}

// DeleteNode removes the node with the given nodeID from the LB pool with the
// given poolID.
func DeleteNode(c *gophercloud.ServiceClient, poolID, nodeID string) DeleteNodeResult {
	var res DeleteNodeResult
	_, res.Err = c.Request("DELETE", deleteNodeURL(c, poolID, nodeID), gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	return res
}

// GetNodeDetails retrieves a specific LB pool node's details based on its unique
// ID and the LB pool's unique ID.
func GetNodeDetails(c *gophercloud.ServiceClient, poolID, nodeID string) GetNodeDetailsResult {
	var res GetNodeDetailsResult
	_, res.Err = c.Request("GET", nodeDetailsURL(c, poolID, nodeID), gophercloud.RequestOpts{
		JSONResponse: &res.Body,
		OkCodes:      []int{200},
	})
	return res
}

// NodesOpts are options for bulk adding/deleting nodes to LB pools.
type NodesOpts struct {
	ServerID string
	PoolID   string
}

// CreateNodes adds the cloud servers with the given serverIDs to the corresponding
// load balancer pools with the given poolIDs.
func CreateNodes(c *gophercloud.ServiceClient, opts []NodesOpts) CreateNodesResult {
	var res CreateNodesResult
	reqBody := make([]map[string]interface{}, len(opts))
	for i := range opts {
		reqBody[i] = map[string]interface{}{
			"cloud_server": map[string]string{
				"id": opts[i].ServerID,
			},
			"load_balancer_pool": map[string]string{
				"id": opts[i].PoolID,
			},
		}
	}
	_, res.Err = c.Request("POST", createNodesURL(c), gophercloud.RequestOpts{
		JSONBody:     &reqBody,
		JSONResponse: &res.Body,
		OkCodes:      []int{201},
	})
	return res
}

// DeleteNodes removes the cloud servers with the given serverIDs to the corresponding
// load balancer pools with the given poolIDs.
func DeleteNodes(c *gophercloud.ServiceClient, opts []NodesOpts) DeleteNodesResult {
	var res DeleteNodesResult
	reqBody := make([]map[string]interface{}, len(opts))
	for i := range opts {
		reqBody[i] = map[string]interface{}{
			"cloud_server": map[string]string{
				"id": opts[i].ServerID,
			},
			"load_balancer_pool": map[string]string{
				"id": opts[i].PoolID,
			},
		}
	}
	_, res.Err = c.Request("DELETE", createNodesURL(c), gophercloud.RequestOpts{
		JSONBody: &reqBody,
		OkCodes:  []int{204},
	})
	return res
}

// ListNodesDetailsForServer returns all load balancer pool nodes that are associated with RackConnect
// for the given LB pool ID with all their details for the server with the given serverID.
func ListNodesDetailsForServer(c *gophercloud.ServiceClient, serverID string) pagination.Pager {
	url := listNodesForServerURL(c, serverID)
	createPage := func(r pagination.PageResult) pagination.Page {
		return NodeDetailsForServerPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(c, url, createPage)
}
