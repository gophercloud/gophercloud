package ec2credentials

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// List enumerates the Credentials to which the current token has access.
func List(client gophercloud.Client, userID string) pagination.Pager {
	url := listURL(client, userID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return CredentialPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single EC2 credential by ID.
func Get(ctx context.Context, client gophercloud.Client, userID string, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, userID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOpts provides options used to create an EC2 credential.
type CreateOpts struct {
	// TenantID is the project ID scope of the EC2 credential.
	TenantID string `json:"tenant_id" required:"true"`
}

// ToCredentialCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToCredentialCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create creates a new EC2 Credential.
func Create(ctx context.Context, client gophercloud.Client, userID string, opts CreateOpts) (r CreateResult) {
	b, err := opts.ToCredentialCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client, userID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes an EC2 credential.
func Delete(ctx context.Context, client gophercloud.Client, userID string, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, userID, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
