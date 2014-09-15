package networks

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

type NetworkProvider struct {
	ProviderSegmentationID  int    `json:"provider:segmentation_id"`
	ProviderPhysicalNetwork string `json:"provider:physical_network"`
	ProviderNetworkType     string `json:"provider:network_type"`
}

type Network struct {
	Status       string        `json:"status"`
	Subnets      []interface{} `json:"subnets"`
	Name         string        `json:"name"`
	AdminStateUp bool          `json:"admin_state_up"`
	TenantID     string        `json:"tenant_id"`
	Shared       bool          `json:"shared"`
	ID           string        `json:"id"`
}

type NetworkResult struct {
	Network
	NetworkProvider
	RouterExternal bool `json:"router:external"`
}

type NetworkCreateResult struct {
	Network
	Segments            []NetworkProvider `json:"segments"`
	PortSecurityEnabled bool              `json:"port_security_enabled"`
}

type APIVersion struct {
	Status string
	ID     string
}

type APIVersionsList struct {
	gophercloud.PaginationLinks `json:"links"`
	Client                      *gophercloud.ServiceClient
	APIVersions                 []APIVersion `json:"versions"`
}

func (list APIVersionsList) Pager() gophercloud.Pager {
	return gophercloud.NewLinkPager(list)
}

func (list APIVersionsList) Concat(other gophercloud.Collection) gophercloud.Collection {
	return APIVersionsList{
		Client:      list.Client,
		APIVersions: append(list.APIVersions, ToAPIVersions(other)...),
	}
}

func (list APIVersionsList) Service() *gophercloud.ServiceClient {
	return list.Client
}

func (list APIVersionsList) Links() gophercloud.PaginationLinks {
	return list.PaginationLinks
}

func (list APIVersionsList) Interpret(json interface{}) (gophercloud.LinkCollection, error) {
	mapped, ok := json.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected JSON response: %#v", json)
	}

	var result APIVersionsList
	err := mapstructure.Decode(mapped, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToAPIVersions(results gophercloud.Collection) []APIVersion {
	return results.(*APIVersionsList).APIVersions
}

type APIResource struct {
	Name       string
	Collection string
}

type APIInfoList struct {
	gophercloud.PaginationLinks `json:"links"`
	Client                      *gophercloud.ServiceClient
	APIResources                []APIResource `json:"resources"`
}

func (list APIInfoList) Pager() gophercloud.Pager {
	return gophercloud.NewLinkPager(list)
}

func (list APIInfoList) Concat(other gophercloud.Collection) gophercloud.Collection {
	return APIInfoList{
		Client:       list.Client,
		APIResources: append(list.APIResources, ToAPIResource(other)...),
	}
}

func (list APIInfoList) Service() *gophercloud.ServiceClient {
	return list.Client
}

func (list APIInfoList) Links() gophercloud.PaginationLinks {
	return list.PaginationLinks
}

func (list APIInfoList) Interpret(json interface{}) (gophercloud.LinkCollection, error) {
	mapped, ok := json.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected JSON response: %#v", json)
	}

	var result APIInfoList
	err := mapstructure.Decode(mapped, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToAPIResource(results gophercloud.Collection) []APIResource {
	return results.(*APIInfoList).APIResources
}

type Extension struct {
	Updated     string        `json:"updated"`
	Name        string        `json:"name"`
	Links       []interface{} `json:"links"`
	Namespace   string        `json:"namespace"`
	Alias       string        `json:"alias"`
	Description string        `json:"description"`
}
