package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request.
type ListOptsBuilder interface {
	ToServicesListQuery() (string, error)
}

// ListOpts represents options to list services.
type ListOpts struct {
	Binary string `q:"binary"`
	Host   string `q:"host"`
}

// ToServicesListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServicesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list services.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToServicesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.SinglePageBase(r)}
	})
}

type ServiceStatus string

const (
	// ServiceEnabled is used to mark a service as being enabled.
	ServiceEnabled ServiceStatus = "enabled"

	// ServiceDisabled is used to mark a service as being disabled.
	ServiceDisabled ServiceStatus = "disabled"
)

// UpdateOpts specifies the base attributes that may be updated on a service.
type UpdateOpts struct {
	// Status represents the new service status. One of enabled or disabled.
	Status ServiceStatus `json:"status,omitempty"`

	// DisabledReason represents the reason for disabling a service.
	DisabledReason string `json:"disabled_reason,omitempty"`

	// ForcedDown is a manual override to tell nova that the service in question
	// has been fenced manually by the operations team.
	ForcedDown bool `json:"forced_down,omitempty"`
}

// ToServiceUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToServiceUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update requests that various attributes of the indicated service be changed.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	b, err := opts.ToServiceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will delete the existing service with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(updateURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
