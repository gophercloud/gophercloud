package testing

import (
	"context"
	"fmt"
	"net/http"
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

	th.AssertEquals(t, "my_new_share_type", st.Name)
	th.AssertTrue(t, st.IsPublic)
}

func TestCreateExtraSpecs(t *testing.T) {
	falseValue, trueValue := false, true
	for _, driverHandlesShareServers := range []bool{false, true} {
		for _, snapshot := range []struct {
			name  string
			value *bool
		}{
			{name: "omitted"},
			{name: "false", value: &falseValue},
			{name: "true", value: &trueValue},
		} {
			t.Run(fmt.Sprintf("driver_handles_share_servers=%t/snapshot_support=%s", driverHandlesShareServers, snapshot.name), func(t *testing.T) {
				fakeServer := th.SetupHTTP()
				defer fakeServer.Teardown()

				snapshotJSON := ""
				if snapshot.value != nil {
					snapshotJSON = fmt.Sprintf(`,"snapshot_support":%t`, *snapshot.value)
				}
				fakeServer.Mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
					th.TestMethod(t, r, "POST")
					th.TestJSONRequest(t, r, fmt.Sprintf(`{
						"share_type": {
							"name": "my_new_share_type",
							"os-share-type-access:is_public": false,
							"extra_specs": {"driver_handles_share_servers": %t%s}
						}
					}`, driverHandlesShareServers, snapshotJSON))
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusAccepted)
					fmt.Fprint(w, `{"share_type":{"id":"share-type-id"}}`)
				})

				opts := sharetypes.CreateOpts{
					Name: "my_new_share_type",
					ExtraSpecs: sharetypes.ExtraSpecsOpts{
						DriverHandlesShareServers: driverHandlesShareServers,
						SnapshotSupport:           snapshot.value,
					},
				}
				shareType, err := sharetypes.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
				th.AssertNoErr(t, err)
				th.AssertEquals(t, "share-type-id", shareType.ID)
			})
		}
	}
}

// Verifies that a share type can't be created without a name.
func TestRequiredCreateOpts(t *testing.T) {
	for _, driverHandlesShareServers := range []bool{false, true} {
		t.Run(fmt.Sprintf("driver_handles_share_servers=%t", driverHandlesShareServers), func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()

			options := &sharetypes.CreateOpts{
				ExtraSpecs: sharetypes.ExtraSpecsOpts{
					DriverHandlesShareServers: driverHandlesShareServers,
				},
			}
			_, err := sharetypes.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
			th.AssertDeepEquals(t, gophercloud.ErrMissingInput{Argument: "Name"}, err)
		})
	}
}

// Verifies that a share type can be updated
func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockUpdateResponse(t, fakeServer)

	name := "my_new_share_type_name"
	description := "my new share type description"
	isPublic := false

	options := &sharetypes.UpdateOpts{
		Name:        &name,
		Description: &description,
		IsPublic:    &isPublic,
	}

	st, err := sharetypes.Update(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "my_new_share_type_name", st.Name)
	th.AssertEquals(t, "my new share type description", *st.Description)
	th.AssertEquals(t, false, st.IsPublic)
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

// Verifies that a share type can be retrieved by ID
func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetResponse(t, fakeServer)

	description := "my share type description"
	isDefault := false
	expected := sharetypes.ShareType{
		ID:                 "1d600d02-26a7-4b23-af3d-7d51860fe858",
		Name:               "my_share_type",
		IsPublic:           true,
		ExtraSpecs:         map[string]any{"snapshot_support": "True", "driver_handles_share_servers": "True"},
		RequiredExtraSpecs: map[string]any{"driver_handles_share_servers": "True"},
		Description:        &description,
		IsDefault:          &isDefault,
	}

	actual, err := sharetypes.Get(context.TODO(), client.ServiceClient(fakeServer), "shareTypeID").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
}

// Verifies that a share type can be retrieved by ID at microversion v6 or greater.
func TestGetV6(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetResponseV6(t, fakeServer)

	description := "my share type description"
	isDefault := false
	expected := sharetypes.ShareType{
		ID:                 "1d600d02-26a7-4b23-af3d-7d51860fe858",
		Name:               "my_share_type",
		IsPublic:           true,
		ExtraSpecs:         map[string]any{"snapshot_support": "True", "driver_handles_share_servers": "True"},
		RequiredExtraSpecs: map[string]any{"driver_handles_share_servers": "True"},
		Description:        &description,
		IsDefault:          &isDefault,
	}

	testClient := client.ServiceClient(fakeServer)
	actual, err := sharetypes.Get(context.TODO(), testClient, "shareTypeID").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, actual)
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

	th.AssertEquals(t, "True", st["snapshot_support"])
	th.AssertEquals(t, "True", st["driver_handles_share_servers"])
	th.AssertEquals(t, "False", st["my_custom_extra_spec"])
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

	th.AssertEquals(t, "my_value", es["my_key"])
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
