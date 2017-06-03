package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/internal"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestRemainingKeys(t *testing.T) {
	type User struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		City      string
	}

	userStruct := User{
		FirstName: "John",
		LastName:  "Doe",
	}

	userMap := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"city":       "Honolulu",
		"state":      "Hawaii",
	}

	expected := map[string]interface{}{
		"city":  "Honolulu",
		"state": "Hawaii",
	}

	actual := internal.RemainingKeys(userStruct, userMap)
	th.AssertDeepEquals(t, expected, actual)
}
