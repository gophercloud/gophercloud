package v2

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

// CreateShares creates `n` shares
func CreateShares(t *testing.T, client *gophercloud.ServiceClient, names []string) ([]shares.Share, error) {
	if testing.Short() {
		t.Skip("Skipping test that requres share creation in short mode.")
	}

	created := []shares.Share{}
	for _, name := range names {
		share, err := CreateShare(t, client, name)
		if err != nil {
			t.Fatalf("Failed to create a share %v", err)
		}

		created = append(created, *share)
	}

	return created, nil
}

// DeleteShares will delete all shares
func DeleteShares(t *testing.T, client *gophercloud.ServiceClient, sharesToDel []shares.Share) {
	for _, s := range sharesToDel {
		DeleteShare(t, client, &s)
	}
}

// CreateShare will create a share with a name, and a size of 1Gb. An
// error will be returned if the share could not be created
func CreateShare(t *testing.T, client *gophercloud.ServiceClient, name string) (*shares.Share, error) {
	if testing.Short() {
		t.Skip("Skipping test that requres share creation in short mode.")
	}

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatalf("Unable to fetch environment information")
	}

	createOpts := shares.CreateOpts{
		Size:           1,
		Name:           name,
		ShareProto:     "nfs",
		ShareNetworkID: choices.ShareNetworkID,
	}

	share, err := shares.Create(client, createOpts).Extract()
	if err != nil {
		return share, err
	}

	err = waitForStatus(client, share.ID, "available", 600)
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

func waitForStatus(c *gophercloud.ServiceClient, id, status string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := shares.Get(c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == "error" {
			return true, fmt.Errorf("An error occurred")
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}
