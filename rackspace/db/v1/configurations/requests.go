package configurations

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/db/v1/instances"
	"github.com/rackspace/gophercloud/pagination"
)

func List(client *gophercloud.ServiceClient) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return ConfigPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, baseURL(client), pageFn)
}

type CreateOptsBuilder interface {
	ToConfigCreateMap() (map[string]interface{}, error)
}

type DatastoreOpts struct {
	Type    string
	Version string
}

func (opts DatastoreOpts) ToMap() (map[string]string, error) {
	datastore := map[string]string{}

	if opts.Type != "" {
		datastore["type"] = opts.Type
	}

	if opts.Version != "" {
		datastore["version"] = opts.Version
	}

	return datastore, nil
}

type CreateOpts struct {
	Datastore   *DatastoreOpts
	Description string
	Name        string
	Values      map[string]interface{}
}

func (opts CreateOpts) ToConfigCreateMap() (map[string]interface{}, error) {
	if opts.Name == "" {
		return nil, errors.New("Name is a required field")
	}
	if len(opts.Values) == 0 {
		return nil, errors.New("Values must be a populated map")
	}

	config := map[string]interface{}{
		"name":   opts.Name,
		"values": opts.Values,
	}

	if opts.Datastore != nil {
		ds, err := opts.Datastore.ToMap()
		if err != nil {
			return config, err
		}
		config["datastore"] = ds
	}

	if opts.Description != "" {
		config["description"] = opts.Description
	}

	return map[string]interface{}{"configuration": config}, nil
}

func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToConfigCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Request("POST", baseURL(client), gophercloud.RequestOpts{
		OkCodes:      []int{201},
		JSONBody:     &reqBody,
		JSONResponse: &res.Body,
	})

	return res
}

func Get(client *gophercloud.ServiceClient, configID string) GetResult {
	var res GetResult

	_, res.Err = client.Request("GET", resourceURL(client, configID), gophercloud.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: &res.Body,
	})

	return res
}

type UpdateOptsBuilder interface {
	ToConfigUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	Datastore   *DatastoreOpts
	Description string
	Name        string
	Values      map[string]interface{}
}

func (opts UpdateOpts) ToConfigUpdateMap() (map[string]interface{}, error) {
	config := map[string]interface{}{}

	if opts.Name != "" {
		config["name"] = opts.Name
	}

	if opts.Description != "" {
		config["description"] = opts.Description
	}

	if opts.Datastore != nil {
		ds, err := opts.Datastore.ToMap()
		if err != nil {
			return config, err
		}
		config["datastore"] = ds
	}

	if len(opts.Values) > 0 {
		config["values"] = opts.Values
	}

	return map[string]interface{}{"configuration": config}, nil
}

func Update(client *gophercloud.ServiceClient, configID string, opts UpdateOptsBuilder) UpdateResult {
	var res UpdateResult

	reqBody, err := opts.ToConfigUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Request("PATCH", resourceURL(client, configID), gophercloud.RequestOpts{
		OkCodes:      []int{200},
		JSONBody:     &reqBody,
		JSONResponse: &res.Body,
	})

	return res
}

func Replace(client *gophercloud.ServiceClient, configID string, opts UpdateOptsBuilder) ReplaceResult {
	var res ReplaceResult

	reqBody, err := opts.ToConfigUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Request("PUT", resourceURL(client, configID), gophercloud.RequestOpts{
		OkCodes:      []int{202},
		JSONBody:     &reqBody,
		JSONResponse: &res.Body,
	})

	return res
}

func Delete(client *gophercloud.ServiceClient, configID string) DeleteResult {
	var res DeleteResult

	_, res.Err = client.Request("DELETE", resourceURL(client, configID), gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	return res
}

func ListInstances(client *gophercloud.ServiceClient, configID string) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return instances.InstancePage{pagination.LinkedPageBase{PageResult: r}}
	}
	return pagination.NewPager(client, instancesURL(client, configID), pageFn)
}

func ListDatastoreParams(client *gophercloud.ServiceClient, datastoreID, versionID string) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return ParamPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, listDSParamsURL(client, datastoreID, versionID), pageFn)
}

func GetDatastoreParam(client *gophercloud.ServiceClient, datastoreID, versionID, paramID string) ParamResult {
	var res ParamResult

	_, res.Err = client.Request("GET", getDSParamURL(client, datastoreID, versionID, paramID), gophercloud.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: &res.Body,
	})

	return res
}

func ListGlobalParams(client *gophercloud.ServiceClient, versionID string) pagination.Pager {
	pageFn := func(r pagination.PageResult) pagination.Page {
		return ParamPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, listGlobalParamsURL(client, versionID), pageFn)
}

func GetGlobalParam(client *gophercloud.ServiceClient, versionID, paramID string) ParamResult {
	var res ParamResult

	_, res.Err = client.Request("GET", getGlobalParamURL(client, versionID, paramID), gophercloud.RequestOpts{
		OkCodes:      []int{200},
		JSONResponse: &res.Body,
	})

	return res
}
