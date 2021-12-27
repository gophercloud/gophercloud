package migrate

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions"
)

// Migrate will initiate a migration of the instance to another host.
func Migrate(client *gophercloud.ServiceClient, id string) (r MigrateResult) {
	resp, err := client.Post(extensions.ActionURL(client, id), map[string]interface{}{"migrate": nil}, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// LiveMigrateOptsBuilder allows extensions to add additional parameters to the
// LiveMigrate request.
type LiveMigrateOptsBuilder interface {
	ToLiveMigrateMap() (map[string]interface{}, error)
}

// LiveMigrateOpts specifies parameters of live migrate action.
type LiveMigrateOpts struct {
	// The host to which to migrate the server.
	// If this parameter is None, the scheduler chooses a host.
	Host *string `json:"host"`

	// Set to True to migrate local disks by using block migration.
	// If the source or destination host uses shared storage and you set
	// this value to True, the live migration fails.
	BlockMigration interface{} `json:"block_migration,omitempty"`

	// Set to True to enable over commit when the destination host is checked
	// for available disk space. Set to False to disable over commit. This setting
	// affects only the libvirt virt driver.
	DiskOverCommit *bool `json:"disk_over_commit,omitempty"`
}

// ToLiveMigrateMap constructs a request body from LiveMigrateOpts.
func (opts LiveMigrateOpts) ToLiveMigrateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "os-migrateLive")
	if err != nil {
		return nil, err
	}

	switch blockMigration := opts.BlockMigration.(type) {
	case string:
		if blockMigration != "auto" && blockMigration != "True" && blockMigration != "False" {
			return nil, fmt.Errorf("block_migration expected to be one of 'auto', 'True', 'False' but got '%s'", blockMigration)
		}
		// valid
	case *bool:
		// valid
	default:
		return nil, fmt.Errorf("expected block_migration type to be string or *bool, got '%t'", blockMigration)
	}

	return b, nil
}

// LiveMigrate will initiate a live-migration (without rebooting) of the instance to another host.
func LiveMigrate(client *gophercloud.ServiceClient, id string, opts LiveMigrateOptsBuilder) (r MigrateResult) {
	b, err := opts.ToLiveMigrateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(extensions.ActionURL(client, id), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
