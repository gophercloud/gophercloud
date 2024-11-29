package testing

import (
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestExtractToken(t *testing.T) {
	result := getGetResult(t)

	token, err := result.ExtractToken()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &ExpectedToken, token)
}

func TestExtractCatalog(t *testing.T) {
	result := getGetResult(t)

	catalog, err := result.ExtractServiceCatalog()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &ExpectedServiceCatalog, catalog)
}

func TestExtractUser(t *testing.T) {
	result := getGetResult(t)

	user, err := result.ExtractUser()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &ExpectedUser, user)
}

func TestExtractRoles(t *testing.T) {
	result := getGetResult(t)

	roles, err := result.ExtractRoles()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ExpectedRoles, roles)
}

func TestExtractProject(t *testing.T) {
	result := getGetResult(t)

	project, err := result.ExtractProject()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &ExpectedProject, project)
}

func TestExtractDomain(t *testing.T) {
	result := getGetDomainResult(t)

	domain, err := result.ExtractDomain()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &ExpectedDomain, domain)
}
