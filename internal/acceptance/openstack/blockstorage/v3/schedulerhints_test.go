//go:build acceptance || blockstorage
// +build acceptance blockstorage

package v3

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSchedulerHints(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volumeName := tools.RandomString("ACPTTEST", 16)
	createOpts := volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}

	volume1, err := volumes.Create(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx, client, volume1.ID, "available")
	th.AssertNoErr(t, err)
	defer volumes.Delete(context.TODO(), client, volume1.ID, volumes.DeleteOpts{})

	volumeName = tools.RandomString("ACPTTEST", 16)
	schedulerHints := volumes.SchedulerHints{
		SameHost: []string{
			volume1.ID,
		},
	}
	createOpts = volumes.CreateOpts{
		Size:           1,
		Name:           volumeName,
		SchedulerHints: schedulerHints,
	}

	volume2, err := volumes.Create(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)

	ctx2, cancel2 := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel2()

	err = volumes.WaitForStatus(ctx2, client, volume2.ID, "available")
	th.AssertNoErr(t, err)

	err = volumes.Delete(context.TODO(), client, volume2.ID, volumes.DeleteOpts{}).ExtractErr()
	th.AssertNoErr(t, err)
}
