package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/manageablevolumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestManageExisting(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockManageExistingResponse(t, fakeServer)

	options := &manageablevolumes.ManageExistingOpts{
		Host:             "host@lvm#LVM",
		Ref:              map[string]string{"source-name": "volume-73796b96-169f-4675-a5bc-73fc0f8f9a17"},
		Name:             "New Volume",
		AvailabilityZone: "nova",
		Description:      "Volume imported from existingLV",
		VolumeType:       "lvm",
		Bootable:         true,
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	n, err := manageablevolumes.ManageExisting(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "host@lvm#LVM", n.Host)
	th.AssertEquals(t, "New Volume", n.Name)
	th.AssertEquals(t, "nova", n.AvailabilityZone)
	th.AssertEquals(t, "Volume imported from existingLV", n.Description)
	th.AssertEquals(t, "true", n.Bootable)
	th.AssertDeepEquals(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, n.Metadata)
	th.AssertEquals(t, "23cf872b-c781-4cd4-847d-5f2ec8cbd91c", n.ID)
}
