package ec2tokens

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// Create authenticates and generates a new token from EC2 credentials.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts auth.AuthOptionsBuilderEC2) (r tokens.CreateResult) {
	request, err := auth.NewRequestEC2(opts)
	if err != nil {
		r.Err = err
		return
	}

	request.JSONResponse = &r.Body
	resp, err := c.Request(ctx, http.MethodPost, ec2tokensURL(c), request)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ValidateS3Token authenticates an S3 request using EC2 credentials. Doesn't
// generate a new token ID, but returns a tokens.CreateResult.
func ValidateS3Token(ctx context.Context, c *gophercloud.ServiceClient, opts auth.AuthOptionsBuilderEC2) (r tokens.CreateResult) {
	request, err := auth.NewRequestS3(opts)
	if err != nil {
		r.Err = err
		return
	}

	request.JSONResponse = &r.Body
	resp, err := c.Request(ctx, http.MethodPost, s3tokensURL(c), request)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
