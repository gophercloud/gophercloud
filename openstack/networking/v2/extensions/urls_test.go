package extensions

import (
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const endpoint = "http://localhost:57909/"

func endpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: endpoint}
}

func TestExtensionURL(t *testing.T) {
	actual := extensionURL(endpointClient(), "agent")
	expected := endpoint + "v2.0/extensions/agent"
	th.AssertEquals(t, expected, actual)
}

func TestListExtensionURL(t *testing.T) {
	actual := listExtensionURL(endpointClient())
	expected := endpoint + "v2.0/extensions"
	th.AssertEquals(t, expected, actual)
}
