package testing

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/backupstrategies"
)

const singleBackupStrategyJSON = `
{
  "project_id": "922b47766bcb448f83a760358337f2b4",
  "instance_id": "0602db72-c63d-11ea-b87c-00224d6b7bc1",
  "backend": "swift",
  "swift_container": "my_trove_backups"
}
`

var (
	ListBackupStrategiesJSON = fmt.Sprintf(`{"backup_strategies": [%s]}`, singleBackupStrategyJSON)
	CreateBackupStrategyJSON = fmt.Sprintf(`{"backup_strategy": %s}`, singleBackupStrategyJSON)
)

var CreateBackupStrategyReq = `
{
  "backup_strategy": {
    "instance_id": "0602db72-c63d-11ea-b87c-00224d6b7bc1",
    "swift_container": "my_trove_backups"
  }
}
`

var ExampleBackupStrategy = backupstrategies.BackupStrategy{
	ProjectID:      "922b47766bcb448f83a760358337f2b4",
	InstanceID:     "0602db72-c63d-11ea-b87c-00224d6b7bc1",
	Backend:        "swift",
	SwiftContainer: "my_trove_backups",
}
