package services

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
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
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     &req,
		Results:     &resp,
		OkCodes:     []int{201},
	})
	if err != nil {
		return nil, err
	}

	return &resp.Service, nil
}
