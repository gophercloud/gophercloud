package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/nodegroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// TestGetNodeGroupSuccess gets a node group successfully.
func TestGetNodeGroupSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleGetNodeGroupSuccess(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	ng, err := nodegroups.Get(context.TODO(), sc, clusterUUID, nodeGroup1UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expectedNodeGroup1, *ng)
}

// TestGetNodeGroupNotFound tries to get a node group which does not exist.
func TestGetNodeGroupNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleGetNodeGroupNotFound(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	_, err := nodegroups.Get(context.TODO(), sc, clusterUUID, badNodeGroupUUID).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestGetNodeGroupClusterNotFound tries to get a node group in
// a cluster which does not exist.
func TestGetNodeGroupClusterNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleGetNodeGroupClusterNotFound(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	_, err := nodegroups.Get(context.TODO(), sc, badClusterUUID, badNodeGroupUUID).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestListNodeGroupsSuccess lists the node groups of a cluster successfully.
func TestListNodeGroupsSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleListNodeGroupsSuccess(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	ngPages, err := nodegroups.List(sc, clusterUUID, nodegroups.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	ngs, err := nodegroups.ExtractNodeGroups(ngPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(ngs))
	th.AssertEquals(t, nodeGroup1UUID, ngs[0].UUID)
	th.AssertEquals(t, nodeGroup2UUID, ngs[1].UUID)
}

// TestListNodeGroupsLimitSuccess tests listing node groups
// with each returned page limited to one node group and
// also giving a URL to get the next page.
func TestListNodeGroupsLimitSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleListNodeGroupsLimitSuccess(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	listOpts := nodegroups.ListOpts{Limit: 1}
	ngPages, err := nodegroups.List(sc, clusterUUID, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	ngs, err := nodegroups.ExtractNodeGroups(ngPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(ngs))
	th.AssertEquals(t, nodeGroup1UUID, ngs[0].UUID)
	th.AssertEquals(t, nodeGroup2UUID, ngs[1].UUID)
}

// TestListNodeGroupsClusterNotFound tries to list node groups
// of a cluster which does not exist.
func TestListNodeGroupsClusterNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleListNodeGroupsClusterNotFound(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	_, err := nodegroups.List(sc, clusterUUID, nodegroups.ListOpts{}).AllPages(context.TODO())
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestCreateNodeGroupSuccess creates a node group successfully.
func TestCreateNodeGroupSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleCreateNodeGroupSuccess(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	createOpts := nodegroups.CreateOpts{
		Name:        "test-ng",
		MergeLabels: gophercloud.Enabled,
	}

	ng, err := nodegroups.Create(context.TODO(), sc, clusterUUID, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expectedCreatedNodeGroup, *ng)
}

// TestCreateNodeGroupDuplicate creates a node group with
// the same name as an existing one.
func TestCreateNodeGroupDuplicate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleCreateNodeGroupDuplicate(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	createOpts := nodegroups.CreateOpts{
		Name: "default-worker",
	}

	_, err := nodegroups.Create(context.TODO(), sc, clusterUUID, createOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

// TestCreateNodeGroupMaster creates a node group with
// role=master which is not allowed.
func TestCreateNodeGroupMaster(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleCreateNodeGroupMaster(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	createOpts := nodegroups.CreateOpts{
		Name: "new-ng",
		Role: "master",
	}

	_, err := nodegroups.Create(context.TODO(), sc, clusterUUID, createOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}

// TestCreateNodeGroupBadSizes creates a node group with
// min_nodes greater than max_nodes.
func TestCreateNodeGroupBadSizes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleCreateNodeGroupBadSizes(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	maxNodes := 3
	createOpts := nodegroups.CreateOpts{
		Name:         "default-worker",
		MinNodeCount: 5,
		MaxNodeCount: &maxNodes,
	}

	_, err := nodegroups.Create(context.TODO(), sc, clusterUUID, createOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

// TestUpdateNodeGroupSuccess updates a node group successfully.
func TestUpdateNodeGroupSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleUpdateNodeGroupSuccess(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	updateOpts := []nodegroups.UpdateOptsBuilder{
		nodegroups.UpdateOpts{
			Op:    nodegroups.ReplaceOp,
			Path:  "/max_node_count",
			Value: 3,
		},
	}

	ng, err := nodegroups.Update(context.TODO(), sc, clusterUUID, nodeGroup2UUID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expectedUpdatedNodeGroup, *ng)
}

// TestUpdateNodeGroupInternal tries to update an internal
// property of the node group.
func TestUpdateNodeGroupInternal(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleUpdateNodeGroupInternal(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	updateOpts := []nodegroups.UpdateOptsBuilder{
		nodegroups.UpdateOpts{
			Op:    nodegroups.ReplaceOp,
			Path:  "/name",
			Value: "newname",
		},
	}

	_, err := nodegroups.Update(context.TODO(), sc, clusterUUID, nodeGroup2UUID, updateOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}

// TestUpdateNodeGroupBadField tries to update a
// field of the node group that does not exist.
func TestUpdateNodeGroupBadField(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleUpdateNodeGroupBadField(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	updateOpts := []nodegroups.UpdateOptsBuilder{
		nodegroups.UpdateOpts{
			Op:    nodegroups.ReplaceOp,
			Path:  "/bad_field",
			Value: "abc123",
		},
	}

	_, err := nodegroups.Update(context.TODO(), sc, clusterUUID, nodeGroup2UUID, updateOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}

// TestUpdateNodeGroupBadMin tries to set a minimum node count
// greater than the current node count
func TestUpdateNodeGroupBadMin(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleUpdateNodeGroupBadMin(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	updateOpts := []nodegroups.UpdateOptsBuilder{
		nodegroups.UpdateOpts{
			Op:    nodegroups.ReplaceOp,
			Path:  "/min_node_count",
			Value: 5,
		},
	}

	_, err := nodegroups.Update(context.TODO(), sc, clusterUUID, nodeGroup2UUID, updateOpts).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

// TestDeleteNodeGroupSuccess deletes a node group successfully.
func TestDeleteNodeGroupSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleDeleteNodeGroupSuccess(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, clusterUUID, nodeGroup2UUID).ExtractErr()
	th.AssertNoErr(t, err)
}

// TestDeleteNodeGroupNotFound tries to delete a node group that does not exist.
func TestDeleteNodeGroupNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleDeleteNodeGroupNotFound(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, clusterUUID, badNodeGroupUUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestDeleteNodeGroupClusterNotFound tries to delete a node group in a cluster that does not exist.
func TestDeleteNodeGroupClusterNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleDeleteNodeGroupClusterNotFound(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, badClusterUUID, badNodeGroupUUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestDeleteNodeGroupDefault tries to delete a protected default node group.
func TestDeleteNodeGroupDefault(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	handleDeleteNodeGroupDefault(t, fakeServer)

	sc := client.ServiceClient(fakeServer)
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, clusterUUID, nodeGroup2UUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}
