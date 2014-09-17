package extensions

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

type Extension struct {
	Updated     string        `json:"updated"`
	Name        string        `json:"name"`
	Links       []interface{} `json:"links"`
	Namespace   string        `json:"namespace"`
	Alias       string        `json:"alias"`
	Description string        `json:"description"`
}

type ExtensionPage struct {
	pagination.SinglePageBase
}

func (r ExtensionPage) IsEmpty() (bool, error) {
	is, err := ExtractExtensions(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractExtensions(page pagination.Page) ([]Extension, error) {
	var resp struct {
		Extensions []Extension `mapstructure:"extensions"`
	}

	err := mapstructure.Decode(page.(ExtensionPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Extensions, nil
}
