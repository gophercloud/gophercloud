package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/manageablevolumes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestManageExisting(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockManageExistingResponse(t)

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
	n, err := manageablevolumes.ManageExisting(context.TODO(), client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Host, "host@lvm#LVM")
	th.AssertEquals(t, n.Name, "New Volume")
	th.AssertEquals(t, n.AvailabilityZone, "nova")
	th.AssertEquals(t, n.Description, "Volume imported from existingLV")
	th.AssertEquals(t, n.Bootable, "true")
	th.AssertDeepEquals(t, n.Metadata, map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	th.AssertEquals(t, n.ID, "23cf872b-c781-4cd4-847d-5f2ec8cbd91c")
}
