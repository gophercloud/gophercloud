package quotas

import (
	"github.com/gophercloud/gophercloud"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToQuotasListQuery() (string, error)
}

type ListOpts struct {
}

// ToZoneListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToQuotasListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get returns information about the quota, given its ID.
func Get(client *gophercloud.ServiceClient, projectID string) (r Result) {
	resp, err := client.Get(URL(client, projectID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToQuotaUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update the DNS Quotas.
type UpdateOpts struct {
	APIExporterSize  *int `json:"api_export_size,omitempty"`
	RecordsetRecords *int `json:"recordset_records,omitempty"`
	ZoneRecords      *int `json:"zone_records,omitempty"`
	ZoneRecordsets   *int `json:"zone_recordsets,omitempty"`
	Zones            *int `json:"zones,omitempty"`
}

// ToQuotaUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToQuotaUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "quota")
}

// Update accepts a UpdateOpts struct and updates an existing DNS Quotas using the
// values provided.
func Update(c *gophercloud.ServiceClient, projectID string, opts UpdateOptsBuilder) (r Result) {
	b, err := opts.ToQuotaUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Patch(URL(c, projectID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
