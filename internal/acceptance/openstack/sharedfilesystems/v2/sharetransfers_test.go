//go:build acceptance || sharedfilesystems || sharetransfers

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetransfers"
)

// minimal microversion for the share transfers
const shareTransfersMicroversion = "2.77"

func TestTransferRequestCRUD(t *testing.T) {
	clients.SkipReleasesBelow(t, "master")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = shareTransfersMicroversion

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	// Create transfers request to a new tenant
	trName := "123"
	transferRequest, err := CreateTransferRequest(t, client, share, trName)
	th.AssertNoErr(t, err)
	defer DeleteTransferRequest(t, client, transferRequest)

	// list transfer requests
	allTransferRequestsPages, err := sharetransfers.ListDetail(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allTransferRequests, err := sharetransfers.ExtractTransfers(allTransferRequestsPages)
	th.AssertNoErr(t, err)

	// finding the transfer request
	var foundRequest bool
	for _, tr := range allTransferRequests {
		tools.PrintResource(t, &tr)
		if tr.ResourceID == share.ID && tr.Name == trName && !tr.Accepted {
			foundRequest = true
		}
	}
	th.AssertEquals(t, foundRequest, true)

	// checking get
	tr, err := sharetransfers.Get(context.TODO(), client, transferRequest.ID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, transferRequest.ID == tr.ID, true)

	// Accept Share Transfer Request
	err = AcceptTransfer(t, client, transferRequest)
	th.AssertNoErr(t, err)
}
