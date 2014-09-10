package services

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// Service is the result of a list or information query.
type Service struct {
	Description *string `json:"description,omitempty"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
}

// ServiceList is a collection of Services.
type ServiceList struct {
	gophercloud.PaginationLinks `json:"links"`

	client   *gophercloud.ServiceClient
	Services []Service `json:"services"`
}

// Pager indicates that the ServiceList is paginated by next and previous links.
func (list ServiceList) Pager() gophercloud.Pager {
	return gophercloud.NewLinkPager(list)
}

// Concat returns a new collection that's the result of appending a new collection at the end of this one.
func (list ServiceList) Concat(other gophercloud.Collection) gophercloud.Collection {
	return ServiceList{
		client:   list.client,
		Services: append(list.Services, AsServices(other)...),
	}
}

// Service returns the ServiceClient used to acquire this list.
func (list ServiceList) Service() *gophercloud.ServiceClient {
	return list.client
}

// Links accesses pagination information for the current page.
func (list ServiceList) Links() gophercloud.PaginationLinks {
	return list.PaginationLinks
}

// Interpret parses a follow-on JSON response as an additional page.
func (list ServiceList) Interpret(json interface{}) (gophercloud.LinkCollection, error) {
	mapped, ok := json.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected JSON response: %#v", json)
	}

	var result ServiceList
	err := mapstructure.Decode(mapped, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AsServices extracts a slice of Services from a Collection acquired from List.
// It panics if the Collection does not actually contain Services.
func AsServices(results gophercloud.Collection) []Service {
	return results.(*ServiceList).Services
}
