package containers

import (
	"testing"
	"github.com/rackspace/gophercloud"
)

func TestAccountURL(t *testing.T) {
	client := gophercloud.ServiceClient{
		Endpoint: "http://localhost:5000/v1/",
	}
	expected := "http://localhost:5000/v1/"
	actual := accountURL(&client)
	if actual != expected {
		t.Errorf("Unexpected service URL generated: [%s]", actual)
	}

}

func TestContainerURL(t *testing.T) {
	client := gophercloud.ServiceClient{
		Endpoint: "http://localhost:5000/v1/",
	}
	expected := "http://localhost:5000/v1/testContainer"
	actual := containerURL(&client, "testContainer")
	if actual != expected {
		t.Errorf("Unexpected service URL generated: [%s]", actual)
	}
}
