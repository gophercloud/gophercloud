package instances

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
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

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	gophercloud.Result
}

func (r CreateResult) Extract() (*Instance, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Instance Instance `mapstructure:"instance"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.Instance, err
}
