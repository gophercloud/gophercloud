package snapshots

import (
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const endpoint = "http://localhost:57909"

func endpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: endpoint}
}

func TestSnapshotsURL(t *testing.T) {
	actual := snapshotsURL(endpointClient())
	expected := endpoint + "snapshots"
	th.AssertEquals(t, expected, actual)
}

func TestSnapshotURL(t *testing.T) {
	actual := snapshotURL(endpointClient(), "foo")
	expected := endpoint + "snapshots/foo"
	th.AssertEquals(t, expected, actual)
}
