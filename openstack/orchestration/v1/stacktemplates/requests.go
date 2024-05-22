package stacktemplates

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// Get retreives data for the given stack template.
func Get(ctx context.Context, c *gophercloud.ServiceClient, stackName, stackID string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, stackName, stackID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ValidateOptsBuilder describes struct types that can be accepted by the Validate call.
// The ValidateOpts struct in this package does.
type ValidateOptsBuilder interface {
	ToStackTemplateValidateMap() (map[string]any, error)
}

// ValidateOpts specifies the template validation parameters.
type ValidateOpts struct {
	Template    string `json:"template" or:"TemplateURL"`
	TemplateURL string `json:"template_url" or:"Template"`
}

// ToStackTemplateValidateMap assembles a request body based on the contents of a ValidateOpts.
func (opts ValidateOpts) ToStackTemplateValidateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Validate validates the given stack template.
func Validate(ctx context.Context, c *gophercloud.ServiceClient, opts ValidateOptsBuilder) (r ValidateResult) {
	b, err := opts.ToStackTemplateValidateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, validateURL(c), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
