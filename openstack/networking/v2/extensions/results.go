package extensions

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// Extension is a struct that represents a Neutron extension.
type Extension struct {
	Updated     string        `json:"updated"`
	Name        string        `json:"name"`
	Links       []interface{} `json:"links"`
	Namespace   string        `json:"namespace"`
	Alias       string        `json:"alias"`
	Description string        `json:"description"`
}

// ExtensionPage is the page returned by a pager when traversing over a
// collection of extensions.
type ExtensionPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether an ExtensionPage struct is empty.
func (r ExtensionPage) IsEmpty() (bool, error) {
	is, err := ExtractExtensions(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

// ExtractExtensions accepts a Page struct, specifically an ExtensionPage
// struct, and extracts the elements into a slice of Extension structs. In other
// words, a generic collection is mapped into a relevant slice.
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
