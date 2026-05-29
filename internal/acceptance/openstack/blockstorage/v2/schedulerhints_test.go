//go:build acceptance || blockstorage || schedulerhints

package v2

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSchedulerHints(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/ocata")
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV2Client()
	th.AssertNoErr(t, err)

	volumeName := tools.RandomString("ACPTTEST", 16)
	createOpts := volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}

	volume1, err := volumes.Create(context.TODO(), client, createOpts, nil).Extract()
	th.AssertNoErr(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx, client, volume1.ID, "available")
	th.AssertNoErr(t, err)
	defer volumes.Delete(context.TODO(), client, volume1.ID, volumes.DeleteOpts{})

	volumeName = tools.RandomString("ACPTTEST", 16)
	createOpts = volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}
	schedulerHintOpts := volumes.SchedulerHintOpts{
		SameHost: []string{
			volume1.ID,
		},
	}

	volume2, err := volumes.Create(context.TODO(), client, createOpts, schedulerHintOpts).Extract()
	th.AssertNoErr(t, err)

	ctx2, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx2, client, volume2.ID, "available")
	th.AssertNoErr(t, err)

	err = volumes.Delete(context.TODO(), client, volume2.ID, volumes.DeleteOpts{}).ExtractErr()
	th.AssertNoErr(t, err)
}
