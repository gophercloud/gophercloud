package rules

import "github.com/gophercloud/gophercloud"

const (
	rootPath = "qos/policies"

	bandwidthLimitRulesResourcePath = "bandwidth_limit_rules"
)

func bandwidthLimitRulesRootURL(c *gophercloud.ServiceClient, policyID string) string {
	return c.ServiceURL(rootPath, policyID, bandwidthLimitRulesResourcePath)
}

func bandwidthLimitRulesResourceURL(c *gophercloud.ServiceClient, policyID, ruleID string) string {
	return c.ServiceURL(rootPath, policyID, bandwidthLimitRulesResourcePath, ruleID)
}

func listBandwidthLimitRulesURL(c *gophercloud.ServiceClient, policyID string) string {
	return bandwidthLimitRulesRootURL(c, policyID)
}

func getBandwidthLimitRuleURL(c *gophercloud.ServiceClient, policyID, ruleID string) string {
	return bandwidthLimitRulesResourceURL(c, policyID, ruleID)
}

func createBandwidthLimitRuleURL(c *gophercloud.ServiceClient, policyID string) string {
	return bandwidthLimitRulesRootURL(c, policyID)
}
