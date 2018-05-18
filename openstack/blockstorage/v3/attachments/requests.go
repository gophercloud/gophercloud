package attachments

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAttachmentCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a Volume. This object is passed to
// the volumes.Create function. For more information about these parameters,
// see the Volume object.
type CreateOpts struct {
	// VolumeID is the UUID of the Cinder volume to create the attachment record for
	VolumeID string `json:"volume_uuid"`
	// ServerID is the ID of the Server to create the attachment for.  When attaching to a
	// Nova Server this is the Nova Server (Instance) UUID
	ServerID string `json:"instance_uuid"`
	// Connector is optional and is a map containing all of the needed atachment information
	// for exmaple initiator iqn etc
	Connector map[string]string `json:"connector,omitempty"`
}

// ToAttachmentCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToAttachmentCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "attachment")
}

// Create will create a new Attachment based on the values in CreateOpts. To extract
// the Attachment object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAttachmentCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

// Delete will delete the existing Attachment with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get retrieves the Volume with the provided ID. To extract the Volume object
// from the response, call the Extract method on the GetResult.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
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

	// TenantID will filter by a specific tenant/project ID.
	// Setting AllTenants is required for this.
	TenantID string `q:"project_id"`

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

// List returns Attachemnts optionally limited by the conditions provided in ListOpts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
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
	ToAttachmentUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contain options for updating an existing Attachment.
// This is used to finalize an attachment that was created without a
// connector (reserve)
type UpdateOpts struct {
	Connector map[string]string `json:"connector,omitempty"`
}

// ToAttachmentUpdateMap assembles a request body based on the contents of an
// UpdateOpts.
func (opts UpdateOpts) ToAttachmentUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "attachment")
}

// Update will update the Volume with provided information. To extract the updated
// Attachment from the response, call the Extract method on the UpdateResult.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAttachmentUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
