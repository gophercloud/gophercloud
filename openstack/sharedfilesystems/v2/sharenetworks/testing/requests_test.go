package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// Verifies that a share network can be created correctly
func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &sharenetworks.CreateOpts{
		Name:            "my_network",
		Description:     "This is my share network",
		NeutronNetID:    "998b42ee-2cee-4d36-8b95-67b5ca1f2109",
		NeutronSubnetID: "53482b62-2c84-4a53-b6ab-30d9d9800d06",
	}

	n, err := sharenetworks.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Name, "my_network")
	th.AssertEquals(t, n.Description, "This is my share network")
	th.AssertEquals(t, n.NeutronNetID, "998b42ee-2cee-4d36-8b95-67b5ca1f2109")
	th.AssertEquals(t, n.NeutronSubnetID, "53482b62-2c84-4a53-b6ab-30d9d9800d06")
}

// Verifies that share network deletion works
func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := sharenetworks.Delete(client.ServiceClient(), "fa158a3d-6d9f-4187-9ca5-abbb82646eb2")
	th.AssertNoErr(t, res.Err)
}
