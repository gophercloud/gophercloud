package v2

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shareaccessrules"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func ShareAccessRuleGet(t *testing.T, client *gophercloud.ServiceClient, accessID string) (*shareaccessrules.ShareAccess, error) {
	accessRule, err := shareaccessrules.Get(client, accessID).Extract()
	if err != nil {
		t.Logf("Failed to get share access rule %s: %v", accessID, err)
		return nil, err
	}

	return accessRule, nil
}

// AccessRightToShareAccess is a helper function that converts
// shares.AccessRight into shareaccessrules.ShareAccess struct.
func AccessRightToShareAccess(accessRight *shares.AccessRight) *shareaccessrules.ShareAccess {
	return &shareaccessrules.ShareAccess{
		ShareID:     accessRight.ShareID,
		AccessType:  accessRight.AccessType,
		AccessTo:    accessRight.AccessTo,
		AccessKey:   accessRight.AccessKey,
		AccessLevel: accessRight.AccessLevel,
		State:       accessRight.State,
		ID:          accessRight.ID,
	}
}

func WaitForShareAccessRule(t *testing.T, client *gophercloud.ServiceClient, accessRule *shareaccessrules.ShareAccess, status string) error {
	if accessRule.State == status {
		return nil
	}

	return tools.WaitFor(func() (bool, error) {
		latest, err := ShareAccessRuleGet(t, client, accessRule.ID)
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return false, nil
			}

			return false, err
		}

		if latest.State == status {
			*accessRule = *latest
			return true, nil
		}

		if latest.State == "error" {
			return false, fmt.Errorf("share access rule %s for share %s is in error state", accessRule.ID, accessRule.ShareID)
		}

		return false, nil
	})
}

func ShareAccessRuleList(t *testing.T, client *gophercloud.ServiceClient, shareID string) ([]shareaccessrules.ShareAccess, error) {
	accessRules, err := shareaccessrules.List(client, shareID).Extract()
	if err != nil {
		t.Logf("Failed to list share access rules for share %s: %v", shareID, err)
		return nil, err
	}

	return accessRules, nil
}
