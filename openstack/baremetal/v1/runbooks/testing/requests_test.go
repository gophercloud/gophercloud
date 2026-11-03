package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/runbooks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListRunbooksFollowsNextLink(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRunbookListSuccessfully(t, fakeServer)

	var names []string
	err := runbooks.List(client.ServiceClient(fakeServer), runbooks.ListOpts{Limit: 1}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		actual, err := runbooks.ExtractRunbooks(page)
		if err != nil {
			return false, err
		}
		for _, r := range actual {
			names = append(names, r.Name)
		}
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, []string{"CUSTOM_AWESOME", "CUSTOM_OTHER"}, names)
}

func TestListRunbooksWithAllOptions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRunbookListDetailSuccessfully(t, fakeServer)

	allPages, err := runbooks.List(client.ServiceClient(fakeServer), runbooks.ListOpts{
		Fields:  []string{"uuid", "name"},
		Limit:   10,
		Marker:  "fc6b1709-8dd5-86b0-2d34-5203d1c29127",
		SortDir: "desc",
		SortKey: "name",
		Detail:  true,
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := runbooks.ExtractRunbooks(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, []runbooks.Runbook{RunbookAwesome}, actual)
}

func TestCreateRunbook(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRunbookCreationSuccessfully(t, fakeServer)

	disableRamdisk := false
	public := false
	actual, err := runbooks.Create(context.TODO(), client.ServiceClient(fakeServer), runbooks.CreateOpts{
		Name:           "CUSTOM_AWESOME",
		UUID:           "fc6b1709-8dd5-86b0-2d34-5203d1c29127",
		Description:    "Enable logical processors",
		DisableRamdisk: &disableRamdisk,
		Extra:          map[string]any{"team": "metal"},
		Owner:          "f5a7c3a8b9e94d7f8a1b2c3d4e5f6a7b",
		Public:         &public,
		Steps: []runbooks.Step{
			{
				Interface: nodes.InterfaceBIOS,
				Step:      "apply_configuration",
				Args: map[string]any{
					"settings": []map[string]string{{"name": "LogicalProc", "value": "Enabled"}},
				},
				// Order 0 is valid and must still be sent.
				Order: 0,
			},
		},
	}).Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, RunbookAwesome, *actual)
}

func TestCreateRunbookRequiresSteps(t *testing.T) {
	_, err := runbooks.CreateOpts{Name: "CUSTOM_AWESOME"}.ToRunbookCreateMap()
	if err == nil {
		t.Fatal("expected an error for a runbook without steps")
	}
}

func TestGetRunbook(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRunbookGetSuccessfully(t, fakeServer)

	actual, err := runbooks.Get(context.TODO(), client.ServiceClient(fakeServer), runbookID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, RunbookAwesome, *actual)
}

func TestUpdateRunbook(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRunbookUpdateSuccessfully(t, fakeServer)

	actual, err := runbooks.Update(context.TODO(), client.ServiceClient(fakeServer), runbookID, runbooks.UpdateOpts{
		runbooks.UpdateOperation{
			Op:    runbooks.ReplaceOp,
			Path:  "/description",
			Value: "Enable logical processors",
		},
		runbooks.UpdateOperation{
			Op:   runbooks.RemoveOp,
			Path: "/extra/team",
		},
	}).Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, RunbookAwesome, *actual)
}

func TestDeleteRunbook(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRunbookDeletionSuccessfully(t, fakeServer)

	err := runbooks.Delete(context.TODO(), client.ServiceClient(fakeServer), runbookID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestRunbookTraits(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRunbookTraitsSuccessfully(t, fakeServer)
	HandleRunbookTraitSuccessfully(t, fakeServer)
	c := client.ServiceClient(fakeServer)

	traits, err := runbooks.ListTraits(context.TODO(), c, runbookID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, []string{"CUSTOM_GPU", "CUSTOM_FPGA"}, traits)

	th.AssertNoErr(t, runbooks.SetTraits(context.TODO(), c, runbookID, []string{"CUSTOM_GPU", "CUSTOM_FPGA"}).ExtractErr())
	th.AssertNoErr(t, runbooks.AddTrait(context.TODO(), c, runbookID, "CUSTOM_GPU").ExtractErr())
	th.AssertNoErr(t, runbooks.DeleteTrait(context.TODO(), c, runbookID, "CUSTOM_GPU").ExtractErr())
	th.AssertNoErr(t, runbooks.DeleteTraits(context.TODO(), c, runbookID).ExtractErr())
}
