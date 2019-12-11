package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/nodegroups"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

// TestGetNodeGroupSuccess gets a node group successfully.
func TestGetNodeGroupSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleGetNodeGroupSuccess(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	ng, err := nodegroups.Get(sc, clusterUUID, nodeGroup1UUID).Extract()
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

	_, err := nodegroups.Get(sc, clusterUUID, badNodeGroupUUID).Extract()
	th.AssertEquals(t, true, err != nil)

	_, isNotFound := err.(gophercloud.ErrDefault404)
	th.AssertEquals(t, true, isNotFound)
}

// TestGetNodeGroupClusterNotFound tries to get a node group in
// a cluster which does not exist.
func TestGetNodeGroupClusterNotFound(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleGetNodeGroupClusterNotFound(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	_, err := nodegroups.Get(sc, badClusterUUID, badNodeGroupUUID).Extract()
	th.AssertEquals(t, true, err != nil)

	_, isNotFound := err.(gophercloud.ErrDefault404)
	th.AssertEquals(t, true, isNotFound)
}

// TestListNodeGroupsSuccess lists the node groups of a cluster successfully.
func TestListNodeGroupsSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	handleListNodeGroupsSuccess(t)

	sc := fake.ServiceClient()
	sc.Endpoint = sc.Endpoint + "v1/"

	ngPages, err := nodegroups.List(sc, clusterUUID, nodegroups.ListOpts{}).AllPages()
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
	ngPages, err := nodegroups.List(sc, clusterUUID, listOpts).AllPages()
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

	_, err := nodegroups.List(sc, clusterUUID, nodegroups.ListOpts{}).AllPages()
	th.AssertEquals(t, true, err != nil)

	_, isNotFound := err.(gophercloud.ErrDefault404)
	th.AssertEquals(t, true, isNotFound)
}
