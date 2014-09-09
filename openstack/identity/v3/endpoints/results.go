package endpoints

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// Endpoint describes the entry point for another service's API.
type Endpoint struct {
	ID        string                `json:"id"`
	Interface gophercloud.Interface `json:"interface"`
	Name      string                `json:"name"`
	Region    string                `json:"region"`
	ServiceID string                `json:"service_id"`
	URL       string                `json:"url"`
}

// EndpointList contains a page of Endpoint results.
type EndpointList struct {
	gophercloud.PaginationLinks `json:"links"`

	client    *gophercloud.ServiceClient
	Endpoints []Endpoint `json:"endpoints"`
}

// Pager marks EndpointList as paged by links.
func (list EndpointList) Pager() gophercloud.Pager {
	return gophercloud.NewLinkPager(list)
}

// Concat adds the contents of another Collection to this one.
func (list EndpointList) Concat(other gophercloud.Collection) gophercloud.Collection {
	return EndpointList{
		client:    list.client,
		Endpoints: append(list.Endpoints, AsEndpoints(other)...),
	}
}

// Service returns the ServiceClient used to acquire this list.
func (list EndpointList) Service() *gophercloud.ServiceClient {
	return list.client
}

// Links accesses pagination information for the current page.
func (list EndpointList) Links() gophercloud.PaginationLinks {
	return list.PaginationLinks
}

// Interpret parses a follow-on JSON response as an additional page.
func (list EndpointList) Interpret(json interface{}) (gophercloud.LinkCollection, error) {
	mapped, ok := json.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected JSON response: %#v", json)
	}

	var result EndpointList
	err := mapstructure.Decode(mapped, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AsEndpoints extracts an Endpoint slice from a Collection.
// Panics if `list` was not returned from a List call.
func AsEndpoints(list gophercloud.Collection) []Endpoint {
	return list.(*EndpointList).Endpoints
}
