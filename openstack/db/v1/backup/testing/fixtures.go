package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/testhelper/fixture"
	"github.com/gophercloud/gophercloud/openstack/db/v1/backup"
	"github.com/gophercloud/gophercloud/openstack/db/v1/datastores"
	"fmt"
)

var (
	instanceID        = "{instanceID}"
	backupID          = "{backupID}"
	resURL            = "/backups/" + backupID
	baseURL           = "/backups"
	instanceBackupURL = "/instances/" + instanceID + "/backups"
	parentID          = "{parentID}"
)

var createBackupReq = `
{
    "backup":{
        "instance":"{instanceID}",
        "incremental":0,
        "name":"testingbackup"
    }
}
`

var expectBackup = backup.Backup{
	Created: "2019-10-30T12:30:00",
	Datastore: datastores.DatastorePartial{
		Type:      "mysql",
		Version:   "5.6",
		VersionID: "b00000b0-00b0-0b00-00b0-000b000000bb"},
	Description: "My Backup",
	ID:          backupID,
	Name:        "testbackup",
	InstanceID:  instanceID,
	LocationRef: "http://localhost/path/to/backup",
	ParentID:    parentID,
	Size:        0.14,
	Status:      "COMPLETED",
	Updated:     "2019-10-31T12:30:00",
}

var oneBackup = `
{
	"created": "2019-10-30T12:30:00",
	"datastore": {
		"type": "mysql",
		"version": "5.6",
		"version_id": "b00000b0-00b0-0b00-00b0-000b000000bb"
	},
	"description": "My Backup",
	"id": "{backupID}",
	"instance_id": "{instanceID}",
	"locationRef": "http://localhost/path/to/backup",
	"name": "testbackup",
	"parent_id": "{parentID}",
	"size": 0.14,
	"status": "COMPLETED",
	"updated": "2019-10-31T12:30:00"
}`

var (
	listBackupResp = fmt.Sprintf(`{"backups":[%s]}`, oneBackup)
	createResp     = fmt.Sprintf(`{"backup": %s}`, oneBackup)
	getResp        = createResp
)

func HandleCreate(t *testing.T) {
	fixture.SetupHandler(t, baseURL, "POST", createBackupReq, createResp, 202)
}

func HandleList(t *testing.T) {
	fixture.SetupHandler(t, baseURL, "GET", "", listBackupResp, 200)
}

func HandleDelete(t *testing.T) {
	fixture.SetupHandler(t, resURL, "DELETE", "", "", 202)
}

func HandleGet(t *testing.T) {
	fixture.SetupHandler(t, resURL, "GET", "", getResp, 200)
}
func HandleListByInstance(t *testing.T) {
	fixture.SetupHandler(t, instanceBackupURL, "GET", "", listBackupResp, 200)
}
