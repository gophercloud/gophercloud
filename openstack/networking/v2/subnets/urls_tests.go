package subnets

import (
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const Endpoint = "http://localhost:57909/"

func EndpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: Endpoint}
}

func TestListURL(t *testing.T) {
	actual := ListURL(EndpointClient())
	expected := Endpoint + "v2.0/subnets"
	th.AssertEquals(t, expected, actual)
}

func TestGetURL(t *testing.T) {
	actual := GetURL(EndpointClient(), "foo")
	expected := Endpoint + "v2.0/subnets/foo"
	th.AssertEquals(t, expected, actual)
}
