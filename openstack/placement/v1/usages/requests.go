package usages

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
)

// GetOptsBuilder allows extensions to add additional parameters to the
// Get request.
type GetOptsBuilder interface {
	ToUsagesGetQuery() (string, error)
}

// GetOpts specifies the query parameters for retrieving total usages.
//
// This requires microversion 1.9 or later.
type GetOpts struct {
	// ProjectID is required: only usages for this project are returned.
	ProjectID string `q:"project_id"`

	// UserID is optional: when set, only usages for this user within the
	// project are returned.
	UserID string `q:"user_id,omitempty"`

	// ConsumerType is optional: when set, results are filtered to this consumer type.
	// Available from microversion 1.38.
	ConsumerType string `q:"consumer_type,omitempty"`
}

// ToUsagesGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToUsagesGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// Get retrieves the total resource usages for a project (and optionally a user).
//
// Requires microversion 1.9 or later.
func Get(ctx context.Context, client *gophercloud.ServiceClient, opts GetOptsBuilder) (r GetResult) {
	url := getURL(client)
	if opts != nil {
		query, err := opts.ToUsagesGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
