package snapshots

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSnapshotCreateMap() (map[string]any, error)
}

// CreateOpts contains options for creating a Snapshot. This object is passed to
// the snapshots.Create function. For more information about these parameters,
// see the Snapshot object.
type CreateOpts struct {
	VolumeID    string            `json:"volume_id" required:"true"`
	Force       bool              `json:"force,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ToSnapshotCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToSnapshotCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "snapshot")
}

// Create will create a new Snapshot based on the values in CreateOpts. To
// extract the Snapshot object from the response, call the Extract method on the
// CreateResult.
func Create(ctx context.Context, client gophercloud.Client, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSnapshotCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete the existing Snapshot with the provided ID.
func Delete(ctx context.Context, client gophercloud.Client, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves the Snapshot with the provided ID. To extract the Snapshot
// object from the response, call the Extract method on the GetResult.
func Get(ctx context.Context, client gophercloud.Client, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToSnapshotListQuery() (string, error)
}

// ListOpts hold options for listing Snapshots. It is passed to the
// snapshots.List function.
type ListOpts struct {
	// AllTenants will retrieve snapshots of all tenants/projects.
	AllTenants bool `q:"all_tenants"`

	// Name will filter by the specified snapshot name.
	Name string `q:"name"`

	// Status will filter by the specified status.
	Status string `q:"status"`

	// TenantID will filter by a specific tenant/project ID.
	// Setting AllTenants is required to use this.
	TenantID string `q:"project_id"`

	// VolumeID will filter by a specified volume ID.
	VolumeID string `q:"volume_id"`
}

// ToSnapshotListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSnapshotListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns Snapshots optionally limited by the conditions provided in
// ListOpts.
func List(client gophercloud.Client, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToSnapshotListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SnapshotPage{pagination.SinglePageBase(r)}
	})
}

// UpdateMetadataOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateMetadataOptsBuilder interface {
	ToSnapshotUpdateMetadataMap() (map[string]any, error)
}

// UpdateMetadataOpts contain options for updating an existing Snapshot. This
// object is passed to the snapshots.Update function. For more information
// about the parameters, see the Snapshot object.
type UpdateMetadataOpts struct {
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ToSnapshotUpdateMetadataMap assembles a request body based on the contents of
// an UpdateMetadataOpts.
func (opts UpdateMetadataOpts) ToSnapshotUpdateMetadataMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// UpdateMetadata will update the Snapshot with provided information. To
// extract the updated Snapshot from the response, call the ExtractMetadata
// method on the UpdateMetadataResult.
func UpdateMetadata(ctx context.Context, client gophercloud.Client, id string, opts UpdateMetadataOptsBuilder) (r UpdateMetadataResult) {
	b, err := opts.ToSnapshotUpdateMetadataMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateMetadataURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
