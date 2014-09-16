package networks

import (
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const Endpoint = "http://localhost:57909/"

func EndpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: Endpoint}
}

func TestExtensionURL(t *testing.T) {
	actual := ExtensionURL(EndpointClient(), "agent")
	expected := Endpoint + "v2.0/extensions/agent"
	th.AssertEquals(t, expected, actual)
}
