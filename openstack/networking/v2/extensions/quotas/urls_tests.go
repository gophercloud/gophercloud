package quotas

import (
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const endpoint = "http://localhost:57909/"

func endpointClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{Endpoint: endpoint}
}

func TestRootURL(t *testing.T) {
	expected := endpoint + "v2.0/quotas"
	actual := rootURL(endpointClient())
	th.AssertEquals(t, expected, actual)
}
