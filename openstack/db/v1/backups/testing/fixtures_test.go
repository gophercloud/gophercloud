package testing

import (
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/backups"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
)

var (
	timestamp  = "2014-10-30T12:30:00"
	timeVal, _ = time.Parse(gophercloud.RFC3339NoZ, timestamp)
)

const singleBackupJSON = `
{
  "created": "2014-10-30T12:30:00",
  "datastore": {
    "type": "mysql",
    "version": "5.5",
    "version_id": "b00000b0-00b0-0b00-00b0-000b000000bb"
  },
  "description": "My Backup",
  "id": "a9832168-7541-4536-b8d9-a8a9b79cf1b4",
  "instance_id": "44b277eb-39be-4921-be31-3d61b43651d7",
  "locationRef": "http://localhost/path/to/backup",
  "name": "snapshot",
  "parent_id": "",
  "project_id": "922b47766bcb448f83a760358337f2b4",
  "size": 0.14,
  "status": "COMPLETED",
  "updated": "2014-10-30T12:30:00",
  "storage_driver": "swift"
}
`

var (
	ListBackupsJSON  = fmt.Sprintf(`{"backups": [%s]}`, singleBackupJSON)
	GetBackupJSON    = fmt.Sprintf(`{"backup": %s}`, singleBackupJSON)
	CreateBackupJSON = fmt.Sprintf(`{"backup": %s}`, singleBackupJSON)
)

var CreateBackupReq = `
{
  "backup": {
    "description": "My Backup",
    "incremental": 0,
    "instance": "44b277eb-39be-4921-be31-3d61b43651d7",
    "name": "snapshot",
    "storage_driver": "swift"
  }
}
`

var ExampleBackup = backups.Backup{
	Created: timeVal,
	Datastore: datastores.DatastorePartial{
		Type:      "mysql",
		Version:   "5.5",
		VersionID: "b00000b0-00b0-0b00-00b0-000b000000bb",
	},
	Description:   "My Backup",
	ID:            "a9832168-7541-4536-b8d9-a8a9b79cf1b4",
	InstanceID:    "44b277eb-39be-4921-be31-3d61b43651d7",
	LocationRef:   "http://localhost/path/to/backup",
	Name:          "snapshot",
	ProjectID:     "922b47766bcb448f83a760358337f2b4",
	Size:          0.14,
	Status:        "COMPLETED",
	StorageDriver: "swift",
	Updated:       timeVal,
}
