package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockCreateResponse(t, fakeServer)

	options := &shares.CreateOpts{Size: 1, Name: "my_test_share", ShareProto: "NFS", SnapshotID: "70bfbebc-d3ff-4528-8bbb-58422daa280b"}
	_, err := shares.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()

	if err == nil {
		t.Fatal("Expected error")
	}

	detailedErr := errors.ErrorDetails{}
	e := errors.ExtractErrorInto(err, &detailedErr)
	th.AssertNoErr(t, e)

	for k, msg := range detailedErr {
		th.AssertEquals(t, k, "itemNotFound")
		th.AssertEquals(t, msg.Code, 404)
		th.AssertEquals(t, msg.Message, "ShareSnapshotNotFound: Snapshot 70bfbebc-d3ff-4528-8bbb-58422daa280b could not be found.")
	}
}
