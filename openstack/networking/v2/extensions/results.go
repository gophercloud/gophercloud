package extensions

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type Extension struct {
	Updated     string        `json:"updated"`
	Name        string        `json:"name"`
	Links       []interface{} `json:"links"`
	Namespace   string        `json:"namespace"`
	Alias       string        `json:"alias"`
	Description string        `json:"description"`
}

func ExtractExtensions(page gophercloud.Page) ([]Extension, error) {
	var resp struct {
		Extensions []Extension `mapstructure:"extensions"`
	}

	err := mapstructure.Decode(page.(gophercloud.LinkedPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Extensions, nil
}
