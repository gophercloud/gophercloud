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

func TestVolumeActions(t *testing.T) {
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

	_, err = volumeactions.Attach(client, cv.ID, &volumeactions.AttachOpts{
		MountPoint:   "/mnt",
		Mode:         "rw",
		InstanceUUID: "50902f4f-a974-46a0-85e9-7efc5e22dfdd",
	}).Extract()
	th.AssertNoErr(t, err)

	err = volumes.WaitForStatus(client, cv.ID, "in-use", 60)
	th.AssertNoErr(t, err)

	_, err = volumeactions.Detach(client, cv.ID).Extract()
	th.AssertNoErr(t, err)
}
