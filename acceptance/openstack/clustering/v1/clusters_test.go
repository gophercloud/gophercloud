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

	err = WaitForAction(client, actionID, 600)
	th.AssertNoErr(t, err)

	newCluster, err := clusters.Get(client, cluster.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newCluster.Name, cluster.Name+"-UPDATED")

	tools.PrintResource(t, newCluster)
}
