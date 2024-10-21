package quotas

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToQuotaCreateMap() (map[string]any, error)
}

// CreateOpts params
type CreateOpts struct {
	ProjectID string `json:"project_id"`
	Resource  string `json:"resource"`
	HardLimit int    `json:"hard_limit"`
}

// ToQuotaCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToQuotaCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create requests the creation of a new quota.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToQuotaCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
