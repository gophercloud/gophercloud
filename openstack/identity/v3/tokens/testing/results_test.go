package testing

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/testhelper"
)

var (
	testID       = "testID"
	testExpireAt = "2033-03-14T15:09:26.535897Z"

	testEndpointIDs = []string{
		"endpoint1ID",
		"endpoint2ID",
		"endpoint3ID",
	}
	testEndpointInterfaces = [...]string{
		"public",
		"internal",
		"admin",
	}
	testEndpointRegions = [...]string{
		"region1",
		"region2",
		"region3",
	}
	testEndpointURLs = [...]string{
		"http://endpoint1.pl",
		"http://endpoint2.io",
		"http://endpoint3.onion",
	}
	testCatalogID   = "catalogID"
	testCatalogType = "identity"
	testCatalogName = "keystone"

	testDomainID   = "domainID"
	testDomainName = "domainName"
	testUserID     = "userID"
	testUserName   = "userName"

	testRoleIDs = [...]string{
		"role1ID",
		"role2ID",
		"role3ID",
	}
	testRoleNames = [...]string{
		"role1Name",
		"role2Name",
		"role3Name",
	}

	testProjectID   = "projectID"
	testProjectName = "projectName"
)

func TestExtractToken(t *testing.T) {
	result := getGetResultFromResponse(t, `{
		"token": {
			"expires_at": "`+testExpireAt+`"
		}
	}`)

	token, err := result.ExtractToken()
	testhelper.AssertNoErr(t, err)

	testhelper.CheckDeepEquals(t, testID, token.ID)
	testhelper.CheckDeepEquals(t, testExpireAt,
		token.ExpiresAt.Format(gophercloud.RFC3339Milli))
}

func TestExtractCatalog(t *testing.T) {
	endpoints := ""
	for i := 0; i < len(testEndpointIDs); i++ {
		endpoints += getEndpoint(testEndpointIDs[i], testEndpointInterfaces[i],
			testEndpointRegions[i], testEndpointURLs[i]) + ","
	}
	result := getGetResultFromResponse(t, `{
		"token": {
			"catalog": [{
				"id": "`+testCatalogID+`",
				"name": "`+testCatalogName+`",
				"type": "`+testCatalogType+`",
				"endpoints": [`+endpoints[:len(endpoints)-1]+`]
			}]
		}
	}`)

	catalog, err := result.ExtractServiceCatalog()
	testhelper.AssertNoErr(t, err)

	testhelper.CheckDeepEquals(t, 1, len(catalog.Entries))
	testhelper.CheckDeepEquals(t, testCatalogID, catalog.Entries[0].ID)
	testhelper.CheckDeepEquals(t, testCatalogName, catalog.Entries[0].Name)
	testhelper.CheckDeepEquals(t, testCatalogType, catalog.Entries[0].Type)

	testhelper.CheckDeepEquals(t, 3, len(catalog.Entries[0].Endpoints))
	for i, endpoint := range catalog.Entries[0].Endpoints {
		testhelper.CheckDeepEquals(t, testEndpointIDs[i], endpoint.ID)
		testhelper.CheckDeepEquals(t, testEndpointInterfaces[i], endpoint.Interface)
		testhelper.CheckDeepEquals(t, testEndpointRegions[i], endpoint.Region)
		testhelper.CheckDeepEquals(t, testEndpointURLs[i], endpoint.URL)
	}
}

func TestExtractUser(t *testing.T) {
	result := getGetResultFromResponse(t, `{
		"token": {
			"user": {
				"domain": {
					"id": "`+testDomainID+`",
					"name": "`+testDomainName+`"
				},
				"id": "`+testUserID+`",
				"name": "`+testUserName+`"
			}
		}
	}`)

	user, err := result.ExtractUser()
	testhelper.AssertNoErr(t, err)

	testhelper.CheckDeepEquals(t, testDomainID, user.Domain.ID)
	testhelper.CheckDeepEquals(t, testDomainName, user.Domain.Name)
	testhelper.CheckDeepEquals(t, testUserID, user.ID)
	testhelper.CheckDeepEquals(t, testUserName, user.Name)
}

func TestExtractRoles(t *testing.T) {
	rawRoles := ""
	for i := 0; i < len(testRoleIDs); i++ {
		rawRoles += getRole(testRoleIDs[i], testRoleNames[i]) + ","
	}
	result := getGetResultFromResponse(t, `{
		"token": {
			"roles": [`+rawRoles[:len(rawRoles)-1]+`]
		}
	}`)

	roles, err := result.ExtractRoles()
	testhelper.AssertNoErr(t, err)

	testhelper.CheckDeepEquals(t, 3, len(roles.Roles))
	for i, role := range roles.Roles {
		testhelper.CheckDeepEquals(t, testRoleIDs[i], role.ID)
		testhelper.CheckDeepEquals(t, testRoleNames[i], role.Name)
	}
}

func TestExtractProject(t *testing.T) {
	result := getGetResultFromResponse(t, `{
		"token": {
			"project": {
				"id": "`+testProjectID+`",
				"name": "`+testProjectName+`"
			}
		}
	}`)

	project, err := result.ExtractProject()
	testhelper.AssertNoErr(t, err)

	testhelper.CheckDeepEquals(t, testProjectID, project.ID)
	testhelper.CheckDeepEquals(t, testProjectName, project.Name)
}

func getGetResultFromResponse(t *testing.T, response string) tokens.GetResult {
	result := tokens.GetResult{}
	result.Header = http.Header{
		"X-Subject-Token": []string{testID},
	}
	err := json.Unmarshal([]byte(response), &result.Body)
	testhelper.AssertNoErr(t, err)
	return result
}

func getEndpoint(id, interfaceName, region, url string) string {
	return `{
		"id": "` + id + `",
		"interface": "` + interfaceName + `",
		"region": "` + region + `",
		"url": "` + url + `"
	}`
}

func getRole(id, name string) string {
	return `{
		"id": "` + id + `",
		"name": "` + name + `"
	}`
}
