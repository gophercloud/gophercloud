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

func TestGetURL(t *testing.T) {
	actual := GetURL(EndpointClient(), "foo")
	expected := Endpoint + "v2.0/networks/foo"
	th.AssertEquals(t, expected, actual)
}

func TestCreateURL(t *testing.T) {
	actual := CreateURL(EndpointClient())
	expected := Endpoint + "v2.0/networks"
	th.AssertEquals(t, expected, actual)
}

func TestListURL(t *testing.T) {
	actual := CreateURL(EndpointClient())
	expected := Endpoint + "v2.0/networks"
	th.AssertEquals(t, expected, actual)
}

func TestDeleteURL(t *testing.T) {
	actual := DeleteURL(EndpointClient(), "foo")
	expected := Endpoint + "v2.0/networks/foo"
	th.AssertEquals(t, expected, actual)
}
