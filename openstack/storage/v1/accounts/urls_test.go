package accounts

import (
	"testing"

	"github.com/rackspace/gophercloud"
)

func TestAccountURL(t *testing.T) {
	client := gophercloud.ServiceClient{
		Endpoint: "http://localhost:5000/v3/",
	}
	url := accountURL(&client)
	if url != "http://localhost:5000/v3/" {
		t.Errorf("Unexpected service URL generated: [%s]", url)
	}
}
