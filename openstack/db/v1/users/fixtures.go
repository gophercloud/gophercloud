package users

import "fmt"

const user1 = `
{"databases": [{"name": "databaseA"}],"name": "dbuser3"%s}
`

const user2 = `
{"databases": [{"name": "databaseB"},{"name": "databaseC"}],"name": "dbuser4"%s}
`

var (
	pUser1 = fmt.Sprintf(user1, `,"password":"secretsecret"`)
	pUser2 = fmt.Sprintf(user2, `,"password":"secretsecret"`)
)

var (
	createReq = fmt.Sprintf(`{"users":[%s, %s]}`, pUser1, pUser2)
	listResp  = fmt.Sprintf(`{"users":[%s, %s]}`, fmt.Sprintf(user1, ""), fmt.Sprintf(user2, ""))
)
