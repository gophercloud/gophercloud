//go:build acceptance || messaging || claims

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/messaging/v2/claims"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCRUDClaim(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734c"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	createdQueueName, err := CreateQueue(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQueue(t, client, createdQueueName)

	clientID = "3381af92-2b9e-11e3-b191-71861300734d"

	client, err = clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}
	for i := 0; i < 3; i++ {
		_, err := CreateMessage(t, client, createdQueueName)
		th.AssertNoErr(t, err)
	}

	claimedMessages, err := CreateClaim(t, client, createdQueueName)
	th.AssertNoErr(t, err)
	claimIDs, _ := ExtractIDs(claimedMessages)

	tools.PrintResource(t, claimedMessages)

	updateOpts := claims.UpdateOpts{
		TTL:   600,
		Grace: 500,
	}

	for _, claimID := range claimIDs {
		t.Logf("Attempting to update claim: %s", claimID)
		err := claims.Update(context.TODO(), client, createdQueueName, claimID, updateOpts).ExtractErr()
		if err != nil {
			t.Fatalf("Unable to update claim %s: %v", claimID, err)
		} else {
			t.Logf("Successfully updated claim: %s", claimID)
		}

		updatedClaim, err := GetClaim(t, client, createdQueueName, claimID)
		if err != nil {
			t.Fatalf("Unable to retrieve claim %s: %v", claimID, err)
		}

		tools.PrintResource(t, updatedClaim)
		err = DeleteClaim(t, client, createdQueueName, claimID)
		th.AssertNoErr(t, err)
	}
}
