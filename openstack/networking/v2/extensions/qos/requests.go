package qos

import "github.com/gophercloud/gophercloud"

// ListRuleTypes returns the list of rule types from the server
func ListRuleTypes(c *gophercloud.ServiceClient) (result ListRuleTypesResult) {
	_, result.Err = c.Get(listRuleTypesURL(c), &result.Body, nil)
	return
}