package backups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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
	VolumeID    string            `json:"volume_id" required:"true"`
	Force       bool              `json:"force,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Container   string            `json:"container,omitempty"`
	Incremental bool              `json:"incremental,omitempty"`
	SnapshotID  string            `json:"snapshot_id,omitempty"`
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
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

// Delete will delete the existing Backup with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get retrieves the Backup with the provided ID. To extract the Backup
// object from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToBackupListQuery() (string, error)
}

// ListOpts hold options for listing Backups. It is passed to the
// backups.List function.
type ListOpts struct {
	// AllTenants will retrieve backups of all tenants/projects.
	AllTenants bool `q:"all_tenants"`

	// Name will filter by the specified backup name.
	Name string `q:"name"`

	// Status will filter by the specified status.
	Status string `q:"status"`

	// TenantID will filter by a specific tenant/project ID.
	// Setting AllTenants is required to use this.
	TenantID string `q:"project_id"`

	// VolumeID will filter by a specified volume ID.
	VolumeID string `q:"volume_id"`
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
		return BackupPage{pagination.SinglePageBase(r)}
	})
}

// UpdateMetadataOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateMetadataOptsBuilder interface {
	ToBackupUpdateMetadataMap() (map[string]interface{}, error)
}

// UpdateMetadataOpts contain options for updating an existing Backup. This
// object is passed to the backups.Update function. For more information
// about the parameters, see the Backup object.
type UpdateMetadataOpts struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToBackupUpdateMetadataMap assembles a request body based on the contents of
// an UpdateMetadataOpts.
func (opts UpdateMetadataOpts) ToBackupUpdateMetadataMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// UpdateMetadata will update the Backup with provided information. To
// extract the updated Backup from the response, call the ExtractMetadata
// method on the UpdateMetadataResult.
func UpdateMetadata(client *gophercloud.ServiceClient, id string, opts UpdateMetadataOptsBuilder) (r UpdateMetadataResult) {
	b, err := opts.ToBackupUpdateMetadataMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateMetadataURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// IDFromName is a convienience function that returns a backup's ID given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""
	pages, err := List(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractBackups(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{Name: name, ResourceType: "backup"}
	case 1:
		return id, nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "backup"}
	}
}
