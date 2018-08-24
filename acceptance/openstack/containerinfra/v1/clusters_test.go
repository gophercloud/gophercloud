// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestClustersCRUD(t *testing.T) {
	client, err := clients.NewContainerInfraV1Client()
	th.AssertNoErr(t, err)

	clusterTemplate, err := CreateClusterTemplate(t, client)
	th.AssertNoErr(t, err)
	defer DeleteClusterTemplate(t, client, clusterTemplate.UUID)

	clusterID, err := CreateCluster(t, client, clusterTemplate.UUID)
	tools.PrintResource(t, clusterID)
}
