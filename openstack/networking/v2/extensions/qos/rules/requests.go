package rules

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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
func GetBandwidthLimitRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r GetBandwidthLimitRuleResult) {
	resp, err := c.Get(ctx, getBandwidthLimitRuleURL(c, policyID, ruleID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateBandwidthLimitRuleOptsBuilder allows to add additional parameters to the
// CreateBandwidthLimitRule request.
type CreateBandwidthLimitRuleOptsBuilder interface {
	ToBandwidthLimitRuleCreateMap() (map[string]any, error)
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
func (opts CreateBandwidthLimitRuleOpts) ToBandwidthLimitRuleCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bandwidth_limit_rule")
}

// CreateBandwidthLimitRule requests the creation of a new BandwidthLimitRule on the server.
func CreateBandwidthLimitRule(ctx context.Context, client *gophercloud.ServiceClient, policyID string, opts CreateBandwidthLimitRuleOptsBuilder) (r CreateBandwidthLimitRuleResult) {
	b, err := opts.ToBandwidthLimitRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createBandwidthLimitRuleURL(client, policyID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateBandwidthLimitRuleOptsBuilder allows to add additional parameters to the
// UpdateBandwidthLimitRule request.
type UpdateBandwidthLimitRuleOptsBuilder interface {
	ToBandwidthLimitRuleUpdateMap() (map[string]any, error)
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
func (opts UpdateBandwidthLimitRuleOpts) ToBandwidthLimitRuleUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "bandwidth_limit_rule")
}

// UpdateBandwidthLimitRule requests the creation of a new BandwidthLimitRule on the server.
func UpdateBandwidthLimitRule(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string, opts UpdateBandwidthLimitRuleOptsBuilder) (r UpdateBandwidthLimitRuleResult) {
	b, err := opts.ToBandwidthLimitRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateBandwidthLimitRuleURL(client, policyID, ruleID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts policy and rule ID and deletes the BandwidthLimitRule associated with them.
func DeleteBandwidthLimitRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r DeleteBandwidthLimitRuleResult) {
	resp, err := c.Delete(ctx, deleteBandwidthLimitRuleURL(c, policyID, ruleID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DSCPMarkingRulesListOptsBuilder allows extensions to add additional parameters to the
// List request.
type DSCPMarkingRulesListOptsBuilder interface {
	ToDSCPMarkingRulesListQuery() (string, error)
}

// DSCPMarkingRulesListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the DSCPMarking attributes you want to see returned.
// SortKey allows you to sort by a particular DSCPMarkingRule attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type DSCPMarkingRulesListOpts struct {
	ID         string `q:"id"`
	TenantID   string `q:"tenant_id"`
	DSCPMark   int    `q:"dscp_mark"`
	Limit      int    `q:"limit"`
	Marker     string `q:"marker"`
	SortKey    string `q:"sort_key"`
	SortDir    string `q:"sort_dir"`
	Tags       string `q:"tags"`
	TagsAny    string `q:"tags-any"`
	NotTags    string `q:"not-tags"`
	NotTagsAny string `q:"not-tags-any"`
}

// ToDSCPMarkingRulesListQuery formats a ListOpts into a query string.
func (opts DSCPMarkingRulesListOpts) ToDSCPMarkingRulesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDSCPMarkingRules returns a Pager which allows you to iterate over a collection of
// DSCPMarkingRules. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListDSCPMarkingRules(c *gophercloud.ServiceClient, policyID string, opts DSCPMarkingRulesListOptsBuilder) pagination.Pager {
	url := listDSCPMarkingRulesURL(c, policyID)
	if opts != nil {
		query, err := opts.ToDSCPMarkingRulesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return DSCPMarkingRulePage{pagination.LinkedPageBase{PageResult: r}}

	})
}

// GetDSCPMarkingRule retrieves a specific DSCPMarkingRule based on its ID.
func GetDSCPMarkingRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r GetDSCPMarkingRuleResult) {
	resp, err := c.Get(ctx, getDSCPMarkingRuleURL(c, policyID, ruleID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateDSCPMarkingRuleOptsBuilder allows to add additional parameters to the
// CreateDSCPMarkingRule request.
type CreateDSCPMarkingRuleOptsBuilder interface {
	ToDSCPMarkingRuleCreateMap() (map[string]any, error)
}

// CreateDSCPMarkingRuleOpts specifies parameters of a new DSCPMarkingRule.
type CreateDSCPMarkingRuleOpts struct {
	// DSCPMark contains DSCP mark value.
	DSCPMark int `json:"dscp_mark"`
}

// ToDSCPMarkingRuleCreateMap constructs a request body from CreateDSCPMarkingRuleOpts.
func (opts CreateDSCPMarkingRuleOpts) ToDSCPMarkingRuleCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "dscp_marking_rule")
}

// CreateDSCPMarkingRule requests the creation of a new DSCPMarkingRule on the server.
func CreateDSCPMarkingRule(ctx context.Context, client *gophercloud.ServiceClient, policyID string, opts CreateDSCPMarkingRuleOptsBuilder) (r CreateDSCPMarkingRuleResult) {
	b, err := opts.ToDSCPMarkingRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createDSCPMarkingRuleURL(client, policyID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateDSCPMarkingRuleOptsBuilder allows to add additional parameters to the
// UpdateDSCPMarkingRule request.
type UpdateDSCPMarkingRuleOptsBuilder interface {
	ToDSCPMarkingRuleUpdateMap() (map[string]any, error)
}

// UpdateDSCPMarkingRuleOpts specifies parameters for the Update call.
type UpdateDSCPMarkingRuleOpts struct {
	// DSCPMark contains DSCP mark value.
	DSCPMark *int `json:"dscp_mark,omitempty"`
}

// ToDSCPMarkingRuleUpdateMap constructs a request body from UpdateDSCPMarkingRuleOpts.
func (opts UpdateDSCPMarkingRuleOpts) ToDSCPMarkingRuleUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "dscp_marking_rule")
}

// UpdateDSCPMarkingRule requests the creation of a new DSCPMarkingRule on the server.
func UpdateDSCPMarkingRule(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string, opts UpdateDSCPMarkingRuleOptsBuilder) (r UpdateDSCPMarkingRuleResult) {
	b, err := opts.ToDSCPMarkingRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateDSCPMarkingRuleURL(client, policyID, ruleID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteDSCPMarkingRule accepts policy and rule ID and deletes the DSCPMarkingRule associated with them.
func DeleteDSCPMarkingRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r DeleteDSCPMarkingRuleResult) {
	resp, err := c.Delete(ctx, deleteDSCPMarkingRuleURL(c, policyID, ruleID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type MinimumBandwidthRulesListOptsBuilder interface {
	ToMinimumBandwidthRulesListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the MinimumBandwidthRules attributes you want to see returned.
// SortKey allows you to sort by a particular MinimumBandwidthRule attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type MinimumBandwidthRulesListOpts struct {
	ID         string `q:"id"`
	TenantID   string `q:"tenant_id"`
	MinKBps    int    `q:"min_kbps"`
	Direction  string `q:"direction"`
	Limit      int    `q:"limit"`
	Marker     string `q:"marker"`
	SortKey    string `q:"sort_key"`
	SortDir    string `q:"sort_dir"`
	Tags       string `q:"tags"`
	TagsAny    string `q:"tags-any"`
	NotTags    string `q:"not-tags"`
	NotTagsAny string `q:"not-tags-any"`
}

// ToMinimumBandwidthRulesListQuery formats a ListOpts into a query string.
func (opts MinimumBandwidthRulesListOpts) ToMinimumBandwidthRulesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListMinimumBandwidthRules returns a Pager which allows you to iterate over a collection of
// MinimumBandwidthRules. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListMinimumBandwidthRules(c *gophercloud.ServiceClient, policyID string, opts MinimumBandwidthRulesListOptsBuilder) pagination.Pager {
	url := listMinimumBandwidthRulesURL(c, policyID)
	if opts != nil {
		query, err := opts.ToMinimumBandwidthRulesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return MinimumBandwidthRulePage{pagination.LinkedPageBase{PageResult: r}}

	})
}

// GetMinimumBandwidthRule retrieves a specific MinimumBandwidthRule based on its ID.
func GetMinimumBandwidthRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r GetMinimumBandwidthRuleResult) {
	resp, err := c.Get(ctx, getMinimumBandwidthRuleURL(c, policyID, ruleID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateMinimumBandwidthRuleOptsBuilder allows to add additional parameters to the
// CreateMinimumBandwidthRule request.
type CreateMinimumBandwidthRuleOptsBuilder interface {
	ToMinimumBandwidthRuleCreateMap() (map[string]any, error)
}

// CreateMinimumBandwidthRuleOpts specifies parameters of a new MinimumBandwidthRule.
type CreateMinimumBandwidthRuleOpts struct {
	// MaxKBps is a minimum kilobits per second. It's a required parameter.
	MinKBps int `json:"min_kbps"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToMinimumBandwidthRuleCreateMap constructs a request body from CreateMinimumBandwidthRuleOpts.
func (opts CreateMinimumBandwidthRuleOpts) ToMinimumBandwidthRuleCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "minimum_bandwidth_rule")
}

// CreateMinimumBandwidthRule requests the creation of a new MinimumBandwidthRule on the server.
func CreateMinimumBandwidthRule(ctx context.Context, client *gophercloud.ServiceClient, policyID string, opts CreateMinimumBandwidthRuleOptsBuilder) (r CreateMinimumBandwidthRuleResult) {
	b, err := opts.ToMinimumBandwidthRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createMinimumBandwidthRuleURL(client, policyID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMinimumBandwidthRuleOptsBuilder allows to add additional parameters to the
// UpdateMinimumBandwidthRule request.
type UpdateMinimumBandwidthRuleOptsBuilder interface {
	ToMinimumBandwidthRuleUpdateMap() (map[string]any, error)
}

// UpdateMinimumBandwidthRuleOpts specifies parameters for the Update call.
type UpdateMinimumBandwidthRuleOpts struct {
	// MaxKBps is a minimum kilobits per second. It's a required parameter.
	MinKBps *int `json:"min_kbps,omitempty"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToMinimumBandwidthRuleUpdateMap constructs a request body from UpdateMinimumBandwidthRuleOpts.
func (opts UpdateMinimumBandwidthRuleOpts) ToMinimumBandwidthRuleUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "minimum_bandwidth_rule")
}

// UpdateMinimumBandwidthRule requests the creation of a new MinimumBandwidthRule on the server.
func UpdateMinimumBandwidthRule(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string, opts UpdateMinimumBandwidthRuleOptsBuilder) (r UpdateMinimumBandwidthRuleResult) {
	b, err := opts.ToMinimumBandwidthRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateMinimumBandwidthRuleURL(client, policyID, ruleID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMinimumBandwidthRule accepts policy and rule ID and deletes the MinimumBandwidthRule associated with them.
func DeleteMinimumBandwidthRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r DeleteMinimumBandwidthRuleResult) {
	resp, err := c.Delete(ctx, deleteMinimumBandwidthRuleURL(c, policyID, ruleID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// MinimumPacketRateRulesListOptsBuilder allows extensions to add additional parameters to the
// List request.
type MinimumPacketRateRulesListOptsBuilder interface {
	ToMinimumPacketRateRulesListQuery() (string, error)
}

// MinimumPacketRateRulesListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the MinimumPacketRateRules attributes you want to see returned.
// SortKey allows you to sort by a particular MinimumPacketRateRule attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type MinimumPacketRateRulesListOpts struct {
	ID         string `q:"id"`
	TenantID   string `q:"tenant_id"`
	MinKPps    int    `q:"min_kpps"`
	Direction  string `q:"direction"`
	Limit      int    `q:"limit"`
	Marker     string `q:"marker"`
	SortKey    string `q:"sort_key"`
	SortDir    string `q:"sort_dir"`
	Tags       string `q:"tags"`
	TagsAny    string `q:"tags-any"`
	NotTags    string `q:"not-tags"`
	NotTagsAny string `q:"not-tags-any"`
}

// ToMinimumPacketRateRulesListQuery formats a ListOpts into a query string.
func (opts MinimumPacketRateRulesListOpts) ToMinimumPacketRateRulesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListMinimumPacketRateRules returns a Pager which allows you to iterate over a collection of
// MinimumPacketRateRules. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListMinimumPacketRateRules(c *gophercloud.ServiceClient, policyID string, opts MinimumPacketRateRulesListOptsBuilder) pagination.Pager {
	url := listMinimumPacketRateRulesURL(c, policyID)
	if opts != nil {
		query, err := opts.ToMinimumPacketRateRulesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return MinimumPacketRateRulePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetMinimumPacketRateRule retrieves a specific MinimumPacketRateRule based on its ID.
func GetMinimumPacketRateRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r GetMinimumPacketRateRuleResult) {
	resp, err := c.Get(ctx, getMinimumPacketRateRuleURL(c, policyID, ruleID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateMinimumPacketRateRuleOptsBuilder allows to add additional parameters to the
// CreateMinimumPacketRateRule request.
type CreateMinimumPacketRateRuleOptsBuilder interface {
	ToMinimumPacketRateRuleCreateMap() (map[string]any, error)
}

// CreateMinimumPacketRateRuleOpts specifies parameters of a new MinimumPacketRateRule.
type CreateMinimumPacketRateRuleOpts struct {
	// MinKPps is a minimum kilopackets per second. It's a required parameter.
	MinKPps int `json:"min_kpps"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToMinimumPacketRateRuleCreateMap constructs a request body from CreateMinimumPacketRateRuleOpts.
func (opts CreateMinimumPacketRateRuleOpts) ToMinimumPacketRateRuleCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "minimum_packet_rate_rule")
}

// CreateMinimumPacketRateRule requests the creation of a new MinimumPacketRateRule on the server.
func CreateMinimumPacketRateRule(ctx context.Context, client *gophercloud.ServiceClient, policyID string, opts CreateMinimumPacketRateRuleOptsBuilder) (r CreateMinimumPacketRateRuleResult) {
	b, err := opts.ToMinimumPacketRateRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createMinimumPacketRateRuleURL(client, policyID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMinimumPacketRateRuleOptsBuilder allows to add additional parameters to the
// UpdateMinimumPacketRateRule request.
type UpdateMinimumPacketRateRuleOptsBuilder interface {
	ToMinimumPacketRateRuleUpdateMap() (map[string]any, error)
}

// UpdateMinimumPacketRateRuleOpts specifies parameters for the Update call.
type UpdateMinimumPacketRateRuleOpts struct {
	// MinKPps is a minimum kilopackets per second.
	MinKPps *int `json:"min_kpps,omitempty"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToMinimumPacketRateRuleUpdateMap constructs a request body from UpdateMinimumPacketRateRuleOpts.
func (opts UpdateMinimumPacketRateRuleOpts) ToMinimumPacketRateRuleUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "minimum_packet_rate_rule")
}

// UpdateMinimumPacketRateRule requests the update of an existing MinimumPacketRateRule on the server.
func UpdateMinimumPacketRateRule(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string, opts UpdateMinimumPacketRateRuleOptsBuilder) (r UpdateMinimumPacketRateRuleResult) {
	b, err := opts.ToMinimumPacketRateRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updateMinimumPacketRateRuleURL(client, policyID, ruleID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMinimumPacketRateRule accepts policy and rule ID and deletes the MinimumPacketRateRule associated with them.
func DeleteMinimumPacketRateRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r DeleteMinimumPacketRateRuleResult) {
	resp, err := c.Delete(ctx, deleteMinimumPacketRateRuleURL(c, policyID, ruleID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// PacketRateLimitRulesListOptsBuilder allows extensions to add additional parameters to the
// List request.
type PacketRateLimitRulesListOptsBuilder interface {
	ToPacketRateLimitRulesListQuery() (string, error)
}

// PacketRateLimitRulesListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the PacketRateLimitRules attributes you want to see returned.
// SortKey allows you to sort by a particular PacketRateLimitRule attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type PacketRateLimitRulesListOpts struct {
	ID           string `q:"id"`
	TenantID     string `q:"tenant_id"`
	MaxKPps      int    `q:"max_kpps"`
	MaxBurstKPps int    `q:"max_burst_kpps"`
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

// ToPacketRateLimitRulesListQuery formats a ListOpts into a query string.
func (opts PacketRateLimitRulesListOpts) ToPacketRateLimitRulesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListPacketRateLimitRules returns a Pager which allows you to iterate over a collection of
// PacketRateLimitRules. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func ListPacketRateLimitRules(c *gophercloud.ServiceClient, policyID string, opts PacketRateLimitRulesListOptsBuilder) pagination.Pager {
	url := listPacketRateLimitRulesURL(c, policyID)
	if opts != nil {
		query, err := opts.ToPacketRateLimitRulesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PacketRateLimitRulePage{pagination.LinkedPageBase{PageResult: r}}

	})
}

// GetPacketRateLimitRule retrieves a specific PacketRateLimitRule based on its ID.
func GetPacketRateLimitRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r GetPacketRateLimitRuleResult) {
	resp, err := c.Get(ctx, getPacketRateLimitRuleURL(c, policyID, ruleID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreatePacketRateLimitRuleOptsBuilder allows to add additional parameters to the
// CreatePacketRateLimitRule request.
type CreatePacketRateLimitRuleOptsBuilder interface {
	ToPacketRateLimitRuleCreateMap() (map[string]any, error)
}

// CreatePacketRateLimitRuleOpts specifies parameters of a new PacketRateLimitRule.
type CreatePacketRateLimitRuleOpts struct {
	// MaxKPps is a maximum kilopackets per second. It's a required parameter.
	MaxKPps int `json:"max_kpps"`

	// MaxBurstKPps is a maximum burst size in kilopackets.
	MaxBurstKPps int `json:"max_burst_kpps,omitempty"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToPacketRateLimitRuleCreateMap constructs a request body from CreatePacketRateLimitRuleOpts.
func (opts CreatePacketRateLimitRuleOpts) ToPacketRateLimitRuleCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "packet_rate_limit_rule")
}

// CreatePacketRateLimitRule requests the creation of a new PacketRateLimitRule on the server.
func CreatePacketRateLimitRule(ctx context.Context, client *gophercloud.ServiceClient, policyID string, opts CreatePacketRateLimitRuleOptsBuilder) (r CreatePacketRateLimitRuleResult) {
	b, err := opts.ToPacketRateLimitRuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createPacketRateLimitRuleURL(client, policyID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdatePacketRateLimitRuleOptsBuilder allows to add additional parameters to the
// UpdatePacketRateLimitRule request.
type UpdatePacketRateLimitRuleOptsBuilder interface {
	ToPacketRateLimitRuleUpdateMap() (map[string]any, error)
}

// UpdatePacketRateLimitRuleOpts specifies parameters for the Update call.
type UpdatePacketRateLimitRuleOpts struct {
	// MaxKPps is a maximum kilopackets per second.
	MaxKPps *int `json:"max_kpps,omitempty"`

	// MaxBurstKPps is a maximum burst size in kilopackets.
	MaxBurstKPps *int `json:"max_burst_kpps,omitempty"`

	// Direction represents the direction of traffic.
	Direction string `json:"direction,omitempty"`
}

// ToPacketRateLimitRuleUpdateMap constructs a request body from UpdatePacketRateLimitRuleOpts.
func (opts UpdatePacketRateLimitRuleOpts) ToPacketRateLimitRuleUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "packet_rate_limit_rule")
}

// UpdatePacketRateLimitRule requests the update of an existing PacketRateLimitRule on the server.
func UpdatePacketRateLimitRule(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string, opts UpdatePacketRateLimitRuleOptsBuilder) (r UpdatePacketRateLimitRuleResult) {
	b, err := opts.ToPacketRateLimitRuleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, updatePacketRateLimitRuleURL(client, policyID, ruleID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeletePacketRateLimitRule accepts policy and rule ID and deletes the PacketRateLimitRule associated with them.
func DeletePacketRateLimitRule(ctx context.Context, c *gophercloud.ServiceClient, policyID, ruleID string) (r DeletePacketRateLimitRuleResult) {
	resp, err := c.Delete(ctx, deletePacketRateLimitRuleURL(c, policyID, ruleID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
