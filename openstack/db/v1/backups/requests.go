package backups

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToBackupListQuery() (string, error)
}

// ListOpts represents options for listing database backups.
type ListOpts struct {
	// Return the list of backups for a particular database instance.
	InstanceID string `q:"instance_id"`
	// Return the list of backups for all projects. Admin only.
	AllProjects bool `q:"all_projects"`
	// Return backups for a particular datastore.
	Datastore string `q:"datastore"`
	// Return backups for a particular project. Admin only.
	ProjectID string `q:"project_id"`
}

// ToBackupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToBackupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List will list database backups.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToBackupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.SinglePageBase(r)}
	})
}

// ListByInstance will list backups for a database instance.
func ListByInstance(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {
	return pagination.NewPager(client, instanceURL(client, instanceID), func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder is the top-level interface for create backup options.
type CreateOptsBuilder interface {
	ToBackupCreateMap() (map[string]any, error)
}

// RestoreFromOpts contains information for restoring a backup from a remote
// location.
type RestoreFromOpts struct {
	RemoteLocation          string  `json:"remote_location,omitempty"`
	LocalDatastoreVersionID string  `json:"local_datastore_version_id,omitempty"`
	Size                    float64 `json:"size,omitempty"`
}

// CreateOpts represents options for creating a database backup.
type CreateOpts struct {
	// Name of the backup.
	Name string `json:"name" required:"true"`
	// ID of the instance to create a backup for.
	InstanceID string `json:"instance,omitempty"`
	// ID of the parent backup to perform an incremental backup from.
	ParentID string `json:"parent_id,omitempty"`
	// Set to 1 to create an incremental backup based on the last full backup.
	Incremental *int `json:"incremental,omitempty"`
	// An optional description for the backup.
	Description string `json:"description,omitempty"`
	// User-defined Swift container name.
	SwiftContainer string `json:"swift_container,omitempty"`
	// Information needed to restore a remote backup.
	RestoreFrom *RestoreFromOpts `json:"restore_from,omitempty"`
	// Backup storage driver.
	StorageDriver string `json:"storage_driver,omitempty"`
}

// ToBackupCreateMap converts a CreateOpts struct into a request body.
func (opts CreateOpts) ToBackupCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "backup")
}

// Create creates a database backup.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBackupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves details for a database backup.
func Get(ctx context.Context, client *gophercloud.ServiceClient, backupID string) (r GetResult) {
	resp, err := client.Get(ctx, resourceURL(client, backupID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a database backup.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, backupID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, resourceURL(client, backupID), &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
