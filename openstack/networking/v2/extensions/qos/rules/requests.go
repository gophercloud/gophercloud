package rules

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type BandwidthLimitRulesListOptsBuilder interface {
	ToBandwidthLimitRulesListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the BandwidthLimitRules attributes you want to see returned.
// SortKey allows you to sort by a particular BandwidthLimitRule attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type BandwidthLimitRulesListOpts struct {
	ID           string `q:"id"`
	TenantID     string `q:"tenant_id"`
	MaxKBps      int    `q:"max_kbps"`
	MaxBurstKBps int    `q:"max_burst_kbps"`
	Direction    string `q:"direction"`
	Limit        int    `q:"limit"`
	Marker       string `q:"marker"`
	SortKey      string `q:"sort_key"`
	SortDir      string `q:"sort_dir"`
	Tags         string `q:"tags"`
	TagsAny      string `q:"tags-any"`
	NotTags      string `q:"not-tags"`
	NotTagsAny   string `q:"not-tags-any"`
}

// ToBandwidthLimitRulesListQuery formats a ListOpts into a query string.
func (opts BandwidthLimitRulesListOpts) ToBandwidthLimitRulesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListBandwidthLimitRules returns a Pager which allows you to iterate over a collection of
// BandwidthLimitRules. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListBandwidthLimitRules(c *gophercloud.ServiceClient, policyID string, opts BandwidthLimitRulesListOptsBuilder) pagination.Pager {
	url := listBandwidthLimitRulesURL(c, policyID)
	if opts != nil {
		query, err := opts.ToBandwidthLimitRulesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BandwidthLimitRulePage{pagination.LinkedPageBase{PageResult: r}}

	})
}

// GetBandwidthLimitRule retrieves a specific BandwidthLimitRule based on its ID.
func GetBandwidthLimitRule(c *gophercloud.ServiceClient, policyID, ruleID string) (r GetBandwidthLimitRuleResult) {
	_, r.Err = c.Get(getBandwidthLimitRuleURL(c, policyID, ruleID), &r.Body, nil)
	return
}

// CreateBandwidthLimitRuleOptsBuilder allows to add additional parameters to the
// CreateBandwidthLimitRule request.
type CreateBandwidthLimitRuleOptsBuilder interface {
	ToBandwidthLimitRuleCreateMap() (map[string]interface{}, error)
}

// CreateBandwidthLimitRuleOpts specifies parameters of a new BandwidthLimitRule.
type CreateBandwidthLimitRuleOpts struct {
	// MaxKBps is a maximum kilobits per second. It's a required parameter.
	MaxKBps int `json:"max_kbps"`

	// MaxBurstKBps is a maximum burst size in kilobits.
	MaxBurstKBps int `json:"max_burst_kbps,omitempty"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToBandwidthLimitRuleCreateMap constructs a request body from CreateBandwidthLimitRuleOpts.
func (opts CreateBandwidthLimitRuleOpts) ToBandwidthLimitRuleCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "bandwidth_limit_rule")
}

// CreateBandwidthLimitRule requests the creation of a new BandwidthLimitRule on the server.
func CreateBandwidthLimitRule(client *gophercloud.ServiceClient, policyID string, opts CreateBandwidthLimitRuleOptsBuilder) (r CreateBandwidthLimitRuleResult) {
	b, err := opts.ToBandwidthLimitRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createBandwidthLimitRuleURL(client, policyID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// UpdateBandwidthLimitRuleOptsBuilder allows to add additional parameters to the
// UpdateBandwidthLimitRule request.
type UpdateBandwidthLimitRuleOptsBuilder interface {
	ToBandwidthLimitRuleUpdateMap() (map[string]interface{}, error)
}

// UpdateBandwidthLimitRuleOpts specifies parameters for the Update call.
type UpdateBandwidthLimitRuleOpts struct {
	// MaxKBps is a maximum kilobits per second.
	MaxKBps *int `json:"max_kbps,omitempty"`

	// MaxBurstKBps is a maximum burst size in kilobits.
	MaxBurstKBps *int `json:"max_burst_kbps,omitempty"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToBandwidthLimitRuleUpdateMap constructs a request body from UpdateBandwidthLimitRuleOpts.
func (opts UpdateBandwidthLimitRuleOpts) ToBandwidthLimitRuleUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "bandwidth_limit_rule")
}

// UpdateBandwidthLimitRule requests the creation of a new BandwidthLimitRule on the server.
func UpdateBandwidthLimitRule(client *gophercloud.ServiceClient, policyID, ruleID string, opts UpdateBandwidthLimitRuleOptsBuilder) (r UpdateBandwidthLimitRuleResult) {
	b, err := opts.ToBandwidthLimitRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateBandwidthLimitRuleURL(client, policyID, ruleID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

// Delete accepts policy and rule ID and deletes the BandwidthLimitRule associated with them.
func DeleteBandwidthLimitRule(c *gophercloud.ServiceClient, policyID, ruleID string) (r DeleteBandwidthLimitRuleResult) {
	_, r.Err = c.Delete(deleteBandwidthLimitRuleURL(c, policyID, ruleID), nil)
	return
}
