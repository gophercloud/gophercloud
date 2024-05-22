package v2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/replicas"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
)

// CreateReplica will create a replica from shareID. An error will be returned
// if the replica could not be created.
func CreateReplica(t *testing.T, client *gophercloud.ServiceClient, share *shares.Share) (*replicas.Replica, error) {
	createOpts := replicas.CreateOpts{
		ShareID:          share.ID,
		AvailabilityZone: share.AvailabilityZone,
	}

	replica, err := replicas.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		t.Logf("Failed to create replica")
		return nil, err
	}

	err = waitForReplicaStatus(t, client, replica.ID, "available")
	if err != nil {
		t.Logf("Failed to get %s replica status", replica.ID)
		DeleteReplica(t, client, replica)
		return replica, err
	}

	return replica, nil
}

// DeleteReplica will delete a replica. A fatal error will occur if the replica
// failed to be deleted. This works best when used as a deferred function.
func DeleteReplica(t *testing.T, client *gophercloud.ServiceClient, replica *replicas.Replica) {
	err := replicas.Delete(context.TODO(), client, replica.ID).ExtractErr()
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return
		}
		t.Errorf("Unable to delete replica %s: %v", replica.ID, err)
	}

	err = waitForReplicaStatus(t, client, replica.ID, "deleted")
	if err != nil {
		t.Errorf("Failed to wait for 'deleted' status for %s replica: %v", replica.ID, err)
	} else {
		t.Logf("Deleted replica: %s", replica.ID)
	}
}

// ListShareReplicas lists all replicas that belong to shareID.
// An error will be returned if the replicas could not be listed..
func ListShareReplicas(t *testing.T, client *gophercloud.ServiceClient, shareID string) ([]replicas.Replica, error) {
	opts := replicas.ListOpts{
		ShareID: shareID,
	}
	pages, err := replicas.List(client, opts).AllPages(context.TODO())
	if err != nil {
		t.Errorf("Unable to list %q share replicas: %v", shareID, err)
	}

	return replicas.ExtractReplicas(pages)
}

func waitForReplicaStatus(t *testing.T, c *gophercloud.ServiceClient, id, status string) error {
	var current *replicas.Replica

	err := tools.WaitFor(func(ctx context.Context) (bool, error) {
		var err error

		current, err = replicas.Get(ctx, c, id).Extract()
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
			return fmt.Errorf("Replica status is '%s' and unable to get manila messages: %s", err, mErr)
		}
	}

	return err
}

func waitForReplicaState(t *testing.T, c *gophercloud.ServiceClient, id, state string) error {
	var current *replicas.Replica

	err := tools.WaitFor(func(ctx context.Context) (bool, error) {
		var err error

		current, err = replicas.Get(ctx, c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.State == state {
			return true, nil
		}

		if strings.Contains(current.State, "error") {
			return true, fmt.Errorf("An error occurred, wrong state: %s", current.State)
		}

		return false, nil
	})

	if err != nil {
		mErr := PrintMessages(t, c, id)
		if mErr != nil {
			return fmt.Errorf("Replica state is '%s' and unable to get manila messages: %s", err, mErr)
		}
	}

	return err
}
