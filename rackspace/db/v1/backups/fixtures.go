package backups

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

var singleBackup = `
{
  "backup": {
    "created": "2014-02-13T21:47:16",
    "description": "My Backup",
    "id": "61f12fef-edb1-4561-8122-e7c00ef26a82",
    "instance_id": "d4603f69-ec7e-4e9b-803f-600b9205576f",
    "locationRef": null,
    "name": "snapshot",
    "parent_id": null,
    "size": 100,
    "status": "NEW",
		"datastore": {
			"version": "5.1",
			"type": "MySQL",
			"version_id": "20000000-0000-0000-0000-000000000002"
		},
    "updated": "2014-02-13T21:47:16"
  }
}
`

func SetupHandler(t *testing.T, url, method, requestBody, responseBody string, status int) {
	th.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, method)
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		if requestBody != "" {
			th.TestJSONRequest(t, r, requestBody)
		}

		if responseBody != "" {
			w.Header().Add("Content-Type", "application/json")
		}

		w.WriteHeader(status)

		if responseBody != "" {
			fmt.Fprintf(w, responseBody)
		}
	})
}

func HandleCreateSuccessfully(t *testing.T) {
	requestJSON := `
{
  "backup": {
    "description": "My Backup",
    "instance": "d4603f69-ec7e-4e9b-803f-600b9205576f",
    "name": "snapshot"
  }
}
`

	SetupHandler(t, "/backups", "POST", requestJSON, singleBackup, 202)
}

func HandleListSuccessfully(t *testing.T) {
	responseJSON := `
{
  "backups": [
    {
      "status": "COMPLETED",
      "updated": "2014-06-18T21:24:39",
      "description": "Backup from Restored Instance",
      "datastore": {
        "version": "5.1",
        "type": "MySQL",
        "version_id": "20000000-0000-0000-0000-000000000002"
      },
      "id": "87972694-4be2-40f5-83f8-501656e0032a",
      "size": 0.141026,
      "name": "restored_backup",
      "created": "2014-06-18T21:23:35",
      "instance_id": "29af2cd9-0674-48ab-b87a-b160f00208e6",
      "parent_id": null,
      "locationRef": "http://localhost/path/to/backup"
    }
  ]
}
`

	SetupHandler(t, "/backups", "GET", "", responseJSON, 200)
}

func HandleGetSuccessfully(t *testing.T, backupID string) {
	SetupHandler(t, "/backups/"+backupID, "GET", "", singleBackup, 200)
}

func HandleDeleteSuccessfully(t *testing.T, backupID string) {
	SetupHandler(t, "/backups/"+backupID, "DELETE", "", "", 202)
}
