package datastores

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type Version struct {
	ID    string
	Links []gophercloud.Link
	Name  string
}

type Datastore struct {
	DefaultVersion string `json:"default_version" mapstructure:"default_version"`
	ID             string
	Links          []gophercloud.Link
	Name           string
	Versions       []Version
}

type DatastorePartial struct {
	Version   string
	Type      string
	VersionID string `json:"version_id" mapstructure:"version_id"`
}

type GetResult struct {
	gophercloud.Result
}

type GetVersionResult struct {
	gophercloud.Result
}

type DatastorePage struct {
	pagination.SinglePageBase
}

func (r DatastorePage) IsEmpty() (bool, error) {
	is, err := ExtractDatastores(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractDatastores(page pagination.Page) ([]Datastore, error) {
	casted := page.(DatastorePage).Body

	var resp struct {
		Datastores []Datastore `mapstructure:"datastores" json:"datastores"`
	}

	err := mapstructure.Decode(casted, &resp)
	return resp.Datastores, err
}

func (r GetResult) Extract() (*Datastore, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Datastore Datastore `mapstructure:"datastore"`
	}

	err := mapstructure.Decode(r.Body, &response)
	return &response.Datastore, err
}

type VersionPage struct {
	pagination.SinglePageBase
}

func (r VersionPage) IsEmpty() (bool, error) {
	is, err := ExtractVersions(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractVersions(page pagination.Page) ([]Version, error) {
	casted := page.(VersionPage).Body

	var resp struct {
		Versions []Version `mapstructure:"versions" json:"versions"`
	}

	err := mapstructure.Decode(casted, &resp)
	return resp.Versions, err
}

func (r GetVersionResult) Extract() (*Version, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Version Version `mapstructure:"version"`
	}

	err := mapstructure.Decode(r.Body, &response)
	return &response.Version, err
}
