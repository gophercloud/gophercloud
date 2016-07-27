// +build acceptance containerorchestration

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/baymodels"
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
}
