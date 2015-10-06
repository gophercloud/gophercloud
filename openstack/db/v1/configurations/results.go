package configurations

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Config represents a configuration group API resource.
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

// ConfigPage contains a page of Config resources in a paginated collection.
type ConfigPage struct {
	pagination.SinglePageBase
}

// IsEmpty indicates whether a ConfigPage is empty.
func (r ConfigPage) IsEmpty() (bool, error) {
	is, err := ExtractConfigs(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

// ExtractConfigs will retrieve a slice of Config structs from a page.
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

// Extract will retrieve a Config resource from an operation result.
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

// GetResult represents the result of a Get operation.
type GetResult struct {
	commonResult
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of an Update operation.
type UpdateResult struct {
	gophercloud.ErrResult
}

// ReplaceResult represents the result of a Replace operation.
type ReplaceResult struct {
	gophercloud.ErrResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Param represents a configuration parameter API resource.
type Param struct {
	Max             int
	Min             int
	Name            string
	RestartRequired bool `mapstructure:"restart_required" json:"restart_required"`
	Type            string
}

// ParamPage contains a page of Param resources in a paginated collection.
type ParamPage struct {
	pagination.SinglePageBase
}

// IsEmpty indicates whether a ParamPage is empty.
func (r ParamPage) IsEmpty() (bool, error) {
	is, err := ExtractParams(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

// ExtractParams will retrieve a slice of Param structs from a page.
func ExtractParams(page pagination.Page) ([]Param, error) {
	casted := page.(ParamPage).Body

	var resp struct {
		Params []Param `mapstructure:"configuration-parameters" json:"configuration-parameters"`
	}

	err := mapstructure.Decode(casted, &resp)
	return resp.Params, err
}

// ParamResult represents the result of an operation which retrieves details
// about a particular configuration param.
type ParamResult struct {
	gophercloud.Result
}

// Extract will retrieve a param from an operation result.
func (r ParamResult) Extract() (*Param, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var param Param

	err := mapstructure.Decode(r.Body, &param)
	return &param, err
}
