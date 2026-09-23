package tokens

import (
	"context"
	"maps"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

const xSubjectTokenHeader = "X-Subject-Token"

// Create authenticates and either generates a new token, or changes the Scope
// of an existing token.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts auth.AuthOptionsBuilderV3) (r CreateResult) {
	request, err := auth.NewRequestV3(opts, auth.WithTokenURL(tokenURL(c)))
	if err != nil {
		r.Err = err
		return
	}

	request.JSONResponse = &r.Body
	resp, err := c.Request(ctx, http.MethodPost, tokenURL(c), request)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetOptsBuilder allows extensions to add additional parameters to
// the Get request.
type GetOptsBuilder interface {
	ToTokenGetParams() (map[string]string, error)
}

// GetOpts provides options for the Get request.
type GetOpts struct {
	// AccessRulesVersion specifies the OpenStack-Identity-Access-Rules header
	// version. This is required when getting tokens that were created using
	// application credentials with access rules. Versions less than 1 will cause
	// Keystone to ignore access rules validation.
	AccessRulesVersion string `h:"OpenStack-Identity-Access-Rules"`
}

// ToTokenGetParams formats GetOpts into request headers.
func (opts GetOpts) ToTokenGetParams() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// Get validates and retrieves information about another token.
func Get(ctx context.Context, c *gophercloud.ServiceClient, token string, opts GetOptsBuilder) (r GetResult) {
	h := map[string]string{xSubjectTokenHeader: token}
	if opts != nil {
		b, err := opts.ToTokenGetParams()
		if err != nil {
			r.Err = err
			return
		}
		maps.Copy(h, b)
	}

	resp, err := c.Get(ctx, tokenURL(c), &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{200, 203},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ValidateOptsBuilder allows extensions to add additional parameters to
// the Validate request.
type ValidateOptsBuilder interface {
	ToTokenValidateParams() (map[string]string, error)
}

// ValidateOpts provides options for the Validate request.
type ValidateOpts struct {
	// AccessRulesVersion specifies the OpenStack-Identity-Access-Rules header
	// version. This is required when validating tokens that were created using
	// application credentials with access rules. Versions less than 1 will cause
	// Keystone to ignore access rules validation.
	AccessRulesVersion string `h:"OpenStack-Identity-Access-Rules"`
}

// ToTokenValidateParams formats ValidateOpts into request headers.
func (opts ValidateOpts) ToTokenValidateParams() (map[string]string, error) {
	return gophercloud.BuildHeaders(opts)
}

// Validate determines if a specified token is valid or not.
func Validate(ctx context.Context, c *gophercloud.ServiceClient, token string, opts ValidateOptsBuilder) (bool, error) {
	h := map[string]string{xSubjectTokenHeader: token}
	if opts != nil {
		b, err := opts.ToTokenValidateParams()
		if err != nil {
			return false, err
		}
		maps.Copy(h, b)
	}

	resp, err := c.Head(ctx, tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{200, 204, 404},
	})
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 200 || resp.StatusCode == 204, nil
}

// Revoke immediately makes specified token invalid.
func Revoke(ctx context.Context, c *gophercloud.ServiceClient, token string) (r RevokeResult) {
	resp, err := c.Delete(ctx, tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{xSubjectTokenHeader: token},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
