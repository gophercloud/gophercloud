package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/testhelper/fixture"
)

var (
	instanceID = "{instanceID}"
	userName   = "{userName}"
	resURL     = "/instances/" + instanceID + "/databases"
	grantURL   = "/instances/" + instanceID + "/users/" + userName + "/databases"
)

var createDBsReq = `
{
	"databases": [
		{
			"character_set": "utf8",
			"collate": "utf8_general_ci",
			"name": "testingdb"
		},
		{
			"name": "sampledb"
		}
	]
}
`

var listDBsResp = `
{
	"databases": [
		{
			"name": "anotherexampledb"
		},
		{
			"name": "exampledb"
		},
		{
			"name": "nextround"
		},
		{
			"name": "sampledb"
		},
		{
			"name": "testingdb"
		}
	]
}
`

var GrantAccessReq = `
{
	"databases": [
		{
			"name": "anotherexampledb"
		},
		{
			"name": "exampledb"
		}
	]
}
`

func HandleCreate(t *testing.T) {
	fixture.SetupHandler(t, resURL, "POST", createDBsReq, "", 202)
}

func HandleList(t *testing.T) {
	fixture.SetupHandler(t, resURL, "GET", "", listDBsResp, 200)
}

func HandleDelete(t *testing.T) {
	fixture.SetupHandler(t, resURL+"/{dbName}", "DELETE", "", "", 202)
}
