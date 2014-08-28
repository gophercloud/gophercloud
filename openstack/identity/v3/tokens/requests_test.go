package tokens

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

func TestCreateUserIDAndPassword(t *testing.T) {
	setup()
	defer teardown()

	client := gophercloud.ServiceClient{
		Endpoint: endpoint(),
		Options:  gophercloud.AuthOptions{UserID: "me", Password: "squirrel!"},
	}

	mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")

		testhelper.TestJSONRequest(t, r, `
			{
				"auth": {
					"identity": {
						"methods": ["password"],
						"password": {
							"user": { "id": "me", "password": "squirrel!" }
						}
					}
				}
			}
		`)

		fmt.Fprintf(w, `{}`)
	})

	_, err := Create(&client, nil)
	if err != nil {
		t.Errorf("Create returned an error: %v", err)
	}
}
