package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List makes a request against the API to list services.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.SinglePageBase(r)}
	})
}

// UpdateOpts specifies the base attributes that may be updated on a service.
type UpdateOpts struct {
	// Status represents the new service status. One of enabled or disabled.
	Status string `json:"status,omitempty"`

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
	_, r.Err = client.Put(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
