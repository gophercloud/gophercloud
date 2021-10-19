package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/impliedroles"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListImpliedRoles(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListImpliedRolesSuccessfully(t)

	count := 0

	err := impliedroles.List(client.ServiceClient(), "b385b97c988f4a649eecbb5cdd52b7e1", nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := impliedroles.ExtractImpliedRoles(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedImpliedRoleSlice, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestCreateImpliedRole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateImpliedRoleSuccessfully(t)

	priorRoleId := "b385b97c988f4a649eecbb5cdd52b7e1"
	impliesRoleId := "ddb5331895e348b0ab78cf0db18e8b78"

	_, err := impliedroles.Create(client.ServiceClient(), priorRoleId, impliesRoleId).Extract()
	th.AssertNoErr(t, err)
}

func TestDeleteImpliedRole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteImpliedRoleSuccessfully(t)

	priorRoleId := "b385b97c988f4a649eecbb5cdd52b7e1"
	impliesRoleId := "ddb5331895e348b0ab78cf0db18e8b78"

	err := impliedroles.Delete(client.ServiceClient(), priorRoleId, impliesRoleId).ExtractErr()

	th.AssertNoErr(t, err)
}
