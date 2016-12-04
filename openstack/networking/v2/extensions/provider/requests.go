package provider

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type ListOptsBuilder interface {
	ToProviderNetworkListQuery() (string, error)
}

type ListOpts struct {
	Status          string `q:"status"`
	Name            string `q:"name"`
	AdminStateUp    *bool  `q:"admin_state_up"`
	TenantID        string `q:"tenant_id"`
	Shared          *bool  `q:"shared"`
	ID              string `q:"id"`
	NetworkType     string `q:"provider:network_type"`
	PhysicalNetwork string `q:"provider:physical_network"`
	SegmentationID  string `q:"provider:segmentation_id"`
	Marker          string `q:"marker"`
	Limit           int    `q:"limit"`
	SortKey         string `q:"sort_key"`
	SortDir         string `q:"sort_dir"`
}

func (opts ListOpts) ToProviderNetworkListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToProviderNetworkListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NetworkPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

type CreateOptsBuilder interface {
	ToProviderNetworkCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	AdminStateUp    *bool  `json:"admin_state_up,omitempty"`
	Name            string `json:"name,omitempty"`
	Shared          *bool  `json:"shared,omitempty"`
	TenantID        string `json:"tenant_id,omitempty"`
	NetworkType     string `json:"provider:network_type,omitempty"`
	PhysicalNetwork string `json:"provider:physical_network,omitempty"`
	SegmentationID  string `json:"provider:segmentation_id,omitempty"`
}

func (opts CreateOpts) ToProviderNetworkCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "network")
}

func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToProviderNetworkCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

type UpdateOptsBuilder interface {
	ToProviderNetworkUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	AdminStateUp    *bool  `json:"admin_state_up,omitempty"`
	Name            string `json:"name,omitempty"`
	Shared          *bool  `json:"shared,omitempty"`
	NetworkType     string `json:"provider:network_type"`
	PhysicalNetwork string `json:"provider:physical_network"`
	SegmentationID  string `json:"provider:segmentation_id"`
}

func (opts UpdateOpts) ToProviderNetworkUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "network")
}

func Update(c *gophercloud.ServiceClient, networkID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToProviderNetworkUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, networkID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

func Delete(c *gophercloud.ServiceClient, networkID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, networkID), nil)
	return
}

func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""
	pages, err := List(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractList(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{Name: name, ResourceType: "network"}
	case 1:
		return id, nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "network"}
	}
}
