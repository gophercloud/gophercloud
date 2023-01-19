package backups

import (
	"github.com/bizflycloud/gophercloud"
	"github.com/bizflycloud/gophercloud/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToBackupCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Backup. This object is passed to
// the backups.Create function. For more information about these parameters,
// see the Backup object.
type CreateOpts struct {
	// VolumeID is the ID of the volume to create the backup from.
	VolumeID string `json:"volume_id" required:"true"`

	// Force will force the creation of a backup regardless of the
	//volume's status.
	Force bool `json:"force,omitempty"`

	// Name is the name of the backup.
	Name string `json:"name,omitempty"`

	// Description is the description of the backup.
	Description string `json:"description,omitempty"`

	// Metadata is metadata for the backup.
	// Requires microversion 3.43 or later.
	Metadata map[string]string `json:"metadata,omitempty"`

	// Container is a container to store the backup.
	Container string `json:"container,omitempty"`

	// Incremental is whether the backup should be incremental or not.
	Incremental bool `json:"incremental,omitempty"`

	// SnapshotID is the ID of a snapshot to backup.
	SnapshotID string `json:"snapshot_id,omitempty"`

	// AvailabilityZone is an availability zone to locate the volume or snapshot.
	// Requires microversion 3.51 or later.
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

// ToBackupCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToBackupCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "backup")
}

// Create will create a new Backup based on the values in CreateOpts. To
// extract the Backup object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBackupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete the existing Backup with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves the Backup with the provided ID. To extract the Backup
// object from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToBackupListQuery() (string, error)
}

type ListOpts struct {
	// AllTenants will retrieve backups of all tenants/projects.
	AllTenants bool `q:"all_tenants"`

	// Name will filter by the specified backup name.
	// This does not work in later microversions.
	Name string `q:"name"`

	// Status will filter by the specified status.
	// This does not work in later microversions.
	Status string `q:"status"`

	// TenantID will filter by a specific tenant/project ID.
	// Setting AllTenants is required to use this.
	TenantID string `q:"project_id"`

	// VolumeID will filter by a specified volume ID.
	// This does not work in later microversions.
	VolumeID string `q:"volume_id"`

	// Comma-separated list of sort keys and optional sort directions in the
	// form of <key>[:<direction>].
	Sort string `q:"sort"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`
}

// ToBackupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToBackupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns Backups optionally limited by the conditions provided in
// ListOpts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToBackupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListDetailOptsBuilder allows extensions to add additional parameters to the ListDetail
// request.
type ListDetailOptsBuilder interface {
	ToBackupListDetailQuery() (string, error)
}

type ListDetailOpts struct {
	// AllTenants will retrieve backups of all tenants/projects.
	AllTenants bool `q:"all_tenants"`

	// Comma-separated list of sort keys and optional sort directions in the
	// form of <key>[:<direction>].
	Sort string `q:"sort"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`

	// True to include `count` in the API response, supported from version 3.45
	WithCount bool `q:"with_count"`
}

// ToBackupListDetailQuery formats a ListDetailOpts into a query string.
func (opts ListDetailOpts) ToBackupListDetailQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail returns more detailed information about Backups optionally
// limited by the conditions provided in ListDetailOpts.
func ListDetail(client *gophercloud.ServiceClient, opts ListDetailOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToBackupListDetailQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToBackupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing Backup.
type UpdateOpts struct {
	// Name is the name of the backup.
	Name *string `json:"name,omitempty"`

	// Description is the description of the backup.
	Description *string `json:"description,omitempty"`

	// Metadata is metadata for the backup.
	// Requires microversion 3.43 or later.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ToBackupUpdateMap assembles a request body based on the contents of
// an UpdateOpts.
func (opts UpdateOpts) ToBackupUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update will update the Backup with provided information. To extract
// the updated Backup from the response, call the Extract method on the
// UpdateResult.
// Requires microversion 3.9 or later.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToBackupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RestoreOpts contains options for restoring a Backup. This object is passed to
// the backups.RestoreFromBackup function.
type RestoreOpts struct {
	// VolumeID is the ID of the existing volume to restore the backup to.
	VolumeID string `json:"volume_id,omitempty"`

	// Name is the name of the new volume to restore the backup to.
	Name string `json:"name,omitempty"`
}

// ToRestoreMap assembles a request body based on the contents of a
// RestoreOpts.
func (opts RestoreOpts) ToRestoreMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "restore")
}

// RestoreFromBackup will restore a Backup to a volume based on the values in
// RestoreOpts. To extract the Restore object from the response, call the
// Extract method on the RestoreResult.
func RestoreFromBackup(client *gophercloud.ServiceClient, id string, opts RestoreOpts) (r RestoreResult) {
	b, err := opts.ToRestoreMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(restoreURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Export will export a Backup information. To extract the Backup export record
// object from the response, call the Extract method on the ExportResult.
func Export(client *gophercloud.ServiceClient, id string) (r ExportResult) {
	resp, err := client.Get(exportURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ImportOpts contains options for importing a Backup. This object is passed to
// the backups.ImportBackup function.
type ImportOpts BackupRecord

// ToBackupImportMap assembles a request body based on the contents of a
// ImportOpts.
func (opts ImportOpts) ToBackupImportMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "backup-record")
}

// Import will import a Backup data to a backup based on the values in
// ImportOpts. To extract the Backup object from the response, call the
// Extract method on the ImportResult.
func Import(client *gophercloud.ServiceClient, opts ImportOpts) (r ImportResult) {
	b, err := opts.ToBackupImportMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(importURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
