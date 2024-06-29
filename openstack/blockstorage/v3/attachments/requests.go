package attachments

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAttachmentCreateMap() (map[string]any, error)
}

// CreateOpts contains options for creating a Volume attachment. This object is
// passed to the Create function. For more information about these parameters,
// see the Attachment object.
type CreateOpts struct {
	// VolumeUUID is the UUID of the Cinder volume to create the attachment
	// record for.
	VolumeUUID string `json:"volume_uuid"`
	// InstanceUUID is the ID of the Server to create the attachment for.
	// When attaching to a Nova Server this is the Nova Server (Instance)
	// UUID.
	InstanceUUID string `json:"instance_uuid"`
	// Connector is an optional map containing all of the needed atachment
	// information for exmaple initiator IQN, etc.
	Connector map[string]any `json:"connector,omitempty"`
	// Mode is an attachment mode. Acceptable values are read-only ('ro')
	// and read-and-write ('rw'). Available only since 3.54 microversion.
	// For APIs from 3.27 till 3.53 use Connector["mode"] = "rw|ro".
	Mode string `json:"mode,omitempty"`
}

// ToAttachmentCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToAttachmentCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "attachment")
}

// Create will create a new Attachment based on the values in CreateOpts. To
// extract the Attachment object from the response, call the Extract method on
// the CreateResult.
func Create(ctx context.Context, client gophercloud.Client, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAttachmentCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete the existing Attachment with the provided ID.
func Delete(ctx context.Context, client gophercloud.Client, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves the Attachment with the provided ID. To extract the Attachment
// object from the response, call the Extract method on the GetResult.
func Get(ctx context.Context, client gophercloud.Client, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToAttachmentListQuery() (string, error)
}

// ListOpts holds options for listing Attachments. It is passed to the attachments.List
// function.
type ListOpts struct {
	// AllTenants will retrieve attachments of all tenants/projects.
	AllTenants bool `q:"all_tenants"`

	// Status will filter by the specified status.
	Status string `q:"status"`

	// ProjectID will filter by a specific tenant/project ID.
	ProjectID string `q:"project_id"`

	// VolumeID will filter by a specific volume ID.
	VolumeID string `q:"volume_id"`

	// InstanceID will filter by a specific instance ID.
	InstanceID string `q:"instance_id"`

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

// ToAttachmentListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAttachmentListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns Attachments optionally limited by the conditions provided in
// ListOpts.
func List(client gophercloud.Client, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToAttachmentListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AttachmentPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToAttachmentUpdateMap() (map[string]any, error)
}

// UpdateOpts contain options for updating an existing Attachment.
// This is used to finalize an attachment that was created without a
// connector (reserve).
type UpdateOpts struct {
	Connector map[string]any `json:"connector"`
}

// ToAttachmentUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToAttachmentUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "attachment")
}

// Update will update the Attachment with provided information. To extract the
// updated Attachment from the response, call the Extract method on the
// UpdateResult.
func Update(ctx context.Context, client gophercloud.Client, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAttachmentUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Complete will complete an attachment for a cinder volume.
// Available starting in the 3.44 microversion.
func Complete(ctx context.Context, client gophercloud.Client, id string) (r CompleteResult) {
	b := map[string]any{
		"os-complete": nil,
	}
	resp, err := client.Post(ctx, completeURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
