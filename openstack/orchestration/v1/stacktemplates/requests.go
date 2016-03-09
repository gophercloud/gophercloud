package stacktemplates

import "github.com/gophercloud/gophercloud"

// Get retreives data for the given stack template.
func Get(c *gophercloud.ServiceClient, stackName, stackID string) GetResult {
	var res GetResult
	_, res.Err = c.Request("GET", getURL(c, stackName, stackID), &gophercloud.RequestOpts{
		JSONResponse: &res.Body,
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
	Template    string
	TemplateURL string
}

// ToStackTemplateValidateMap assembles a request body based on the contents of a ValidateOpts.
func (opts ValidateOpts) ToStackTemplateValidateMap() (map[string]interface{}, error) {
	vo := make(map[string]interface{})
	if opts.Template != "" {
		vo["template"] = opts.Template
		return vo, nil
	}
	if opts.TemplateURL != "" {
		vo["template_url"] = opts.TemplateURL
		return vo, nil
	}
	err := gophercloud.ErrMissingInput{}
	err.Argument = "stacktemplates.ValidateOpts.Template/stacktemplates.ValidateOpts.TemplateURL"
	err.Info = "One of Template or TemplateURL is required."
	return nil, err
}

// Validate validates the given stack template.
func Validate(c *gophercloud.ServiceClient, opts ValidateOptsBuilder) ValidateResult {
	var res ValidateResult

	reqBody, err := opts.ToStackTemplateValidateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = c.Post(validateURL(c), reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return res
}
