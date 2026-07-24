package shareattach

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List returns a Pager that allows you to iterate over a collection of
// ShareAttachments.
func List(client *gophercloud.ServiceClient, serverID string) pagination.Pager {
	return pagination.NewPager(client, listURL(client, serverID), func(r pagination.PageResult) pagination.Page {
		return ShareAttachmentPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	ToShareAttachmentCreateMap() (map[string]any, error)
}

// CreateOpts specifies share attachment creation parameters.
type CreateOpts struct {
	// ShareID is the UUID of the Manila share to attach to the instance.
	ShareID string `json:"share_id" required:"true"`

	// Tag is the device tag used to mount the share within the instance.
	Tag string `json:"tag,omitempty"`
}

// ToShareAttachmentCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToShareAttachmentCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "share")
}

// Create requests the creation of a new share attachment on the server.
func Create(ctx context.Context, client *gophercloud.ServiceClient, serverID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareAttachmentCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client, serverID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get returns public data about a previously created ShareAttachment.
func Get(ctx context.Context, client *gophercloud.ServiceClient, serverID, shareID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, serverID, shareID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete requests the deletion of a previously stored ShareAttachment from the server.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, serverID, shareID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, serverID, shareID), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
