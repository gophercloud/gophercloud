package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Verifies that a share type can be created correctly
func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockCreateResponse(t, fakeServer)

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

	st, err := sharetypes.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, st.Name, "my_new_share_type")
	th.AssertEquals(t, st.IsPublic, true)
}

// Verifies that a share type can't be created if the required parameters are missing
func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	options := &sharetypes.CreateOpts{
		Name: "my_new_share_type",
	}

	_, err := sharetypes.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}

	extraSpecs := sharetypes.ExtraSpecsOpts{
		DriverHandlesShareServers: true,
	}

	options = &sharetypes.CreateOpts{
		ExtraSpecs: extraSpecs,
	}

	_, err = sharetypes.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

// Verifies that share type deletion works
func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockDeleteResponse(t, fakeServer)
	res := sharetypes.Delete(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID")
	th.AssertNoErr(t, res.Err)
}

// Verifies that share types can be listed correctly
func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListResponse(t, fakeServer)

	allPages, err := sharetypes.List(client.ServiceClient(fakeServer), &sharetypes.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := sharetypes.ExtractShareTypes(allPages)
	th.AssertNoErr(t, err)
	expected := []sharetypes.ShareType{
		{
			ID:                 "be27425c-f807-4500-a056-d00721db45cf",
			Name:               "default",
			IsPublic:           true,
			ExtraSpecs:         map[string]any{"snapshot_support": "True", "driver_handles_share_servers": "True"},
			RequiredExtraSpecs: map[string]any{"driver_handles_share_servers": "True"},
		},
		{
			ID:                 "f015bebe-c38b-4c49-8832-00143b10253b",
			Name:               "d",
			IsPublic:           true,
			ExtraSpecs:         map[string]any{"driver_handles_share_servers": "false", "snapshot_support": "True"},
			RequiredExtraSpecs: map[string]any{"driver_handles_share_servers": "false"},
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}

// Verifies that it is possible to get the default share type
func TestGetDefault(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetDefaultResponse(t, fakeServer)

	expected := sharetypes.ShareType{
		ID:                 "be27425c-f807-4500-a056-d00721db45cf",
		Name:               "default",
		ExtraSpecs:         map[string]any{"snapshot_support": "True", "driver_handles_share_servers": "True"},
		RequiredExtraSpecs: map[string]any(nil),
	}

	actual, err := sharetypes.GetDefault(context.TODO(), client.ServiceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

// Verifies that it is possible to get the extra specifications for a share type
func TestGetExtraSpecs(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetExtraSpecsResponse(t, fakeServer)

	st, err := sharetypes.GetExtraSpecs(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, st["snapshot_support"], "True")
	th.AssertEquals(t, st["driver_handles_share_servers"], "True")
	th.AssertEquals(t, st["my_custom_extra_spec"], "False")
}

// Verifies that an extra specs can be added to a share type
func TestSetExtraSpecs(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockSetExtraSpecsResponse(t, fakeServer)

	options := &sharetypes.SetExtraSpecsOpts{
		ExtraSpecs: map[string]any{"my_key": "my_value"},
	}

	es, err := sharetypes.SetExtraSpecs(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, es["my_key"], "my_value")
}

// Verifies that an extra specification can be unset for a share type
func TestUnsetExtraSpecs(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockUnsetExtraSpecsResponse(t, fakeServer)
	res := sharetypes.UnsetExtraSpecs(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID", "my_key")
	th.AssertNoErr(t, res.Err)
}

// Verifies that it is possible to see the access for a share type
func TestShowAccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockShowAccessResponse(t, fakeServer)

	expected := []sharetypes.ShareTypeAccess{
		{
			ShareTypeID: "1732f284-401d-41d9-a494-425451e8b4b8",
			ProjectID:   "818a3f48dcd644909b3fa2e45a399a27",
		},
		{
			ShareTypeID: "1732f284-401d-41d9-a494-425451e8b4b8",
			ProjectID:   "e1284adea3ee4d2482af5ed214f3ad90",
		},
	}

	shareType, err := sharetypes.ShowAccess(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID").Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, expected, shareType)
}

// Verifies that an access can be added to a share type
func TestAddAccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockAddAccessResponse(t, fakeServer)

	options := &sharetypes.AccessOpts{
		Project: "e1284adea3ee4d2482af5ed214f3ad90",
	}

	err := sharetypes.AddAccess(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID", options).ExtractErr()
	th.AssertNoErr(t, err)
}

// Verifies that an access can be removed from a share type
func TestRemoveAccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockRemoveAccessResponse(t, fakeServer)

	options := &sharetypes.AccessOpts{
		Project: "e1284adea3ee4d2482af5ed214f3ad90",
	}

	err := sharetypes.RemoveAccess(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID", options).ExtractErr()
	th.AssertNoErr(t, err)
}
