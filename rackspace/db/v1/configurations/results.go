package configurations

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type Config struct {
	Created              string
	Updated              string
	DatastoreName        string `mapstructure:"datastore_name"`
	DatastoreVersionID   string `mapstructure:"datastore_version_id"`
	DatastoreVersionName string `mapstructure:"datastore_version_name"`
	Description          string
	ID                   string
	Name                 string
	Values               map[string]interface{}
}

type ConfigPage struct {
	pagination.SinglePageBase
}

func (r ConfigPage) IsEmpty() (bool, error) {
	is, err := ExtractConfigs(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractConfigs(page pagination.Page) ([]Config, error) {
	casted := page.(ConfigPage).Body

	var resp struct {
		Configs []Config `mapstructure:"configurations" json:"configurations"`
	}

	err := mapstructure.Decode(casted, &resp)
	return resp.Configs, err
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Config, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Config Config `mapstructure:"configuration"`
	}

	err := mapstructure.Decode(r.Body, &response)
	return &response.Config, err
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	gophercloud.ErrResult
}

type ReplaceResult struct {
	gophercloud.ErrResult
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type Param struct {
	Max             int
	Min             int
	Name            string
	RestartRequired bool `mapstructure:"restart_required" json:"restart_required"`
	Type            string
}

type ParamPage struct {
	pagination.SinglePageBase
}

func (r ParamPage) IsEmpty() (bool, error) {
	is, err := ExtractParams(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractParams(page pagination.Page) ([]Param, error) {
	casted := page.(ParamPage).Body

	var resp struct {
		Params []Param `mapstructure:"configuration-parameters" json:"configuration-parameters"`
	}

	err := mapstructure.Decode(casted, &resp)
	return resp.Params, err
}

type ParamResult struct {
	gophercloud.Result
}

func (r ParamResult) Extract() (*Param, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var param Param

	err := mapstructure.Decode(r.Body, &param)
	return &param, err
}
