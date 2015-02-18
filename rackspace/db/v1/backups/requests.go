package backups

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	ToBackupCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name string

	InstanceID string

	Description string
}

func (opts CreateOpts) ToBackupCreateMap() (map[string]interface{}, error) {
	if opts.Name == "" {
		return nil, errors.New("Name is a required field")
	}
	if opts.InstanceID == "" {
		return nil, errors.New("InstanceID is a required field")
	}

	backup := map[string]interface{}{
		"name":     opts.Name,
		"instance": opts.InstanceID,
	}

	if opts.Description != "" {
		backup["description"] = opts.Description
	}

	return map[string]interface{}{"backup": backup}, nil
}

func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToBackupCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Request("POST", baseURL(client), gophercloud.RequestOpts{
		JSONBody:     &reqBody,
		JSONResponse: &res.Body,
		OkCodes:      []int{202},
	})

	return res
}

type ListOptsBuilder interface {
	ToBackupListQuery() (string, error)
}

type ListOpts struct {
	Datastore string `q:"datastore"`
}

func (opts ListOpts) ToBackupListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)

	if opts != nil {
		query, err := opts.ToBackupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	pageFn := func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, pageFn)
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult

	_, res.Err = client.Request("GET", resourceURL(client, id), gophercloud.RequestOpts{
		JSONResponse: &res.Body,
		OkCodes:      []int{200},
	})

	return res
}

func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult

	_, res.Err = client.Request("DELETE", resourceURL(client, id), gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	return res
}
