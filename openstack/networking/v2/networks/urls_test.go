package networks

import (
	"testing"

	"github.com/rackspace/gophercloud"
)

const Endpoint = "http://localhost:57909/"

func EndpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: Endpoint}
}

func TestAPIVersionsURL(t *testing.T) {
	actual := APIVersionsURL(EndpointClient())
	expected := Endpoint
	if expected != actual {
		t.Errorf("[%s] does not match expected [%s]", actual, expected)
	}
}

func TestAPIInfoURL(t *testing.T) {
	actual := APIInfoURL(EndpointClient(), "v2.0")
	expected := Endpoint + "v2.0/"
	if expected != actual {
		t.Fatalf("[%s] does not match expected [%s]", actual, expected)
	}
}

func TestExtensionURL(t *testing.T) {
	actual := ExtensionURL(EndpointClient(), "agent")
	expected := Endpoint + "v2.0/extensions/agent"
	if expected != actual {
		t.Fatalf("[%s] does not match expected [%s]", actual, expected)
	}
}

func TestNetworkURL(t *testing.T) {
	actual := NetworkURL(EndpointClient(), "foo")
	expected := Endpoint + "v2.0/networks/foo"
	if expected != actual {
		t.Fatalf("[%s] does not match expected [%s]", actual, expected)
	}
}
