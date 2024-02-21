package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
)

// CreateShareNetwork will create a share network with a random name. An
// error will be returned if the share network was unable to be created.
func CreateShareNetwork(t *testing.T, client *gophercloud.ServiceClient) (*sharenetworks.ShareNetwork, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires share network creation in short mode.")
	}

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		return nil, err
	}

	shareNetworkName := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create share network: %s", shareNetworkName)

	createOpts := sharenetworks.CreateOpts{
		Name:            shareNetworkName,
		NeutronNetID:    choices.NetworkID,
		NeutronSubnetID: choices.SubnetID,
		Description:     "This is a shared network",
	}

	shareNetwork, err := sharenetworks.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return shareNetwork, err
	}

	return shareNetwork, nil
}

// DeleteShareNetwork will delete a share network. An error will occur if
// the share network was unable to be deleted.
func DeleteShareNetwork(t *testing.T, client *gophercloud.ServiceClient, shareNetworkID string) {
	err := sharenetworks.Delete(context.TODO(), client, shareNetworkID).ExtractErr()
	if err != nil {
		t.Fatalf("Failed to delete share network %s: %v", shareNetworkID, err)
	}

	t.Logf("Deleted share network: %s", shareNetworkID)
}
