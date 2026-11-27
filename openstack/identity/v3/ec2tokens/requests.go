package ec2tokens

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// CreateOptsBuilder builds a complete token creation body.
type CreateOptsBuilder interface {
	ToTokenCreateMap() (map[string]any, error)
}

// CreateOptsHeadersBuilder optionally supplies token creation headers.
type CreateOptsHeadersBuilder interface {
	ToTokenCreateHeaders(method, url string) (map[string]string, error)
}

// Create authenticates and generates a new token from EC2 credentials.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r tokens.CreateResult) {
	body, err := opts.ToTokenCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var headers map[string]string
	if optsWithHeaders, ok := opts.(CreateOptsHeadersBuilder); ok {
		headers, err = optsWithHeaders.ToTokenCreateHeaders(http.MethodPost, ec2tokensURL(c))
		if err != nil {
			r.Err = err
			return
		}
	}
	requestOpts := &gophercloud.RequestOpts{
		JSONBody:     body,
		MoreHeaders:  headers,
		OkCodes:      []int{200},
		JSONResponse: &r.Body,
	}
	resp, err := c.Request(ctx, http.MethodPost, ec2tokensURL(c), requestOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ValidateS3TokenOptsBuilder builds a complete S3 token validation body.
type ValidateS3TokenOptsBuilder interface {
	ToS3TokenValidateMap() (map[string]any, error)
}

// ValidateS3TokenOptsHeadersBuilder optionally supplies S3 token validation headers.
type ValidateS3TokenOptsHeadersBuilder interface {
	ToS3TokenValidateHeaders(method, url string) (map[string]string, error)
}

// ValidateS3Token authenticates an S3 request using EC2 credentials. Doesn't
// generate a new token ID, but returns a tokens.CreateResult.
func ValidateS3Token(ctx context.Context, c *gophercloud.ServiceClient, opts ValidateS3TokenOptsBuilder) (r tokens.CreateResult) {
	body, err := opts.ToS3TokenValidateMap()
	if err != nil {
		r.Err = err
		return
	}
	var headers map[string]string
	if optsWithHeaders, ok := opts.(ValidateS3TokenOptsHeadersBuilder); ok {
		headers, err = optsWithHeaders.ToS3TokenValidateHeaders(http.MethodPost, s3tokensURL(c))
		if err != nil {
			r.Err = err
			return
		}
	}
	requestOpts := &gophercloud.RequestOpts{
		JSONBody:     body,
		MoreHeaders:  headers,
		OkCodes:      []int{200},
		JSONResponse: &r.Body,
	}
	resp, err := c.Request(ctx, http.MethodPost, s3tokensURL(c), requestOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
