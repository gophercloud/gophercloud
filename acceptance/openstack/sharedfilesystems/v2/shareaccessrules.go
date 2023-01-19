package v2

import (
	"testing"

	"github.com/bizflycloud/gophercloud"
	"github.com/bizflycloud/gophercloud/openstack/sharedfilesystems/v2/shareaccessrules"
)

func ShareAccessRuleGet(t *testing.T, client *gophercloud.ServiceClient, accessID string) (*shareaccessrules.ShareAccess, error) {
	accessRule, err := shareaccessrules.Get(client, accessID).Extract()
	if err != nil {
		t.Logf("Failed to get share access rule %s: %v", accessID, err)
		return nil, err
	}

	return accessRule, nil
}
