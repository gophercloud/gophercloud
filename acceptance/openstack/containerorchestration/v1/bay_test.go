// +build acceptance containerorchestration

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"strconv"
)

func TestBayCRUDOperations(t *testing.T) {
	Setup(t)
	defer Teardown()

	// TODO: Once Create and List are complete, remove the hard-coded bay id
	bayID := "a56a6cd8-0779-461b-b1eb-26cec904284a"

	// List bays
	pager := bays.List(Client, bays.ListOpts{Limit: 1})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page ---")

		bayList, err := bays.ExtractBays(page)
		th.AssertNoErr(t, err)

		for _, n := range bayList {
			t.Logf("Bay: ID [%s] Name [%s] Status [%s] Nodes [%s]",
				n.ID, n.Name, n.Status, strconv.Itoa(n.Nodes))
		}

		return true, nil
	})
	th.CheckNoErr(t, err)

	// Get a bay
	if bayID == "" {
		t.Fatalf("In order to retrieve a bay, the BayID must be set")
	}
	bay, err := bays.Get(Client, bayID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bay.Status, "CREATE_COMPLETE")
	th.AssertEquals(t, bay.Name, "k8sbay")
	th.AssertEquals(t, bay.Nodes, 1)
}
