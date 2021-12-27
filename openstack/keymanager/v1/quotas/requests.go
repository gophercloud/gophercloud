package quotas

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Get the effective quotas for the project of the requester. The project id of the requester is derived from the authentication token provided in the X-Auth-Token header.
func Get(client *gophercloud.ServiceClient) (r GetResult) {
	resp, err := client.Get(getURL(client), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToOrderListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Limit is the amount of containers to retrieve.
	Limit int `q:"limit"`

	// Offset is the index within the list to retrieve.
	Offset int `q:"offset"`
}

// ToOrderListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToOrderListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List retrieves a list of orders.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToOrderListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ProjectQuotaPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetProjectQuota returns key manager Quotas for a project.
func GetProjectQuota(client *gophercloud.ServiceClient, projectID string) (r GetProjectResult) {
	resp, err := client.Get(getProjectURL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToQuotaUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update the key manager Quotas.
type UpdateOpts struct {
	// Secrets represents the number of secrets. A "-1" value means no limit.
	Secrets *int `json:"secrets"`

	// Orders represents the number of orders. A "-1" value means no limit.
	Orders *int `json:"orders"`

	// Containers represents the number of containers. A "-1" value means no limit.
	Containers *int `json:"containers"`

	// Consumers represents the number of consumers. A "-1" value means no limit.
	Consumers *int `json:"consumers"`

	// CAS represents the number of cas. A "-1" value means no limit.
	CAS *int `json:"cas"`
}

// ToQuotaUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToQuotaUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "project_quotas")
}

// Update accepts a UpdateOpts struct and updates an existing key manager Quotas using the
// values provided.
func Update(c *gophercloud.ServiceClient, projectID string, opts UpdateOptsBuilder) (r gophercloud.Result) {
	b, err := opts.ToQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(updateProjectURL(c, projectID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete a key manager Quotas for a project.
func Delete(c *gophercloud.ServiceClient, projectID string) (r gophercloud.Result) {
	resp, err := c.Delete(deleteProjectURL(c, projectID), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
