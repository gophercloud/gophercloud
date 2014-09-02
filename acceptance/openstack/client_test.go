// +build acceptance

package openstack

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/utils"
)

func TestAuthenticatedClient(t *testing.T) {
	// Obtain credentials from the environment.
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Fatalf("Unable to acquire credentials: %v", err)
	}

	client, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		t.Fatalf("Unable to authenticate: %v", err)
	}

	if client.TokenID == "" {
		t.Errorf("No token ID assigned to the client")
	}

	t.Logf("Client successfully acquired a token: %v", client.TokenID)
}
