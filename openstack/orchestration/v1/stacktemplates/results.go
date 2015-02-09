package stacktemplates

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type Template struct {
	Description         string                 `mapstructure:"description"`
	HeatTemplateVersion string                 `mapstructure:"heat_template_version"`
	Parameters          map[string]interface{} `mapstructure:"parameters"`
	Resources           map[string]interface{} `mapstructure:"resources"`
}

type GetResult struct {
	gophercloud.Result
}

func (r GetResult) Extract() (*Template, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res Template
	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type ValidatedTemplate struct {
	Description string
	Parameters  map[string]interface{}
}

type ValidateResult struct {
	gophercloud.Result
}

func (r ValidateResult) Extract() (*ValidatedTemplate, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res ValidatedTemplate
	if err := mapstructure.Decode(r.Body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
