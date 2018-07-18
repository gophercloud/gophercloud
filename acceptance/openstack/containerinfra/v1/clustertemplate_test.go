// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestClusterTemplateCRUDOperations(t *testing.T) {
	client, err := clients.NewContainerInfraV1Client()

	ct, err = CreateClusterTemplate(t, client)

	//defer clustertemplates.Delete(client, ct.ID)
	th.AssertEquals(t, "flannel", ct.NetworkDriver)
	clusterTemplateID := ct.ID
	t.Logf(clusterTemplateID)

	// Get a cluster
	ct_get, err := clustertemplates.Get(client, ct.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "k8s", ct_get.Name)
	th.AssertEquals(t, "kubernetes", ct_get.COE)
}
