package volumeTypes

import (
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const endpoint = "http://localhost:57909"

func endpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: endpoint}
}

func TestVolumeTypesURL(t *testing.T) {
	actual := volumeTypesURL(endpointClient())
	expected := endpoint + "types"
	th.AssertEquals(t, expected, actual)
}

func TestVolumeTypeURL(t *testing.T) {
	actual := volumeTypeURL(endpointClient(), "foo")
	expected := endpoint + "types/foo"
	th.AssertEquals(t, expected, actual)
}
