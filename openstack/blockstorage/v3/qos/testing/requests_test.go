package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/qos"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := qos.CreateOpts{
		Name:     "qos-001",
		Consumer: qos.ConsumerFront,
		Specs: map[string]string{
			"read_iops_sec": "20000",
		},
	}
	actual, err := qos.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &createQoSExpected, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := qos.Delete(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", qos.DeleteOpts{})
	th.AssertNoErr(t, res.Err)
}
