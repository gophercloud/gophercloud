package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/nodegroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// TestGetNodeGroupSuccess gets a node group successfully.
func TestGetNodeGroupSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleGetNodeGroupSuccess(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	ng, err := nodegroups.Get(context.TODO(), sc, clusterUUID, nodeGroup1UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expectedNodeGroup1, *ng)
}

// TestGetNodeGroupNotFound tries to get a node group which does not exist.
func TestGetNodeGroupNotFound(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleGetNodeGroupNotFound(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	_, err := nodegroups.Get(context.TODO(), sc, clusterUUID, badNodeGroupUUID).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestGetNodeGroupClusterNotFound tries to get a node group in
// a cluster which does not exist.
func TestGetNodeGroupClusterNotFound(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleGetNodeGroupClusterNotFound(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	_, err := nodegroups.Get(context.TODO(), sc, badClusterUUID, badNodeGroupUUID).Extract()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestListNodeGroupsSuccess lists the node groups of a cluster successfully.
func TestListNodeGroupsSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleListNodeGroupsSuccess(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleListNodeGroupsLimitSuccess(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleListNodeGroupsClusterNotFound(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	_, err := nodegroups.List(sc, clusterUUID, nodegroups.ListOpts{}).AllPages(context.TODO())
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestCreateNodeGroupSuccess creates a node group successfully.
func TestCreateNodeGroupSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleCreateNodeGroupSuccess(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleCreateNodeGroupDuplicate(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleCreateNodeGroupMaster(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleCreateNodeGroupBadSizes(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleUpdateNodeGroupSuccess(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleUpdateNodeGroupInternal(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleUpdateNodeGroupBadField(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleUpdateNodeGroupBadMin(t)

	sc := fake.ServiceClient()
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleDeleteNodeGroupSuccess(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, clusterUUID, nodeGroup2UUID).ExtractErr()
	th.AssertNoErr(t, err)
}

// TestDeleteNodeGroupNotFound tries to delete a node group that does not exist.
func TestDeleteNodeGroupNotFound(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleDeleteNodeGroupNotFound(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, clusterUUID, badNodeGroupUUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestDeleteNodeGroupClusterNotFound tries to delete a node group in a cluster that does not exist.
func TestDeleteNodeGroupClusterNotFound(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleDeleteNodeGroupClusterNotFound(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, badClusterUUID, badNodeGroupUUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

// TestDeleteNodeGroupDefault tries to delete a protected default node group.
func TestDeleteNodeGroupDefault(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleDeleteNodeGroupDefault(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	err := nodegroups.Delete(context.TODO(), sc, clusterUUID, nodeGroup2UUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}
