// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/nodegroups"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestNodeGroupsCRUD(t *testing.T) {
	// API not available until Magnum train
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.SkipRelease(t, "stable/rocky")
	clients.SkipRelease(t, "stable/stein")

	client, err := clients.NewContainerInfraV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.9"

	clusterTemplate, err := CreateKubernetesClusterTemplate(t, client)
	th.AssertNoErr(t, err)
	defer DeleteClusterTemplate(t, client, clusterTemplate.UUID)

	clusterID, err := CreateKubernetesCluster(t, client, clusterTemplate.UUID)
	th.AssertNoErr(t, err)
	defer DeleteCluster(t, client, clusterID)

	t.Run("list", func(t *testing.T) { testNodeGroupsList(t, client, clusterID) })
	t.Run("listone-get", func(t *testing.T) { testNodeGroupGet(t, client, clusterID) })
}

func testNodeGroupsList(t *testing.T, client *gophercloud.ServiceClient, clusterID string) {
	allPages, err := nodegroups.List(client, clusterID, nil).AllPages()
	th.AssertNoErr(t, err)

	allNodeGroups, err := nodegroups.ExtractNodeGroups(allPages)
	th.AssertNoErr(t, err)

	// By default two node groups should be created
	th.AssertEquals(t, 2, len(allNodeGroups))
}

func testNodeGroupGet(t *testing.T, client *gophercloud.ServiceClient, clusterID string) {
	listOpts := nodegroups.ListOpts{
		Role: "worker",
	}
	allPages, err := nodegroups.List(client, clusterID, listOpts).AllPages()
	th.AssertNoErr(t, err)

	allNodeGroups, err := nodegroups.ExtractNodeGroups(allPages)
	th.AssertNoErr(t, err)

	// Should be one worker node group
	th.AssertEquals(t, 1, len(allNodeGroups))

	ngID := allNodeGroups[0].UUID

	ng, err := nodegroups.Get(client, clusterID, ngID).Extract()
	th.AssertNoErr(t, err)

	// Should have got the same node group as from the list
	th.AssertEquals(t, ngID, ng.UUID)
	th.AssertEquals(t, "worker", ng.Role)
}
