package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListRoles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRolesSuccessfully(t, fakeServer)

	count := 0
	err := roles.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := roles.ExtractRoles(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedRolesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListRolesAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRolesSuccessfully(t, fakeServer)

	allPages, err := roles.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := roles.ExtractRoles(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedRolesSlice, actual)
	th.AssertEquals(t, ExpectedRolesSlice[1].Extra["description"], "read-only support role")
}

func TestListUsersFiltersCheck(t *testing.T) {
	type test struct {
		filterName string
		wantErr    bool
	}
	tests := []test{
		{"foo__contains", false},
		{"foo", true},
		{"foo_contains", true},
		{"foo__", true},
		{"__foo", true},
	}

	var listOpts roles.ListOpts
	for _, _test := range tests {
		listOpts.Filters = map[string]string{_test.filterName: "bar"}
		_, err := listOpts.ToRoleListQuery()

		if !_test.wantErr {
			th.AssertNoErr(t, err)
		} else {
			switch _t := err.(type) {
			case nil:
				t.Fatal("error expected but got a nil")
			case roles.InvalidListFilter:
			default:
				t.Fatalf("unexpected error type: [%T]", _t)
			}
		}
	}
}

func TestGetRole(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetRoleSuccessfully(t, fakeServer)

	actual, err := roles.Get(context.TODO(), client.ServiceClient(fakeServer), "9fe1d3").Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondRole, *actual)
}

func TestCreateRole(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateRoleSuccessfully(t, fakeServer)

	createOpts := roles.CreateOpts{
		Name:     "support",
		DomainID: "1789d1",
		Extra: map[string]any{
			"description": "read-only support role",
		},
	}

	actual, err := roles.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondRole, *actual)
}

func TestUpdateRole(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateRoleSuccessfully(t, fakeServer)

	updateOpts := roles.UpdateOpts{
		Extra: map[string]any{
			"description": "admin read-only support role",
		},
	}

	actual, err := roles.Update(context.TODO(), client.ServiceClient(fakeServer), "9fe1d3", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondRoleUpdated, *actual)
}

func TestDeleteRole(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteRoleSuccessfully(t, fakeServer)

	res := roles.Delete(context.TODO(), client.ServiceClient(fakeServer), "9fe1d3")
	th.AssertNoErr(t, res.Err)
}

func TestListAssignmentsSinglePage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRoleAssignmentsSuccessfully(t, fakeServer)

	count := 0
	err := roles.ListAssignments(client.ServiceClient(fakeServer), roles.ListAssignmentsOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := roles.ExtractRoleAssignments(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedRoleAssignmentsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAssignmentsWithNamesSinglePage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRoleAssignmentsWithNamesSuccessfully(t, fakeServer)

	var includeNames = true
	listOpts := roles.ListAssignmentsOpts{
		IncludeNames: &includeNames,
	}

	count := 0
	err := roles.ListAssignments(client.ServiceClient(fakeServer), listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := roles.ExtractRoleAssignments(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedRoleAssignmentsWithNamesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAssignmentsWithSubtreeSinglePage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRoleAssignmentsWithSubtreeSuccessfully(t, fakeServer)

	var includeSubtree = true
	listOpts := roles.ListAssignmentsOpts{
		IncludeSubtree: &includeSubtree,
	}

	count := 0
	err := roles.ListAssignments(client.ServiceClient(fakeServer), listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := roles.ExtractRoleAssignments(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedRoleAssignmentsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAssignmentsOnResource_ProjectsUsers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAssignmentsOnResourceSuccessfully_ProjectsUsers(t, fakeServer)

	count := 0
	err := roles.ListAssignmentsOnResource(client.ServiceClient(fakeServer), roles.ListAssignmentsOnResourceOpts{
		UserID:    "{user_id}",
		ProjectID: "{project_id}",
	}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := roles.ExtractRoles(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedRolesOnResourceSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAssignmentsOnResource_DomainsUsers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAssignmentsOnResourceSuccessfully_DomainsUsers(t, fakeServer)

	count := 0
	err := roles.ListAssignmentsOnResource(client.ServiceClient(fakeServer), roles.ListAssignmentsOnResourceOpts{
		UserID:   "{user_id}",
		DomainID: "{domain_id}",
	}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := roles.ExtractRoles(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedRolesOnResourceSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAssignmentsOnResource_ProjectsGroups(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAssignmentsOnResourceSuccessfully_ProjectsGroups(t, fakeServer)

	count := 0
	err := roles.ListAssignmentsOnResource(client.ServiceClient(fakeServer), roles.ListAssignmentsOnResourceOpts{
		GroupID:   "{group_id}",
		ProjectID: "{project_id}",
	}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := roles.ExtractRoles(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedRolesOnResourceSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAssignmentsOnResource_DomainsGroups(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAssignmentsOnResourceSuccessfully_DomainsGroups(t, fakeServer)

	count := 0
	err := roles.ListAssignmentsOnResource(client.ServiceClient(fakeServer), roles.ListAssignmentsOnResourceOpts{
		GroupID:  "{group_id}",
		DomainID: "{domain_id}",
	}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := roles.ExtractRoles(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedRolesOnResourceSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestAssign(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAssignSuccessfully(t, fakeServer)

	err := roles.Assign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.AssignOpts{
		UserID:    "{user_id}",
		ProjectID: "{project_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = roles.Assign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.AssignOpts{
		UserID:   "{user_id}",
		DomainID: "{domain_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = roles.Assign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.AssignOpts{
		GroupID:   "{group_id}",
		ProjectID: "{project_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = roles.Assign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.AssignOpts{
		GroupID:  "{group_id}",
		DomainID: "{domain_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUnassign(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUnassignSuccessfully(t, fakeServer)

	err := roles.Unassign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.UnassignOpts{
		UserID:    "{user_id}",
		ProjectID: "{project_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = roles.Unassign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.UnassignOpts{
		UserID:   "{user_id}",
		DomainID: "{domain_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = roles.Unassign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.UnassignOpts{
		GroupID:   "{group_id}",
		ProjectID: "{project_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	err = roles.Unassign(context.TODO(), client.ServiceClient(fakeServer), "{role_id}", roles.UnassignOpts{
		GroupID:  "{group_id}",
		DomainID: "{domain_id}",
	}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestCreateRoleInferenceRule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateRoleInferenceRule(t, fakeServer)

	actual, err := roles.CreateRoleInferenceRule(context.TODO(), client.ServiceClient(fakeServer), "7ceab6192ea34a548cc71b24f72e762c", "97e2f5d38bc94842bc3da818c16762ed").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedRoleInferenceRule, *actual)
}

func TestListRoleInferenceRules(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListRoleInferenceRules(t, fakeServer)

	actual, err := roles.ListRoleInferenceRules(context.TODO(), client.ServiceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedRoleInferenceRuleList, *actual)
}

func TestDeleteRoleInferenceRule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteRoleInferenceRule(t, fakeServer)

	err := roles.DeleteRoleInferenceRule(context.TODO(), client.ServiceClient(fakeServer), "7ceab6192ea34a548cc71b24f72e762c", "97e2f5d38bc94842bc3da818c16762ed").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetInferenceRule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetRoleInferenceRule(t, fakeServer)

	actual, err := roles.GetRoleInferenceRule(context.TODO(), client.ServiceClient(fakeServer), "7ceab6192ea34a548cc71b24f72e762c", "97e2f5d38bc94842bc3da818c16762ed").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedRoleInferenceRule, *actual)
}
