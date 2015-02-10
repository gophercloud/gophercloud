package instances

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/db/v1/instances"
)

type Datastore struct {
	Type    string
	Version string
}

type Instance struct {
	Created   string //time.Time
	Updated   string //time.Time
	Datastore Datastore
	Flavor    os.Flavor
	Hostname  string
	ID        string
	Links     []gophercloud.Link
	Name      string
	Status    string
	Volume    os.Volume
}

func commonExtract(err error, body interface{}) (*Instance, error) {
	if err != nil {
		return nil, err
	}

	var response struct {
		Instance Instance `mapstructure:"instance"`
	}

	err = mapstructure.Decode(body, &response)
	return &response.Instance, err
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	os.CreateResult
}

func (r CreateResult) Extract() (*Instance, error) {
	return commonExtract(r.Err, r.Body)
}

type GetResult struct {
	os.GetResult
}

func (r GetResult) Extract() (*Instance, error) {
	return commonExtract(r.Err, r.Body)
}

type ConfigResult struct {
	gophercloud.Result
}

func (r ConfigResult) Extract() (map[string]string, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Instance struct {
			Config map[string]string `mapstructure:"configuration"`
		} `mapstructure:"instance"`
	}

	err := mapstructure.Decode(r.Body, &response)
	return response.Instance.Config, err
}

type UpdateResult struct {
	gophercloud.ErrResult
}
