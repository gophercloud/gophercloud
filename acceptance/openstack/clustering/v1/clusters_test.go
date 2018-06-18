// +build acceptance clustering policies

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/clusters"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestClustersCRUD(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	profile, err := CreateProfile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteProfile(t, client, profile.ID)

	cluster, err := CreateCluster(t, client, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, cluster.ID)

	// Test clusters list
	allPages, err := clusters.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allClusters, err := clusters.ExtractClusters(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allClusters {
		if v.ID == cluster.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	// Test cluster update
	updateOpts := clusters.UpdateOpts{
		Name: cluster.Name + "-UPDATED",
	}

	res := clusters.Update(client, cluster.ID, updateOpts)
	th.AssertNoErr(t, res.Err)

	actionID, err := GetActionID(res.Header)
	th.AssertNoErr(t, err)

	err = WaitForAction(client, actionID)
	th.AssertNoErr(t, err)

	newCluster, err := clusters.Get(client, cluster.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newCluster.Name, cluster.Name+"-UPDATED")

	tools.PrintResource(t, newCluster)
}

func TestClustersResize(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	profile, err := CreateProfile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteProfile(t, client, profile.ID)

	cluster, err := CreateCluster(t, client, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, cluster.ID)

	iTrue := true
	resizeOpts := clusters.ResizeOpts{
		AdjustmentType: clusters.ChangeInCapacityAdjustment,
		Number:         1,
		Strict:         &iTrue,
	}

	actionID, err := clusters.Resize(client, cluster.ID, resizeOpts).Extract()
	th.AssertNoErr(t, err)

	err = WaitForAction(client, actionID)
	th.AssertNoErr(t, err)

	newCluster, err := clusters.Get(client, cluster.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newCluster.DesiredCapacity, 2)

	tools.PrintResource(t, newCluster)
}

func TestClustersScale(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	profile, err := CreateProfile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteProfile(t, client, profile.ID)

	cluster, err := CreateCluster(t, client, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, cluster.ID)

	// reduce cluster size to 0
	count := 1
	scaleInOpts := clusters.ScaleInOpts{
		Count: &count,
	}

	actionID, err := clusters.ScaleIn(client, cluster.ID, scaleInOpts).Extract()
	th.AssertNoErr(t, err)

	err = WaitForAction(client, actionID)
	th.AssertNoErr(t, err)

	newCluster, err := clusters.Get(client, cluster.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newCluster.DesiredCapacity, 0)

	tools.PrintResource(t, newCluster)
}

func TestClustersPolicies(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)

	profile, err := CreateProfile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteProfile(t, client, profile.ID)

	cluster, err := CreateCluster(t, client, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, cluster.ID)

	allPages, err := clusters.ListPolicies(client, cluster.ID, nil).AllPages()
	th.AssertNoErr(t, err)

	allPolicies, err := clusters.ExtractClusterPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allPolicies {
		tools.PrintResource(t, v)
	}
}
