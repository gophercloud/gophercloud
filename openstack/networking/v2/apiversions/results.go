package apiversions

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

type APIVersion struct {
	Status string `mapstructure:"status" json:"status"`
	ID     string `mapstructure:"id" json:"id"`
}

type APIVersionPage struct {
	pagination.SinglePageBase
}

func (r APIVersionPage) IsEmpty() (bool, error) {
	is, err := ExtractAPIVersions(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractAPIVersions(page pagination.Page) ([]APIVersion, error) {
	var resp struct {
		Versions []APIVersion `mapstructure:"versions"`
	}

	err := mapstructure.Decode(page.(APIVersionPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Versions, nil
}

type APIVersionResource struct {
	Name       string `mapstructure:"name" json:"name"`
	Collection string `mapstructure:"collection" json:"collection"`
}

type APIVersionResourcePage struct {
	pagination.SinglePageBase
}

func (r APIVersionResourcePage) IsEmpty() (bool, error) {
	is, err := ExtractVersionResources(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractVersionResources(page pagination.Page) ([]APIVersionResource, error) {
	var resp struct {
		APIVersionResources []APIVersionResource `mapstructure:"resources"`
	}

	err := mapstructure.Decode(page.(APIVersionResourcePage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.APIVersionResources, nil
}
