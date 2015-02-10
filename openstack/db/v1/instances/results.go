package instances

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type Flavor struct {
	ID    string
	Links []gophercloud.Link
}

type Volume struct {
	Size int
}

type Instance struct {
	Created  string //time.Time
	Updated  string //time.Time
	Flavor   Flavor
	Hostname string
	ID       string
	Links    []gophercloud.Link
	Name     string
	Status   string
	Volume   Volume
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

func (r commonResult) Extract() (*Instance, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Instance Instance `mapstructure:"instance"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.Instance, err
}

type InstancePage struct {
	pagination.LinkedPageBase
}

func (page InstancePage) IsEmpty() (bool, error) {
	instances, err := ExtractInstances(page)
	if err != nil {
		return true, err
	}
	return len(instances) == 0, nil
}

func (page InstancePage) NextPageURL() (string, error) {
	type resp struct {
		Links []gophercloud.Link `mapstructure:"instances_links"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(r.Links)
}

func ExtractInstances(page pagination.Page) ([]Instance, error) {
	casted := page.(InstancePage).Body

	var response struct {
		Instances []Instance `mapstructure:"instances"`
	}

	err := mapstructure.Decode(casted, &response)

	return response.Instances, err
}
