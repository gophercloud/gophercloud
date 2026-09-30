//go:build acceptance || sharedfilesystems || shareaccessrules

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shareaccessrules"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
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

func TestShareAccessRulesMetadata(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	// Access rule metadata requires microversion 2.45 or later.
	client.Microversion = "2.45"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	addedAccessRight, err := shares.GrantAccess(context.TODO(), client, share.ID, shares.GrantAccessOpts{
		AccessType:  "ip",
		AccessTo:    "0.0.0.0/32",
		AccessLevel: "ro",
		Metadata: map[string]string{
			"key1": "value1",
		},
	}).Extract()
	if err != nil {
		t.Fatalf("Unable to grant access to share %s: %v", share.ID, err)
	}

	th.AssertDeepEquals(t, map[string]string{"key1": "value1"}, addedAccessRight.Metadata)

	addedShareAccess := AccessRightToShareAccess(addedAccessRight)

	if err = WaitForShareAccessRule(t, client, addedShareAccess, "active"); err != nil {
		t.Fatalf("Unable to wait for share access rule to achieve 'active' state: %v", err)
	}

	tools.PrintResource(t, addedAccessRight)

	updated, err := shareaccessrules.UpdateMetadata(context.TODO(), client, addedAccessRight.ID, shareaccessrules.UpdateMetadataOpts{
		Metadata: map[string]string{
			"key2": "value2",
		},
	}).Extract()
	if err != nil {
		t.Fatalf("Unable to update metadata for share access rule %s: %v", addedAccessRight.ID, err)
	}

	th.AssertDeepEquals(t, map[string]string{"key2": "value2"}, updated)

	accessRule, err := ShareAccessRuleGet(t, client, addedAccessRight.ID)
	if err != nil {
		t.Fatalf("Unable to get share access rule for share %s: %v", share.ID, err)
	}

	tools.PrintResource(t, accessRule)

	// key1 was set at creation time and must still be present; key2 was
	// added via UpdateMetadata, which only touches the keys it is given.
	th.AssertDeepEquals(t, map[string]string{"key1": "value1", "key2": "value2"}, accessRule.Metadata)

	err = shareaccessrules.DeleteMetadatum(context.TODO(), client, addedAccessRight.ID, "key1").ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete metadatum for share access rule %s: %v", addedAccessRight.ID, err)
	}

	accessRule, err = ShareAccessRuleGet(t, client, addedAccessRight.ID)
	if err != nil {
		t.Fatalf("Unable to get share access rule for share %s: %v", share.ID, err)
	}

	th.AssertDeepEquals(t, map[string]string{"key2": "value2"}, accessRule.Metadata)
}
