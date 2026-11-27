package tokens

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
)

// CreateOptsBuilder builds a complete token creation body.
type CreateOptsBuilder interface {
	ToTokenCreateMap() (map[string]any, error)
}

// CreateOptsHeadersBuilder optionally supplies token creation headers.
type CreateOptsHeadersBuilder interface {
	ToTokenCreateHeaders(method, url string) (map[string]string, error)
}

// Create authenticates to the identity service and attempts to acquire a Token.
// Generally, applications should use the auth package to authenticate and
// configure a provider client instead of calling this operation directly.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	body, err := opts.ToTokenCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var headers map[string]string
	if optsWithHeaders, ok := opts.(CreateOptsHeadersBuilder); ok {
		headers, err = optsWithHeaders.ToTokenCreateHeaders(http.MethodPost, CreateURL(client))
		if err != nil {
			r.Err = err
			return
		}
	}
	requestOpts := &gophercloud.RequestOpts{
		JSONBody:     body,
		MoreHeaders:  headers,
		OmitHeaders:  []string{"X-Auth-Token", "Authorization"},
		OkCodes:      []int{200, 203},
		JSONResponse: &r.Body,
	}
	resp, err := client.Request(ctx, http.MethodPost, CreateURL(client), requestOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get validates and retrieves information for user's token.
func Get(ctx context.Context, client *gophercloud.ServiceClient, token string) (r GetResult) {
	resp, err := client.Get(ctx, GetURL(client, token), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
