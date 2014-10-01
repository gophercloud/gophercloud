package volumes

import (
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const endpoint = "http://localhost:57909"

func endpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: endpoint}
}

func TestVolumesURL(t *testing.T) {
	actual := volumesURL(endpointClient())
	expected := endpoint + "volumes"
	th.AssertEquals(t, expected, actual)
}

func TestVolumeURL(t *testing.T) {
	actual := volumeURL(endpointClient(), "foo")
	expected := endpoint + "volumes/foo"
	th.AssertEquals(t, expected, actual)
}
