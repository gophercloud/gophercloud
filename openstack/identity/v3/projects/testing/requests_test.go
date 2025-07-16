package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAvailableProjects(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAvailableProjectsSuccessfully(t, fakeServer)

	count := 0
	err := projects.ListAvailable(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := projects.ExtractProjects(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedAvailableProjectsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListProjects(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListProjectsSuccessfully(t, fakeServer)

	count := 0
	err := projects.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := projects.ExtractProjects(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedProjectSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListGroupsFiltersCheck(t *testing.T) {
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

	var listOpts projects.ListOpts
	for _, _test := range tests {
		listOpts.Filters = map[string]string{_test.filterName: "bar"}
		_, err := listOpts.ToProjectListQuery()

		if !_test.wantErr {
			th.AssertNoErr(t, err)
		} else {
			switch _t := err.(type) {
			case nil:
				t.Fatal("error expected but got a nil")
			case projects.InvalidListFilter:
			default:
				t.Fatalf("unexpected error type: [%T]", _t)
			}
		}
	}
}

func TestGetProject(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetProjectSuccessfully(t, fakeServer)

	actual, err := projects.Get(context.TODO(), client.ServiceClient(fakeServer), "1234").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, RedTeam, *actual)
}

func TestCreateProject(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateProjectSuccessfully(t, fakeServer)

	createOpts := projects.CreateOpts{
		Name:        "Red Team",
		Description: "The team that is red",
		Tags:        []string{"Red", "Team"},
		Extra:       map[string]any{"test": "old"},
	}

	actual, err := projects.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, RedTeam, *actual)
}

func TestDeleteProject(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteProjectSuccessfully(t, fakeServer)

	res := projects.Delete(context.TODO(), client.ServiceClient(fakeServer), "1234")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateProject(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateProjectSuccessfully(t, fakeServer)

	var description = "The team that is bright red"
	updateOpts := projects.UpdateOpts{
		Name:        "Bright Red Team",
		Description: &description,
		Tags:        &[]string{"Red"},
		Extra:       map[string]any{"test": "new"},
	}

	actual, err := projects.Update(context.TODO(), client.ServiceClient(fakeServer), "1234", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, UpdatedRedTeam, *actual)
	t.Log(projects.Update(context.TODO(), client.ServiceClient(fakeServer), "1234", updateOpts))
}

func TestListProjectTags(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListProjectTagsSuccessfully(t, fakeServer)

	actual, err := projects.ListTags(context.TODO(), client.ServiceClient(fakeServer), "966b3c7d36a24facaf20b7e458bf2192").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedTags, *actual)
}

func TestModifyProjectTags(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleModifyProjectTagsSuccessfully(t, fakeServer)

	modifyOpts := projects.ModifyTagsOpts{
		Tags: []string{"foo", "bar"},
	}
	actual, err := projects.ModifyTags(context.TODO(), client.ServiceClient(fakeServer), "966b3c7d36a24facaf20b7e458bf2192", modifyOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedProjects, *actual)
}

func TestDeleteTags(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteProjectTagsSuccessfully(t, fakeServer)

	err := projects.DeleteTags(context.TODO(), client.ServiceClient(fakeServer), "966b3c7d36a24facaf20b7e458bf2192").ExtractErr()
	th.AssertNoErr(t, err)
}
