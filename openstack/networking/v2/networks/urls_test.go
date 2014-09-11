package networks

import (
	"testing"

	"github.com/rackspace/gophercloud"
)

const Endpoint = "http://localhost:57909/v3/"

func TestAPIVersionsURL(t *testing.T) {
	client := gophercloud.ServiceClient{Endpoint: Endpoint}
	actual := APIVersionsURL(&client)
	expected := Endpoint
	if expected != actual {
		t.Errorf("[%s] does not match expected [%s]", actual, expected)
	}
}
