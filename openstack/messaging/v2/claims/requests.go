package claims

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToClaimCreateRequest() (map[string]any, string, error)
}

// CreateOpts params to be used with Create.
type CreateOpts struct {
	// Sets the TTL for the claim. When the claim expires un-deleted messages will be able to be claimed again.
	TTL int `json:"ttl,omitempty"`

	// Sets the Grace period for the claimed messages. The server extends the lifetime of claimed messages
	// to be at least as long as the lifetime of the claim itself, plus the specified grace period.
	Grace int `json:"grace,omitempty"`

	// Set the limit of messages returned by create.
	Limit int `q:"limit" json:"-"`
}

// ToClaimCreateRequest assembles a body and URL for a Create request based on
// the contents of a CreateOpts.
func (opts CreateOpts) ToClaimCreateRequest() (map[string]any, string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return nil, q.String(), err
	}

	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return b, "", err
	}
	return b, q.String(), err
}

// Create creates a Claim that claims messages on a specified queue.
func Create(ctx context.Context, client *gophercloud.ServiceClient, queueName string, opts CreateOptsBuilder) (r CreateResult) {
	b, q, err := opts.ToClaimCreateRequest()
	if err != nil {
		r.Err = err
		return
	}

	url := createURL(client, queueName)
	if q != "" {
		url += q
	}

	resp, err := client.Post(ctx, url, b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get queries the specified claim for the specified queue.
func Get(ctx context.Context, client *gophercloud.ServiceClient, queueName string, claimID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, queueName, claimID), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToClaimUpdateMap() (map[string]any, error)
}

// UpdateOpts implements UpdateOpts.
type UpdateOpts struct {
	// Update the TTL for the specified Claim.
	TTL int `json:"ttl,omitempty"`

	// Update the grace period for Messages in a specified Claim.
	Grace int `json:"grace,omitempty"`
}

// ToClaimUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToClaimUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Update will update the options for a specified claim.
func Update(ctx context.Context, client *gophercloud.ServiceClient, queueName string, claimID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClaimUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}
	resp, err := client.Patch(ctx, updateURL(client, queueName, claimID), &b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete a Claim for a specified Queue.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, queueName string, claimID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, queueName, claimID), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
