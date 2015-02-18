package users

const singleDB = `{"databases": [{"name": "databaseE"}]}`

var changePwdReq = `
{
  "users": [
    {
      "name": "dbuser1",
      "password": "newpassword"
    },
    {
      "name": "dbuser2",
      "password": "anotherpassword"
    }
  ]
}
`

var updateReq = `
{
	"user": {
		"name": "new_username",
		"password": "new_password"
	}
}
`

var getResp = `
{
	"user": {
		"name": "exampleuser",
		"host": "foo",
		"databases": [
			{
				"name": "databaseA"
			},
			{
				"name": "databaseB"
			}
		]
	}
}
`

var (
	listUserAccessResp = singleDB
	grantUserAccessReq = singleDB
)
