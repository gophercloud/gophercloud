package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharetypes"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// Verifies that a share type can be created correctly
func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	snapshotSupport := true
	extraSpecs := sharetypes.ExtraSpecsOpts{
		DriverHandlesShareServers: true,
		SnapshotSupport:           &snapshotSupport,
	}

	options := &sharetypes.CreateOpts{
		Name:       "my_new_share_type",
		IsPublic:   true,
		ExtraSpecs: extraSpecs,
	}

	st, err := sharetypes.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, st.Name, "my_new_share_type")
	th.AssertEquals(t, st.IsPublic, true)
}

// Verifies that a share type can't be created if the required parameters are missing
func TestCreateFails(t *testing.T) {
	options := &sharetypes.CreateOpts{
		Name: "my_new_share_type",
	}

	_, err := sharetypes.Create(client.ServiceClient(), options).Extract()
	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}

	extraSpecs := sharetypes.ExtraSpecsOpts{
		DriverHandlesShareServers: true,
	}

	options = &sharetypes.CreateOpts{
		ExtraSpecs: extraSpecs,
	}

	_, err = sharetypes.Create(client.ServiceClient(), options).Extract()
	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

// Verifies that share type deletion works
func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)
	res := sharetypes.Delete(client.ServiceClient(), "shareTypeID")
	th.AssertNoErr(t, res.Err)
}
