// +build acceptance dns transfers

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	identity "github.com/gophercloud/gophercloud/acceptance/openstack/identity/v3"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	transferAccepts "github.com/gophercloud/gophercloud/openstack/dns/v2/transfer/accept"
	transferRequests "github.com/gophercloud/gophercloud/openstack/dns/v2/transfer/request"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestTransferRequestCRUD(t *testing.T) {
	// Create new Zone
	clients.RequireDNS(t)

	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	zone, err := CreateZone(t, client)
	th.AssertNoErr(t, err)
	defer DeleteZone(t, client, zone)

	// Create transfers request to new tenant
	transferRequest, err := CreateTransferRequest(t, client, zone, "123")
	th.AssertNoErr(t, err)
	defer DeleteTransferRequest(t, client, transferRequest)

	allTransferRequestsPages, err := transferRequests.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allTransferRequests, err := transferRequests.ExtractTransferRequests(allTransferRequestsPages)
	th.AssertNoErr(t, err)

	var foundRequest bool
	for _, tr := range allTransferRequests {
		tools.PrintResource(t, &tr)
		if transferRequest.ZoneID == tr.ZoneID {
			foundRequest = true
		}
	}
	th.AssertEquals(t, foundRequest, true)

	description := "new description"
	updateOpts := transferRequests.UpdateOpts{
		Description: description,
	}

	newTransferRequest, err := transferRequests.Update(client, transferRequest.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, &newTransferRequest)
	th.AssertEquals(t, newTransferRequest.Description, description)
}

func TestTransferRequestAccept(t *testing.T) {
	// Create new project
	clients.RequireAdmin(t)

	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := identity.CreateProject(t, identityClient, nil)
	th.AssertNoErr(t, err)
	defer identity.DeleteProject(t, identityClient, project.ID)

	// Create new Zone
	clients.RequireDNS(t)

	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	zone, err := CreateZone(t, client)
	th.AssertNoErr(t, err)
	defer DeleteZone(t, client, zone)

	// Create transfers request to new tenant
	transferRequest, err := CreateTransferRequest(t, client, zone, project.ID)

	// Accept Zone Transfer Request
	transferAccept, err := CreateTransferAccept(t, client, transferRequest.ID, transferRequest.Key)

	allTransferAcceptsPages, err := transferAccepts.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allTransferAccepts, err := transferAccepts.ExtractTransferAccepts(allTransferAcceptsPages)
	th.AssertNoErr(t, err)

	var foundAccept bool
	for _, ta := range allTransferAccepts {
		tools.PrintResource(t, &ta)
		if transferAccept.ZoneID == ta.ZoneID {
			foundAccept = true
		}
	}
	th.AssertEquals(t, foundAccept, true)
}
