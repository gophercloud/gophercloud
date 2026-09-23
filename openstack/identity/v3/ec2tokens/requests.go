package ec2tokens

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// Create authenticates and generates a new token from EC2 credentials.
func Create(ctx context.Context, c *gophercloud.ServiceClient, request gophercloud.RawRequest) (r tokens.CreateResult) {
	requestOpts := &gophercloud.RequestOpts{
		JSONBody:     request.Body,
		MoreHeaders:  request.Headers,
		OkCodes:      []int{200},
		JSONResponse: &r.Body,
	}
	resp, err := c.Request(ctx, http.MethodPost, ec2tokensURL(c), requestOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ValidateS3Token authenticates an S3 request using EC2 credentials. Doesn't
// generate a new token ID, but returns a tokens.CreateResult.
func ValidateS3Token(ctx context.Context, c *gophercloud.ServiceClient, request gophercloud.RawRequest) (r tokens.CreateResult) {
	requestOpts := &gophercloud.RequestOpts{
		JSONBody:     request.Body,
		MoreHeaders:  request.Headers,
		OkCodes:      []int{200},
		JSONResponse: &r.Body,
	}
	resp, err := c.Request(ctx, http.MethodPost, s3tokensURL(c), requestOpts)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
