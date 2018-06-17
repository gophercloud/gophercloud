// +build acceptance clustering policies

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
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

	_, err = CreateNode(t, client, cluster.ID, profile.ID)
	th.AssertNoErr(t, err)
}
