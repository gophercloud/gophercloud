package networks

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// A Network represents a a virtual layer-2 broadcast domain.
type Network struct {
	// Id is the unique identifier for the network.
	Id string `json:"id"`
	// Name is the (not necessarily unique) human-readable identifier for the network.
	Name string `json:"name"`
	// AdminStateUp is administrative state of the network. If false, network is down.
	AdminStateUp bool `json:"admin_state_up"`
	// Status indicates if the network is operational. Possible values: active, down, build, error.
	Status string `json:"status"`
	// Subnets are IP address blocks that can be used to assign IP addresses to virtual instances.
	Subnets []string `json:"subnets"`
	// Shared indicates whether the network can be accessed by any tenant or not.
	Shared bool `json:"shared"`
	// TenantId is the owner of the network. Admins may specify TenantId other than their own.
	TenantId string `json:"tenant_id"`
	// RouterExternal indicates if the network is connected to an external router.
	RouterExternal bool `json:"router:external"`
	// ProviderPhysicalNetwork is the name of the provider physical network.
	ProviderPhysicalNetwork string `json:"provider:physical_network"`
	// ProviderNetworkType is the type of provider network (eg "vlan").
	ProviderNetworkType string `json:"provider:network_type"`
	// ProviderSegmentationId is the provider network identifier (such as the vlan id).
	ProviderSegmentationId string `json:"provider:segmentation_id"`
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
