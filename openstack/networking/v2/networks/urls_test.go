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

func TestAPIVersionsURL(t *testing.T) {
	actual := APIVersionsURL(EndpointClient())
	expected := Endpoint
	th.AssertEquals(expected, actual)
}

func TestAPIInfoURL(t *testing.T) {
	actual := APIInfoURL(EndpointClient(), "v2.0")
	expected := Endpoint + "v2.0/"
	th.AssertEquals(expected, actual)
}

func TestExtensionURL(t *testing.T) {
	actual := ExtensionURL(EndpointClient(), "agent")
	expected := Endpoint + "v2.0/extensions/agent"
	th.AssertEquals(expected, actual)
}

func TestNetworkURL(t *testing.T) {
	actual := NetworkURL(EndpointClient(), "foo")
	expected := Endpoint + "v2.0/networks/foo"
	th.AssertEquals(expected, actual)
}

func TestCreateURL(t *testing.T) {
	actual := CreateURL(EndpointClient())
	expected := Endpoint + "v2.0/networks"
	th.AssertEquals(expected, actual)
}
