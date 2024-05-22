package policies

import (
	"context"
	"net/url"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

const policyTypeMaxLength = 255

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToPolicyListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Type filters the response by MIME media type
	// of the serialized policy blob.
	Type string `q:"type"`

	// Filters filters the response by custom filters such as
	// 'type__contains=foo'
	Filters map[string]string `q:"-"`
}

// ToPolicyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToPolicyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()
	for k, v := range opts.Filters {
		i := strings.Index(k, "__")
		if i > 0 && i < len(k)-2 {
			params.Add(k, v)
		} else {
			return "", InvalidListFilter{FilterName: k}
		}
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}

// List enumerates the policies to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToPolicyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a policy.
type CreateOpts struct {
	// Type is the MIME media type of the serialized policy blob.
	Type string `json:"type" required:"true"`

	// Blob is the policy rule as a serialized blob.
	Blob []byte `json:"-" required:"true"`

	// Extra is free-form extra key/value pairs to describe the policy.
	Extra map[string]any `json:"-"`
}

// ToPolicyCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToPolicyCreateMap() (map[string]any, error) {
	if len(opts.Type) > policyTypeMaxLength {
		return nil, StringFieldLengthExceedsLimit{
			Field: "type",
			Limit: policyTypeMaxLength,
		}
	}

	b, err := gophercloud.BuildRequestBody(opts, "policy")
	if err != nil {
		return nil, err
	}

	if v, ok := b["policy"].(map[string]any); ok {
		v["blob"] = string(opts.Blob)

		if opts.Extra != nil {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create creates a new Policy.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPolicyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves details on a single policy, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, policyID string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, policyID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToPolicyUpdateMap() (map[string]any, error)
}

// UpdateOpts provides options for updating a policy.
type UpdateOpts struct {
	// Type is the MIME media type of the serialized policy blob.
	Type string `json:"type,omitempty"`

	// Blob is the policy rule as a serialized blob.
	Blob []byte `json:"-"`

	// Extra is free-form extra key/value pairs to describe the policy.
	Extra map[string]any `json:"-"`
}

// ToPolicyUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]any, error) {
	if len(opts.Type) > policyTypeMaxLength {
		return nil, StringFieldLengthExceedsLimit{
			Field: "type",
			Limit: policyTypeMaxLength,
		}
	}

	b, err := gophercloud.BuildRequestBody(opts, "policy")
	if err != nil {
		return nil, err
	}

	if v, ok := b["policy"].(map[string]any); ok {
		if len(opts.Blob) != 0 {
			v["blob"] = string(opts.Blob)
		}

		if opts.Extra != nil {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update updates an existing Role.
func Update(ctx context.Context, client *gophercloud.ServiceClient, policyID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPolicyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, updateURL(client, policyID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a policy.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, policyID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, policyID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
