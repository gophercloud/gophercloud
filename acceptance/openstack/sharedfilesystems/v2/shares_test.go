// +build acceptance

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestShareCreate(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	created, err := shares.Get(client, share.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve share: %v", err)
	}
	tools.PrintResource(t, created)
}

func TestShareExportLocations(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	client.Microversion = "2.9"

	exportLocations, err := shares.ListExportLocations(client, share.ID).Extract()
	if err != nil {
		t.Errorf("Unable to list share export locations: %v", err)
	}
	tools.PrintResource(t, exportLocations)

	exportLocation, err := shares.GetExportLocation(client, share.ID, exportLocations[0].ID).Extract()
	if err != nil {
		t.Errorf("Unable to get share export location: %v", err)
	}
	tools.PrintResource(t, exportLocation)
	th.AssertEquals(t, exportLocations[0], *exportLocation)
}

func TestShareUpdate(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create share: %v", err)
	}

	defer DeleteShare(t, client, share)

	expectedShare, err := shares.Get(client, share.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve share: %v", err)
	}

	name := "NewName"
	description := ""
	iFalse := false
	options := shares.UpdateOpts{
		DisplayName:        &name,
		DisplayDescription: &description,
		IsPublic:           &iFalse,
	}

	expectedShare.Name = name
	expectedShare.Description = description
	expectedShare.IsPublic = iFalse

	_, err = shares.Update(client, share.ID, options).Extract()
	if err != nil {
		t.Errorf("Unable to update share: %v", err)
	}

	updatedShare, err := shares.Get(client, share.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve share: %v", err)
	}

	// Update time has to be set in order to get the assert equal to pass
	expectedShare.UpdatedAt = updatedShare.UpdatedAt

	tools.PrintResource(t, share)

	th.CheckDeepEquals(t, expectedShare, updatedShare)
}

func TestShareListDetail(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	ss, err := ListShares(t, client)
	if err != nil {
		t.Fatalf("Unable to list shares: %v", err)
	}

	for i := range ss {
		tools.PrintResource(t, &ss[i])
	}
}

func TestGrantAndRevokeAccess(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

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

	accessRight, err := GrantAccess(t, client, share)
	if err != nil {
		t.Fatalf("Unable to grant access: %v", err)
	}

	tools.PrintResource(t, accessRight)

	if err = RevokeAccess(t, client, share, accessRight); err != nil {
		t.Fatalf("Unable to revoke access: %v", err)
	}
}

func TestListAccessRights(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = "2.7"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	_, err = GrantAccess(t, client, share)
	if err != nil {
		t.Fatalf("Unable to grant access: %v", err)
	}

	rs, err := GetAccessRightsSlice(t, client, share)
	if err != nil {
		t.Fatalf("Unable to retrieve list of access rules for share %s: %v", share.ID, err)
	}

	if len(rs) != 1 {
		t.Fatalf("Unexpected number of access rules for share %s: got %d, expected 1", share.ID, len(rs))
	}

	t.Logf("Share %s has %d access rule(s):", share.ID, len(rs))

	for _, r := range rs {
		tools.PrintResource(t, &r)
	}
}

func TestExtendAndShrink(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = "2.7"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	err = ExtendShare(t, client, share, 2)
	if err != nil {
		t.Fatalf("Unable to extend a share: %v", err)
	}

	// We need to wait till the Extend operation is done
	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	t.Logf("Share %s successfuly extended", share.ID)

	/* disable shrinking for the LVM dhss=false
	err = ShrinkShare(t, client, share, 1)
	if err != nil {
		t.Fatalf("Unable to shrink a share: %v", err)
	}

	// We need to wait till the Shrink operation is done
	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	t.Logf("Share %s successfuly shrunk", share.ID)
	*/
}

func TestShareMetadata(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = "2.7"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	const (
		k  = "key"
		v1 = "value1"
		v2 = "value2"
	)

	checkMetadataEq := func(m map[string]string, value string) {
		if m == nil || len(m) != 1 || m[k] != value {
			t.Fatalf("Unexpected metadata contents %v", m)
		}
	}

	metadata, err := shares.SetMetadata(client, share.ID, shares.SetMetadataOpts{Metadata: map[string]string{k: v1}}).Extract()
	if err != nil {
		t.Fatalf("Unable to set share metadata: %v", err)
	}
	checkMetadataEq(metadata, v1)

	metadata, err = shares.UpdateMetadata(client, share.ID, shares.UpdateMetadataOpts{Metadata: map[string]string{k: v2}}).Extract()
	if err != nil {
		t.Fatalf("Unable to update share metadata: %v", err)
	}
	checkMetadataEq(metadata, v2)

	metadata, err = shares.GetMetadatum(client, share.ID, k).Extract()
	if err != nil {
		t.Fatalf("Unable to get share metadatum: %v", err)
	}
	checkMetadataEq(metadata, v2)

	err = shares.DeleteMetadatum(client, share.ID, k).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete share metadatum: %v", err)
	}

	metadata, err = shares.GetMetadata(client, share.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get share metadata: %v", err)
	}

	if metadata == nil || len(metadata) != 0 {
		t.Fatalf("Unexpected metadata contents %v, expected an empty map", metadata)
	}
}

func TestRevert(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = "2.27"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	snapshot, err := CreateSnapshot(t, client, share.ID)
	if err != nil {
		t.Fatalf("Unable to create a snapshot: %v", err)
	}
	defer DeleteSnapshot(t, client, snapshot)

	err = waitForSnapshotStatus(t, client, snapshot.ID, "available")
	if err != nil {
		t.Fatalf("Snapshot status error: %v", err)
	}

	revertOpts := &shares.RevertOpts{
		SnapshotID: snapshot.ID,
	}
	err = shares.Revert(client, share.ID, revertOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to revert a snapshot: %v", err)
	}

	// We need to wait till the Extend operation is done
	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	err = waitForSnapshotStatus(t, client, snapshot.ID, "available")
	if err != nil {
		t.Fatalf("Snapshot status error: %v", err)
	}

	t.Logf("Share %s successfuly reverted", share.ID)
}

func TestResetStatus(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = "2.7"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	resetStatusOpts := &shares.ResetStatusOpts{
		Status: "error",
	}
	err = shares.ResetStatus(client, share.ID, resetStatusOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to reset a share status: %v", err)
	}

	// We need to wait till the Extend operation is done
	err = waitForStatus(t, client, share.ID, "error")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	t.Logf("Share %s status successfuly reset", share.ID)
}

func TestForceDelete(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = "2.7"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	err = shares.ForceDelete(client, share.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to force delete a share: %v", err)
	}

	err = waitForStatus(t, client, share.ID, "deleted")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	t.Logf("Share %s was successfuly deleted", share.ID)
}

func TestUnmanage(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.RequireAdmin(t)

	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = "2.7"

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	err = waitForStatus(t, client, share.ID, "available")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	err = shares.Unmanage(client, share.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to unmanage a share: %v", err)
	}

	err = waitForStatus(t, client, share.ID, "deleted")
	if err != nil {
		t.Fatalf("Share status error: %v", err)
	}

	t.Logf("Share %s was successfuly unmanaged", share.ID)
}
