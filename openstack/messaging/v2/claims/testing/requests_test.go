package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/messaging/v2/claims"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer)

	createOpts := claims.CreateOpts{
		TTL:   3600,
		Grace: 3600,
		Limit: 10,
	}

	actual, err := claims.Create(context.TODO(), client.ServiceClient(fakeServer), QueueName, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, CreatedClaim, actual)
}

func TestCreateNoContent(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateNoContent(t, fakeServer)

	createOpts := claims.CreateOpts{
		TTL:   3600,
		Grace: 3600,
		Limit: 10,
	}

	actual, err := claims.Create(context.TODO(), client.ServiceClient(fakeServer), QueueName, createOpts).Extract()
	th.AssertNoErr(t, err)
	var expected []claims.Messages
	th.CheckDeepEquals(t, expected, actual)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)

	actual, err := claims.Get(context.TODO(), client.ServiceClient(fakeServer), QueueName, ClaimID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstClaim, actual)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	updateOpts := claims.UpdateOpts{
		Grace: 1600,
		TTL:   1200,
	}

	err := claims.Update(context.TODO(), client.ServiceClient(fakeServer), QueueName, ClaimID, updateOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	err := claims.Delete(context.TODO(), client.ServiceClient(fakeServer), QueueName, ClaimID).ExtractErr()
	th.AssertNoErr(t, err)
}
