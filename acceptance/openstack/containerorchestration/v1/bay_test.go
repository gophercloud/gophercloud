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

	// Create a bay
	bayModelID := "5b793604-fc76-4886-a834-ed522812cdcb"
	b, err := bays.Create(Client, bays.CreateOpts{BayModelID: bayModelID}).Extract()
	th.AssertNoErr(t, err)
	defer bays.Delete(Client, b.ID)
	th.AssertEquals(t, b.Status, "CREATE_IN_PROGRESS")
	th.AssertEquals(t, b.BayModelID, bayModelID)
	th.AssertEquals(t, b.Masters, 1)
	th.AssertEquals(t, b.Nodes, 1)
	bayID := b.ID
	bayName := b.Name

	// List bays
	pager := bays.List(Client, bays.ListOpts{Limit: 1})
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page ---")

		bayList, err := bays.ExtractBays(page)
		th.AssertNoErr(t, err)

		for _, b := range bayList {
			t.Logf("Bay: ID [%s] Name [%s] Status [%s] Nodes [%s]",
				b.ID, b.Name, b.Status, strconv.Itoa(b.Nodes))
		}

		return true, nil
	})
	th.CheckNoErr(t, err)

	// Get a bay
	if bayID == "" {
		t.Fatalf("In order to retrieve a bay, the BayID must be set")
	}
	b, err = bays.Get(Client, bayID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, b.Name, bayName)
	th.AssertEquals(t, b.BayModelID, bayModelID)
	th.AssertEquals(t, b.Masters, 1)
	th.AssertEquals(t, b.Nodes, 1)
}
