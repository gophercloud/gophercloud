package apiversions

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type APIVersion struct {
	Status string `mapstructure:"status" json:"status"`
	ID     string `mapstructure:"id" json:"id"`
}

func ExtractAPIVersions(page gophercloud.Page) ([]APIVersion, error) {
	var resp struct {
		Versions []APIVersion `mapstructure:"versions"`
	}

	err := mapstructure.Decode(page.(gophercloud.LinkedPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Versions, nil
}

type APIVersionResource struct {
	Name       string `mapstructure:"name" json:"name"`
	Collection string `mapstructure:"collection" json:"collection"`
}

func ExtractVersionResources(page gophercloud.Page) ([]APIVersionResource, error) {
	var resp struct {
		APIVersionResources []APIVersionResource `mapstructure:"resources"`
	}

	err := mapstructure.Decode(page.(gophercloud.LinkedPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.APIVersionResources, nil
}
