// +build acceptance blockstorage

package extensions

import (
	"os"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v2/extensions/volumeactions"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v2/volumes"
	th "github.com/rackspace/gophercloud/testhelper"
)

func newClient(t *testing.T) (*gophercloud.ServiceClient, error) {
	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	client, err := openstack.AuthenticatedClient(ao)
	th.AssertNoErr(t, err)

	return openstack.NewBlockStorageV2(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

func TestVolumeAttach(t *testing.T) {
	client, err := newClient(t)
	th.AssertNoErr(t, err)

	cv, err := volumes.Create(client, &volumes.CreateOpts{
		Size: 1,
		Name: "blockv2-volume",
	}).Extract()
	th.AssertNoErr(t, err)
	defer func() {
		err = volumes.WaitForStatus(client, cv.ID, "available", 60)
		th.AssertNoErr(t, err)
		err = volumes.Delete(client, cv.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()

	err = volumes.WaitForStatus(client, cv.ID, "available", 60)
	th.AssertNoErr(t, err)

	instanceID := os.Getenv("OS_INSTANCE_ID")
	if instanceID == "" {
		t.Fatal("Environment variable OS_INSTANCE_ID is required")
	}

	_, err = volumeactions.Attach(client, cv.ID, &volumeactions.AttachOpts{
		MountPoint:   "/mnt",
		Mode:         "rw",
		InstanceUUID: instanceID,
	}).Extract()
	th.AssertNoErr(t, err)

	err = volumes.WaitForStatus(client, cv.ID, "in-use", 60)
	th.AssertNoErr(t, err)

	_, err = volumeactions.Detach(client, cv.ID).Extract()
	th.AssertNoErr(t, err)
}

func TestVolumeReserve(t *testing.T) {
	client, err := newClient(t)
	th.AssertNoErr(t, err)

	cv, err := volumes.Create(client, &volumes.CreateOpts{
		Size: 1,
		Name: "blockv2-volume",
	}).Extract()
	th.AssertNoErr(t, err)
	defer func() {
		err = volumes.WaitForStatus(client, cv.ID, "available", 60)
		th.AssertNoErr(t, err)
		err = volumes.Delete(client, cv.ID).ExtractErr()
		th.AssertNoErr(t, err)
	}()

	err = volumes.WaitForStatus(client, cv.ID, "available", 60)
	th.AssertNoErr(t, err)

	_, err = volumeactions.Reserve(client, cv.ID).Extract()
	th.AssertNoErr(t, err)

	err = volumes.WaitForStatus(client, cv.ID, "attaching", 60)
	th.AssertNoErr(t, err)

	_, err = volumeactions.Unreserve(client, cv.ID).Extract()
	th.AssertNoErr(t, err)

	err = volumes.WaitForStatus(client, cv.ID, "available", 60)
	th.AssertNoErr(t, err)
}
