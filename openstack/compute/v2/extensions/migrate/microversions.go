package migrate

import "github.com/gophercloud/gophercloud"

// LiveMigrateOpts specifies parameters of live migrate action.
type LiveMigrate225Opts struct {
	// The host to which to migrate the server.
	// If this parameter is None, the scheduler chooses a host.
	Host *string `json:"host"`

	// Set to True to migrate local disks by using block migration.
	// If the source or destination host uses shared storage and you set
	// this value to True, the live migration fails.
	BlockMigration string `json:"block_migration"`

	// Set to True to enable over commit when the destination host is checked
	// for available disk space. Set to False to disable over commit. This setting
	// affects only the libvirt virt driver.
	DiskOverCommit *bool `json:"disk_over_commit,omitempty"`
}

// ToLiveMigrateMap constructs a request body from LiveMigrateOpts.
func (opts LiveMigrate225Opts) ToLiveMigrateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "os-migrateLive")
}
