// +build acceptance clustering policies

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestNodesCRUD(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	profile, err := CreateProfile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteProfile(t, client, profile.ID)

	cluster, err := CreateCluster(t, client, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, cluster.ID)

	node, err := CreateNode(t, client, cluster.ID, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node.ID)

	// Test nodes list
	allPages, err := nodes.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allNodes, err := nodes.ExtractNodes(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allNodes {
		if v.ID == node.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	// Test nodes update
	t.Logf("Attempting to update node %s", node.ID)

	updateOpts := nodes.UpdateOpts{
		Metadata: map[string]interface{}{
			"bar": "baz",
		},
	}

	res := nodes.Update(client, node.ID, updateOpts)
	th.AssertNoErr(t, res.Err)

	actionID, err := GetActionID(res.Header)
	th.AssertNoErr(t, err)

	err = WaitForAction(client, actionID)
	th.AssertNoErr(t, err)

	node, err = nodes.Get(client, node.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, node)
	tools.PrintResource(t, node.Metadata)
}

// Performs an operation on a node
func TestNodesOps(t *testing.T) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.4"

	profile, err := CreateProfile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteProfile(t, client, profile.ID)

	cluster, err := CreateCluster(t, client, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, cluster.ID)

	node, err := CreateNode(t, client, cluster.ID, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteNode(t, client, node.ID)

	ops := []nodes.OperationOpts{
		// TODO: Commented out due to backend returns error, as of 2018-12-14
		//{Operation: nodes.RebuildOperation},
		//{Operation: nodes.EvacuateOperation, Params: nodes.OperationParams{"EvacuateHost": node.ID, "EvacuateForce", "True"}},
		{Operation: nodes.RebootOperation, Params: nodes.OperationParams{"type": "SOFT"}},
		{Operation: nodes.ChangePasswordOperation, Params: nodes.OperationParams{"admin_pass": "test"}},
		{Operation: nodes.LockOperation},
		{Operation: nodes.UnlockOperation},
		{Operation: nodes.SuspendOperation},
		{Operation: nodes.ResumeOperation},
		{Operation: nodes.RescueOperation, Params: nodes.OperationParams{"image_ref": choices.ImageID}},
		{Operation: nodes.PauseOperation},
		{Operation: nodes.UnpauseOperation},
		{Operation: nodes.StopOperation},
		{Operation: nodes.StartOperation},
	}

	for _, op := range ops {
		opName := string(op.Operation)
		t.Logf("Attempting to perform '%s' on node: %s", opName, node.ID)
		actionID, res := nodes.Ops(client, node.ID, op).Extract()
		th.AssertNoErr(t, res)

		err = WaitForAction(client, actionID)
		th.AssertNoErr(t, err)

		node, err = nodes.Get(client, node.ID).Extract()
		th.AssertNoErr(t, err)
		th.AssertEquals(t, "Operation '"+opName+"' succeeded", node.StatusReason)
		t.Logf("Successfully performed '%s' on node: %s", opName, node.ID)
	}
}
