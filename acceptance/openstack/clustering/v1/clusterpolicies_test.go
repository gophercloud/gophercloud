// +build acceptance clustering actions

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	_ "github.com/gophercloud/gophercloud/openstack/clustering/v1/actions"
	"github.com/gophercloud/gophercloud/openstack/clustering/v1/clusterpolicies"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestClusterPoliciesList(t *testing.T) {
	client, err := clients.NewClusteringV1Client()
	th.AssertNoErr(t, err)
	client.Microversion = "1.5"

	profile, err := CreateProfile(t, client)
	th.AssertNoErr(t, err)
	defer DeleteProfile(t, client, profile.ID)

	cluster, err := CreateCluster(t, client, profile.ID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, cluster.ID)

	policy, err := CreatePolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeletePolicy(t, client, policy.ID)

	err = AttachPolicy(t, client, cluster.ID, policy.ID)
	th.AssertNoErr(t, err)
	defer DetachPolicy(t, client, cluster.ID, policy.ID)

	opts := clusterpolicies.ListOpts{
		Sort: "enabled:asc",
	}

	allPages, err := clusterpolicies.List(client, cluster.ID, opts).AllPages()
	th.AssertNoErr(t, err)

	allClusterPolicies, err := clusterpolicies.ExtractClusterPolicies(allPages)
	th.AssertNoErr(t, err)

	found := false
	for _, clusterpolicies := range allClusterPolicies {
		tools.PrintResource(t, clusterpolicies)
		if (clusterpolicies.ClusterID == cluster.ID && clusterpolicies.PolicyID == policy.ID) {
			found = true
		}
	}
	th.AssertEquals(t, found, true)
}
