package v2

import (
	"encoding/json"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"testing"
)

// CreateShare will create a share with a name, and a size of 1Gb. An
// error will be returned if the share could not be created
func CreateShare(t *testing.T, client *gophercloud.ServiceClient) (*shares.Share, error) {
	if testing.Short() {
		t.Skip("Skipping test that requres share creation in short mode.")
	}

	// TODO: Need a way to create a share network, which is attached
	// to the private neutron network.
	createOpts := shares.CreateOpts{
		Size:           1,
		Name:           "My Test Share",
		ShareProto:     "NFS",
		ShareNetworkID: "cc92a11e-a108-4a25-96d3-ef9281e93e4e",
	}

	share, err := shares.Create(client, createOpts).Extract()
	if err != nil {
		return share, err
	}

	err = shares.WaitForStatus(client, share.ID, "available", 60)
	if err != nil {
		return share, err
	}

	return share, nil
}

// DeleteShare will delete a share. A fatal error will occur if the share
// failed to be deleted. This works best when used as a deferred function.
func DeleteShare(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share) {
	err := shares.Delete(client, share.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete share %s: %v", share.ID, err)
	}

	t.Logf("Deleted share: %s", share.ID)
}

// PrintShare prints some information of the share
func PrintShare(t *testing.T, share *shares.Share) {
	asJSON, err := json.MarshalIndent(share, "", " ")
	if err != nil {
		t.Logf("Cannot print the contents of %s", share.ID)
	}

	t.Logf("Share %s", string(asJSON))
}
