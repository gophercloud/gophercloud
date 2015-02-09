package instances

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func HandleCreateInstanceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		th.TestJSONRequest(t, r, `
{
  "instance": {
    "databases": [
      {
        "character_set": "utf8",
        "collate": "utf8_general_ci",
        "name": "sampledb"
      },
      {
        "name": "nextround"
      }
    ],
    "flavorRef": "1",
    "name": "json_rack_instance",
    "users": [
      {
        "databases": [
          {
            "name": "sampledb"
          }
        ],
        "name": "demouser",
        "password": "demopassword"
      }
    ],
    "volume": {
      "size": 2
    },
		"restorePoint": "1234567890"
  }
}
`)

		fmt.Fprintf(w, `
{
  "instance": {
    "flavorRef": 1,
    "name": "json_restore",
    "restorePoint": {
      "backupRef": "61f12fef-edb1-4561-8122-e7c00ef26a82"
    },
    "volume": {
      "size": 2
    }
  }
}
`)
	})
}
