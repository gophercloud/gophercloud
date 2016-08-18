// +build acceptance containerorchestration

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/baymodels"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestBayModelCRUDOperations(t *testing.T) {
	Setup(t)
	defer Teardown()

	// List baymodels
	pager := baymodels.List(Client, baymodels.ListOpts{Limit: 1})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page ---")

		bayModelList, err := baymodels.ExtractBayModels(page)
		th.AssertNoErr(t, err)

		for _, m := range bayModelList {
			t.Logf("BayModel: ID [%s] Name [%s] COE [%s] Flavor [%s] Image [%s]",
				m.ID, m.Name, m.COE, m.FlavorID, m.ImageID)
		}

		return true, nil
	})
	th.CheckNoErr(t, err)

	// Get a baymodel
	baymodelID := "5b793604-fc76-4886-a834-ed522812cdcb"
	b, err := baymodels.Get(Client, baymodelID).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, baymodelID, b.ID)
	th.AssertEquals(t, "k8sbaymodel-2", b.Name)
	th.AssertEquals(t, "kubernetes", b.COE)
	th.AssertEquals(t, b.ImageID, "fedora-atomic-latest")
	th.AssertEquals(t, b.FlavorID, "m1.small")
}
