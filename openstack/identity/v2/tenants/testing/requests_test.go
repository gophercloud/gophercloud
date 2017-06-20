package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v2/tenants"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListTenants(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListTenantsSuccessfully(t)

	count := 0
	err := tenants.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := tenants.ExtractTenants(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedTenantSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestCreateTenant(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockCreateTenantResponse(t)

	opts := tenants.CreateOpts{
		Name:        "new_tenant",
		Description: "This is new tenant",
		Enabled:     gophercloud.Enabled,
	}

	tenant, err := tenants.Create(client.ServiceClient(), opts).Extract()

	th.AssertNoErr(t, err)

	expected := &tenants.Tenant{
		Name:        "new_tenant",
		Description: "This is new tenant",
		Enabled:     true,
		ID:          "5c62ef576dc7444cbb73b1fe84b97648",
	}

	th.AssertDeepEquals(t, expected, tenant)
}
