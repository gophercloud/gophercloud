//go:build acceptance
// +build acceptance

package v2

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestShareAccessRulesGet(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	client.Microversion = "2.49"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	shareAccessRight, err := GrantAccess(t, client, share)
	if err != nil {
		t.Fatalf("Unable to grant access to share %s: %v", share.ID, err)
	}

	accessRule, err := ShareAccessRuleGet(t, client, shareAccessRight.ID)
	if err != nil {
		t.Logf("Unable to get share access rule for share %s: %v", share.ID, err)
	}

	tools.PrintResource(t, accessRule)

	th.AssertEquals(t, shareAccessRight.ID, accessRule.ID)
	th.AssertEquals(t, shareAccessRight.ShareID, accessRule.ShareID)
	th.AssertEquals(t, shareAccessRight.AccessType, accessRule.AccessType)
	th.AssertEquals(t, shareAccessRight.AccessLevel, accessRule.AccessLevel)
	th.AssertEquals(t, shareAccessRight.AccessTo, accessRule.AccessTo)
	th.AssertEquals(t, shareAccessRight.AccessKey, accessRule.AccessKey)
	th.AssertEquals(t, shareAccessRight.State, accessRule.State)
}
