package stacktemplates

import (
	"fmt"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// Get retreives data for the given stack template.
func Get(c *gophercloud.ServiceClient, stackName, stackID string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(c, stackName, stackID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// ValidateOptsBuilder describes struct types that can be accepted by the Validate call.
// The ValidateOpts struct in this package does.
type ValidateOptsBuilder interface {
	ToStackTemplateValidateMap() (map[string]interface{}, error)
}

// ValidateOpts specifies the template validation parameters.
type ValidateOpts struct {
	Template    map[string]interface{}
	TemplateURL string
}

// ToStackTemplateValidateMap assembles a request body based on the contents of a ValidateOpts.
func (opts ValidateOpts) ToStackTemplateValidateMap() (map[string]interface{}, error) {
	vo := make(map[string]interface{})
	if opts.Template != nil {
		vo["template"] = opts.Template
		return vo, nil
	}
	if opts.TemplateURL != "" {
		vo["template_url"] = opts.TemplateURL
		return vo, nil
	}
	return vo, fmt.Errorf("One of Template or TemplateURL is required.")
}

// Validate validates the given stack template.
func Validate(c *gophercloud.ServiceClient, opts ValidateOptsBuilder) ValidateResult {
	var res ValidateResult

	reqBody, err := opts.ToStackTemplateValidateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", validateURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}
