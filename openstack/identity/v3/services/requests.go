package services

import (
	"strconv"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
)

// Create adds a new service of the requested type to the catalog.
func Create(client *gophercloud.ServiceClient, serviceType string) (*ServiceResult, error) {
	type request struct {
		Type string `json:"type"`
	}

	type response struct {
		Service ServiceResult `json:"service"`
	}

	req := request{Type: serviceType}
	var resp response

	_, err := perigee.Request("POST", getListURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		ReqBody:     &req,
		Results:     &resp,
		OkCodes:     []int{201},
	})
	if err != nil {
		return nil, err
	}

	return &resp.Service, nil
}

// ListOpts allows you to query the List method.
type ListOpts struct {
	ServiceType string
	PerPage     int
	Page        int
}

// List enumerates the services available to a specific user.
func List(client *gophercloud.ServiceClient, opts ListOpts) (*ServiceListResult, error) {
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
	u := getListURL(client) + utils.BuildQuery(q)

	var resp ServiceListResult
	_, err := perigee.Request("GET", u, perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		Results:     &resp,
		OkCodes:     []int{200},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
