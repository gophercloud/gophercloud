package apiversions

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
	th.AssertEquals(t, expected, actual)
}

func TestAPIInfoURL(t *testing.T) {
	actual := APIInfoURL(EndpointClient(), "v2.0")
	expected := Endpoint + "v2.0/"
	th.AssertEquals(t, expected, actual)
}
