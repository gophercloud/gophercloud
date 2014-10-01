package objects

import (
	"testing"
	"github.com/rackspace/gophercloud"
)

func TestContainerURL(t *testing.T) {
	client := gophercloud.ServiceClient{
		Endpoint: "http://localhost:5000/v1/",
	}
	expected := "http://localhost:5000/v1/testContainer"
	actual := containerURL(&client, "testContainer")
	if actual != expected {
		t.Errorf("Unexpected service URL generated: %s", actual)
	}
}

func TestObjectURL(t *testing.T) {
	client := gophercloud.ServiceClient{
		Endpoint: "http://localhost:5000/v1/",
	}
	expected := "http://localhost:5000/v1/testContainer/testObject"
	actual := objectURL(&client, "testContainer", "testObject")
	if actual != expected {
		t.Errorf("Unexpected service URL generated: %s", actual)
	}
}
