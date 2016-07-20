// +build acceptance containerorchestration

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestBayCRUDOperations(t *testing.T) {
	Setup(t)
	defer Teardown()

	// TODO: Once Create and List are complete, remove the hard-coded bay id
	bayID := "a56a6cd8-0779-461b-b1eb-26cec904284a"

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
