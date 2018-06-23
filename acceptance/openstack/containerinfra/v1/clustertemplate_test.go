// +build acceptance containerinfra

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestClusterTemplateCRUDOperations(t *testing.T) {
	client, err := clients.NewContainerInfraV1Client()

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new cluster template
	opts := clustertemplates.CreateOpts{Name: "k8s",
		COE:               "kubernetes",
		ImageID:           choices.COEImageID,
		KeypairID:         "default",
		DockerVolumeSize:  5,
		FlavorID:          choices.FlavorID,
		MasterFlavorID:    choices.FlavorID,
		ExternalNetworkID: choices.ExternalNetworkID,
		NetworkDriver:     "calico",
	}
	ct, err := clustertemplates.Create(client, opts).Extract()
	th.AssertNoErr(t, err)
	defer clustertemplates.Delete(client, ct.ID)
	th.AssertEquals(t, "calico", ct.NetworkDriver)
	clusterTemplateID := ct.ID
	t.Logf(clusterTemplateID)

	// Get a cluster
	ct_get, err := clustertemplates.Get(client, ct.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "k8s", ct_get.Name)
	th.AssertEquals(t, "kubernetes", ct_get.COE)

	// Delete cluster
	clustertemplates.Delete(client, ct.ID)
	th.AssertNoErr(t, err)
	ct_delete, err := clustertemplates.Get(client, ct.ID).Extract()
	th.AssertEquals(t, (*clustertemplates.ClusterTemplate)(nil), ct_delete)

	// List clustertemplates
	pager := clustertemplates.List(client, clustertemplates.ListOpts{Limit: 1})
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page ---")
		clusterTemplateList, err := clustertemplates.ExtractClusterTemplates(page)
		th.AssertNoErr(t, err)

		for _, m := range clusterTemplateList {
			t.Logf("ClusterTemplate: ID [%s] Name [%s] COE [%s] Flavor [%s] Image [%s]",
				m.ID, m.Name, m.COE, m.FlavorID, m.ImageID)
		}
		return true, nil
	})
	th.CheckNoErr(t, err)
}
