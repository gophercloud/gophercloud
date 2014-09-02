package services

import (
	"testing"

	"github.com/rackspace/gophercloud"
)

func TestGetListURL(t *testing.T) {
	client := gophercloud.ServiceClient{Endpoint: "http://localhost:5000/v3/"}
	url := getListURL(&client)
	if url != "http://localhost:5000/v3/services" {
		t.Errorf("Unexpected list URL generated: [%s]", url)
	}
}

func TestGetServiceURL(t *testing.T) {
	client := gophercloud.ServiceClient{Endpoint: "http://localhost:5000/v3/"}
	url := getServiceURL(&client, "1234")
	if url != "http://localhost:5000/v3/services/1234" {
		t.Errorf("Unexpected service URL generated: [%s]", url)
	}
}
