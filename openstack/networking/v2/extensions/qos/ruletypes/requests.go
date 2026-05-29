package ruletypes

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListRuleTypes returns the list of rule types from the server
func ListRuleTypes(c *gophercloud.ServiceClient) (result pagination.Pager) {
	return pagination.NewPager(c, listRuleTypesURL(c), func(r pagination.PageResult) pagination.Page {
		return ListRuleTypesPage{pagination.SinglePageBase(r)}
	})
}

// GetRuleType retrieves a specific QoS RuleType based on its name.
func GetRuleType(ctx context.Context, c *gophercloud.ServiceClient, name string) (r GetResult) {
	resp, err := c.Get(ctx, getRuleTypeURL(c, name), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
