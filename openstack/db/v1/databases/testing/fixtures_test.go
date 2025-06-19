package testing

import (
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

var (
	instanceID = "{instanceID}"
	resURL     = "/instances/" + instanceID + "/databases"
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

func HandleCreate(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL, "POST", createDBsReq, "", 202)
}

func HandleList(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL, "GET", "", listDBsResp, 200)
}

func HandleDelete(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, resURL+"/{dbName}", "DELETE", "", "", 202)
}
