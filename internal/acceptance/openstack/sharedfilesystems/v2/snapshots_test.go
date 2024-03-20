//go:build acceptance || sharedfilesystems || snapshots

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/snapshots"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// 2.7 is required for a /v2/snapshots/XXX/action URL support
// otherwise we need to set "X-OpenStack-Manila-API-Experimental: true"
const snapshotsPathMicroversion = "2.7"

func TestSnapshotCreate(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	snapshot, err := CreateSnapshot(t, client, share.ID)
	if err != nil {
		t.Fatalf("Unable to create a snapshot: %v", err)
	}

	defer DeleteSnapshot(t, client, snapshot)

	created, err := snapshots.Get(context.TODO(), client, snapshot.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve a snapshot: %v", err)
	}

	tools.PrintResource(t, created)
}

func TestSnapshotUpdate(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create share: %v", err)
	}

	defer DeleteShare(t, client, share)

	snapshot, err := CreateSnapshot(t, client, share.ID)
	if err != nil {
		t.Fatalf("Unable to create a snapshot: %v", err)
	}

	defer DeleteSnapshot(t, client, snapshot)

	expectedSnapshot, err := snapshots.Get(context.TODO(), client, snapshot.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve snapshot: %v", err)
	}

	name := "NewName"
	description := ""
	options := snapshots.UpdateOpts{
		DisplayName:        &name,
		DisplayDescription: &description,
	}

	expectedSnapshot.Name = name
	expectedSnapshot.Description = description

	_, err = snapshots.Update(context.TODO(), client, snapshot.ID, options).Extract()
	if err != nil {
		t.Errorf("Unable to update snapshot: %v", err)
	}

	updatedSnapshot, err := snapshots.Get(context.TODO(), client, snapshot.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve snapshot: %v", err)
	}

	tools.PrintResource(t, snapshot)

	th.CheckDeepEquals(t, expectedSnapshot, updatedSnapshot)
}

func TestSnapshotListDetail(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	snapshot, err := CreateSnapshot(t, client, share.ID)
	if err != nil {
		t.Fatalf("Unable to create a snapshot: %v", err)
	}

	defer DeleteSnapshot(t, client, snapshot)

	ss, err := ListSnapshots(t, client)
	if err != nil {
		t.Fatalf("Unable to list snapshots: %v", err)
	}

	for i := range ss {
		tools.PrintResource(t, &ss[i])
	}
}

func TestSnapshotResetStatus(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = snapshotsPathMicroversion

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	snapshot, err := CreateSnapshot(t, client, share.ID)
	if err != nil {
		t.Fatalf("Unable to create a snapshot: %v", err)
	}

	defer DeleteSnapshot(t, client, snapshot)

	resetStatusOpts := &snapshots.ResetStatusOpts{
		Status: "error",
	}
	err = snapshots.ResetStatus(context.TODO(), client, snapshot.ID, resetStatusOpts).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to reset a snapshot status: %v", err)
	}

	err = waitForSnapshotStatus(t, client, snapshot.ID, "error")
	if err != nil {
		t.Fatalf("Snapshot status error: %v", err)
	}

	t.Logf("Snapshot %s status successfuly reset", snapshot.ID)
}

func TestSnapshotForceDelete(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}
	client.Microversion = snapshotsPathMicroversion

	share, err := CreateShare(t, client)
	if err != nil {
		t.Fatalf("Unable to create a share: %v", err)
	}

	defer DeleteShare(t, client, share)

	snapshot, err := CreateSnapshot(t, client, share.ID)
	if err != nil {
		t.Fatalf("Unable to create a snapshot: %v", err)
	}

	defer DeleteSnapshot(t, client, snapshot)

	err = snapshots.ForceDelete(context.TODO(), client, snapshot.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to force delete a snapshot: %v", err)
	}

	err = waitForSnapshotStatus(t, client, snapshot.ID, "deleted")
	if err != nil {
		t.Fatalf("Snapshot status error: %v", err)
	}

	t.Logf("Snapshot %s was successfuly deleted", snapshot.ID)
}
