package services

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"
)

type response struct {
	Service Service `json:"service"`
}

// Create adds a new service of the requested type to the catalog.
func Create(client *gophercloud.ServiceClient, serviceType string) CreateResult {
	type request struct {
		Type string `json:"type"`
	}

	req := request{Type: serviceType}

	var result CreateResult
	_, result.Err = perigee.Request("POST", listURL(client), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &req,
		Results:     &result.Body,
		OkCodes:     []int{201},
	})
	return result
}

// ListOpts allows you to query the List method.
type ListOpts struct {
	ServiceType string
	PerPage     int
	Page        int
}

// List enumerates the services available to a specific user.
func List(client *gophercloud.ServiceClient, opts ListOpts) pagination.Pager {
	q := make(map[string]string)
	if opts.ServiceType != "" {
		q["type"] = opts.ServiceType
	}
	if opts.Page != 0 {
		q["page"] = strconv.Itoa(opts.Page)
	}
	if opts.PerPage != 0 {
		q["perPage"] = strconv.Itoa(opts.PerPage)
	}
	u := listURL(client) + utils.BuildQuery(q)

	createPage := func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, u, createPage)
}

// Get returns additional information about a service, given its ID.
func Get(client *gophercloud.ServiceClient, serviceID string) GetResult {
	var result GetResult
	_, result.Err = perigee.Request("GET", serviceURL(client, serviceID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &result.Body,
		OkCodes:     []int{200},
	})
	return result
}

// Update changes the service type of an existing service.
func Update(client *gophercloud.ServiceClient, serviceID string, serviceType string) UpdateResult {
	type request struct {
		Type string `json:"type"`
	}

	req := request{Type: serviceType}

	var result UpdateResult
	_, result.Err = perigee.Request("PATCH", serviceURL(client, serviceID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &req,
		Results:     &result.Body,
		OkCodes:     []int{200},
	})
	return result
}

// Delete removes an existing service.
// It either deletes all associated endpoints, or fails until all endpoints are deleted.
func Delete(client *gophercloud.ServiceClient, serviceID string) error {
	_, err := perigee.Request("DELETE", serviceURL(client, serviceID), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}
