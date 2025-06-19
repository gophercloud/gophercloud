package testing

import (
	"fmt"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

const user1 = `
{"databases": [{"name": "databaseA"}],"name": "dbuser3"%s}
`

const user2 = `
{"databases": [{"name": "databaseB"},{"name": "databaseC"}],"name": "dbuser4"%s}
`

var (
	instanceID = "{instanceID}"
	_rootURL   = "/instances/" + instanceID + "/users"
	pUser1     = fmt.Sprintf(user1, `,"password":"secretsecret"`)
	pUser2     = fmt.Sprintf(user2, `,"password":"secretsecret"`)
	createReq  = fmt.Sprintf(`{"users":[%s, %s]}`, pUser1, pUser2)
	listResp   = fmt.Sprintf(`{"users":[%s, %s]}`, fmt.Sprintf(user1, ""), fmt.Sprintf(user2, ""))
)

func HandleCreate(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, _rootURL, "POST", createReq, "", 202)
}

func HandleList(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, _rootURL, "GET", "", listResp, 200)
}

func HandleDelete(t *testing.T, fakeServer th.FakeServer) {
	fixture.SetupHandler(t, fakeServer, _rootURL+"/{userName}", "DELETE", "", "", 202)
}
