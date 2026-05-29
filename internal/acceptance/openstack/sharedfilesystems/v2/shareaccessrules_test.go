//go:build acceptance || sharedfilesystems || shareaccessrules

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
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

	addedAccessRight, err := GrantAccess(t, client, share)
	if err != nil {
		t.Fatalf("Unable to grant access to share %s: %v", share.ID, err)
	}

	addedShareAccess := AccessRightToShareAccess(addedAccessRight)

	accessRule, err := ShareAccessRuleGet(t, client, addedShareAccess.ID)
	if err != nil {
		t.Fatalf("Unable to get share access rule for share %s: %v", share.ID, err)
	}

	if err = WaitForShareAccessRule(t, client, accessRule, "active"); err != nil {
		t.Fatalf("Unable to wait for share access rule to achieve 'active' state: %v", err)
	}

	tools.PrintResource(t, accessRule)

	th.AssertEquals(t, addedShareAccess.ID, accessRule.ID)
	th.AssertEquals(t, addedShareAccess.AccessType, accessRule.AccessType)
	th.AssertEquals(t, addedShareAccess.AccessLevel, accessRule.AccessLevel)
	th.AssertEquals(t, addedShareAccess.AccessTo, accessRule.AccessTo)
	th.AssertEquals(t, addedShareAccess.AccessKey, accessRule.AccessKey)
	th.AssertEquals(t, share.ID, accessRule.ShareID)
	th.AssertEquals(t, "active", accessRule.State)
}

func TestShareAccessRulesList(t *testing.T) {
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

	addedAccessRight, err := GrantAccess(t, client, share)
	if err != nil {
		t.Fatalf("Unable to grant access to share %s: %v", share.ID, err)
	}

	addedShareAccess := AccessRightToShareAccess(addedAccessRight)

	if err = WaitForShareAccessRule(t, client, addedShareAccess, "active"); err != nil {
		t.Fatalf("Unable to wait for share access rule to achieve 'active' state: %v", err)
	}

	accessRules, err := ShareAccessRuleList(t, client, share.ID)
	if err != nil {
		t.Logf("Unable to list share access rules for share %s: %v", share.ID, err)
	}

	tools.PrintResource(t, accessRules)

	th.AssertEquals(t, 1, len(accessRules))

	accessRule := accessRules[0]

	if err = WaitForShareAccessRule(t, client, &accessRule, "active"); err != nil {
		t.Fatalf("Unable to wait for share access rule to achieve 'active' state: %v", err)
	}

	th.AssertEquals(t, addedShareAccess.ID, accessRule.ID)
	th.AssertEquals(t, addedShareAccess.AccessType, accessRule.AccessType)
	th.AssertEquals(t, addedShareAccess.AccessLevel, accessRule.AccessLevel)
	th.AssertEquals(t, addedShareAccess.AccessTo, accessRule.AccessTo)
	th.AssertEquals(t, addedShareAccess.AccessKey, accessRule.AccessKey)
	th.AssertEquals(t, addedShareAccess.State, accessRule.State)
}
