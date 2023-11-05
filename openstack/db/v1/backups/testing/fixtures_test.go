package testing

import (
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/db/v1/backups"
	"github.com/gophercloud/gophercloud/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/testhelper/fixture"
)

var (
	timestamp  = "2015-11-12T14:22:42"
	timeVal, _ = time.Parse(gophercloud.RFC3339NoZ, timestamp)
)

var backup = `
{
	"created": "` + timestamp + `",
	"datastore": {
		"type": "mysql",
		"version": "5.5",
		"version_id": "b00000b0-00b0-0b00-00b0-000b000000bb"
	},
	"description": "My Backup",
	"id": "{backupId}",
	"instance_id": "44b277eb-39be-4921-be31-3d61b43651d7",
	"locationRef": null,
	"name": "snapshot",
	"parent_id": null,
	"size": null,
	"status": "NEW",
	"updated": "` + timestamp + `"
}
`

var backupGet = `
{
	"created": "` + timestamp + `",
	"datastore": {
		"type": "mysql",
		"version": "5.5",
		"version_id": "b00000b0-00b0-0b00-00b0-000b000000bb"
	},
	"description": "My Incremental Backup",
	"id": "{backupId}",
	"instance_id": "44b277eb-39be-4921-be31-3d61b43651d7",
	"locationRef": "https://openstack.example.com/path/to/backup",
	"name": "Incremental Snapshot",
	"parent_id": "a9832168-7541-4536-b8d9-a8a9b79cf1b4",
	"size": 0.14,
	"status": "COMPLETED",
	"updated": "` + timestamp + `"
}
`

var createReq = `
{
	"backup": {
		"description": "My Backup",
		"incremental": 0,
		"instance": "44b277eb-39be-4921-be31-3d61b43651d7",
		"name": "snapshot"
	}
}
`

var (
	backupId      = "{backupId}"
	_rootURL      = "/backups"
	createResp    = fmt.Sprintf(`{"backup": %s}`, backup)
	listResp      = fmt.Sprintf(`{"backups":[%s]}`, backup)
	getBackupResp = fmt.Sprintf(`{"backup": %s}`, backupGet)
)

var expectedBackup = backups.Backup{
	Created:     timeVal,
	Updated:     timeVal,
	Description: "My Backup",
	ID:          backupId,
	Datastore: datastores.DatastorePartial{
		Type:      "mysql",
		Version:   "5.5",
		VersionID: "b00000b0-00b0-0b00-00b0-000b000000bb",
	},
	InstanceID: "44b277eb-39be-4921-be31-3d61b43651d7",
	Name:       "snapshot",
	Status:     "NEW",
}

var expectedGetBackup = backups.Backup{
	Created:     timeVal,
	Updated:     timeVal,
	Description: "My Incremental Backup",
	ID:          backupId,
	LocationRef: "https://openstack.example.com/path/to/backup",
	Datastore: datastores.DatastorePartial{
		Type:      "mysql",
		Version:   "5.5",
		VersionID: "b00000b0-00b0-0b00-00b0-000b000000bb",
	},
	InstanceID: "44b277eb-39be-4921-be31-3d61b43651d7",
	Name:       "Incremental Snapshot",
	Status:     "COMPLETED",
	ParentId:   "a9832168-7541-4536-b8d9-a8a9b79cf1b4",
	Size:       0.14,
}

func HandleCreate(t *testing.T) {
	fixture.SetupHandler(t, _rootURL, "POST", createReq, createResp, 202)
}

func HandleList(t *testing.T) {
	fixture.SetupHandler(t, _rootURL, "GET", "", listResp, 200)
}

func HandleGet(t *testing.T) {
	fixture.SetupHandler(t, _rootURL+"/{backupId}", "GET", "", getBackupResp, 200)
}

func HandleDelete(t *testing.T) {
	fixture.SetupHandler(t, _rootURL+"/{backupId}", "DELETE", "", "", 202)
}
