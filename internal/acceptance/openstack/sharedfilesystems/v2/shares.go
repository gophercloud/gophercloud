package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/messages"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
)

// CreateShare will create a share with a name, and a size of 1Gb. An
// error will be returned if the share could not be created
func CreateShare(t *testing.T, client *gophercloud.ServiceClient, optShareType ...string) (*shares.Share, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires share creation in short mode.")
	}

	iTrue := true
	shareType := "dhss_false"
	if len(optShareType) > 0 {
		shareType = optShareType[0]
	}
	createOpts := shares.CreateOpts{
		Size:        1,
		Name:        "My Test Share",
		Description: "My Test Description",
		ShareProto:  "NFS",
		ShareType:   shareType,
		IsPublic:    &iTrue,
	}

	share, err := shares.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		t.Logf("Failed to create share")
		return nil, err
	}

	_, err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Logf("Failed to get %s share status", share.ID)
		DeleteShare(t, client, share)
		return share, err
	}

	return share, nil
}

// ListShares lists all shares that belong to this tenant's project.
// An error will be returned if the shares could not be listed..
func ListShares(t *testing.T, client *gophercloud.ServiceClient) ([]shares.Share, error) {
	r, err := shares.ListDetail(client, &shares.ListOpts{}).AllPages(context.TODO())
	if err != nil {
		return nil, err
	}

	return shares.ExtractShares(r)
}

// GrantAccess will grant access to an existing share. A fatal error will occur if
// this operation fails.
func GrantAccess(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share) (*shares.AccessRight, error) {
	return shares.GrantAccess(context.TODO(), client, share.ID, shares.GrantAccessOpts{
		AccessType:  "ip",
		AccessTo:    "0.0.0.0/32",
		AccessLevel: "ro",
	}).Extract()
}

// RevokeAccess will revoke an exisiting access of a share. A fatal error will occur
// if this operation fails.
func RevokeAccess(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share, accessRight *shares.AccessRight) error {
	return shares.RevokeAccess(context.TODO(), client, share.ID, shares.RevokeAccessOpts{
		AccessID: accessRight.ID,
	}).ExtractErr()
}

// GetAccessRightsSlice will retrieve all access rules assigned to a share.
// A fatal error will occur if this operation fails.
func GetAccessRightsSlice(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share) ([]shares.AccessRight, error) {
	return shares.ListAccessRights(context.TODO(), client, share.ID).Extract()
}

// DeleteShare will delete a share. A fatal error will occur if the share
// failed to be deleted. This works best when used as a deferred function.
func DeleteShare(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share) {
	err := shares.Delete(context.TODO(), client, share.ID).ExtractErr()
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return
		}
		t.Errorf("Unable to delete share %s: %v", share.ID, err)
	}

	_, err = waitForStatus(t, client, share.ID, "deleted")
	if err != nil {
		t.Errorf("Failed to wait for 'deleted' status for %s share: %v", share.ID, err)
	} else {
		t.Logf("Deleted share: %s", share.ID)
	}
}

// ExtendShare extends the capacity of an existing share
func ExtendShare(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share, newSize int) error {
	return shares.Extend(context.TODO(), client, share.ID, &shares.ExtendOpts{NewSize: newSize}).ExtractErr()
}

// ShrinkShare shrinks the capacity of an existing share
func ShrinkShare(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share, newSize int) error {
	return shares.Shrink(context.TODO(), client, share.ID, &shares.ShrinkOpts{NewSize: newSize}).ExtractErr()
}

func PrintMessages(t *testing.T, c *gophercloud.ServiceClient, id string) error {
	c.Microversion = "2.37"

	allPages, err := messages.List(c, messages.ListOpts{ResourceID: id}).AllPages(context.TODO())
	if err != nil {
		return fmt.Errorf("Unable to retrieve messages: %v", err)
	}

	allMessages, err := messages.ExtractMessages(allPages)
	if err != nil {
		return fmt.Errorf("Unable to extract messages: %v", err)
	}

	for _, message := range allMessages {
		tools.PrintResource(t, message)
	}

	return nil
}

func waitForStatus(t *testing.T, c *gophercloud.ServiceClient, id, status string) (*shares.Share, error) {
	var current *shares.Share

	err := tools.WaitFor(func(ctx context.Context) (bool, error) {
		var err error

		current, err = shares.Get(ctx, c, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				switch status {
				case "deleted":
					return true, nil
				default:
					return false, err
				}
			}
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		if strings.Contains(current.Status, "error") {
			return true, fmt.Errorf("An error occurred, wrong status: %s", current.Status)
		}

		return false, nil
	})

	if err != nil {
		mErr := PrintMessages(t, c, id)
		if mErr != nil {
			return current, fmt.Errorf("Share status is '%s' and unable to get manila messages: %s", err, mErr)
		}
	}

	return current, err
}
