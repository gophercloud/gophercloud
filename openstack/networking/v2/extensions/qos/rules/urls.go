package rules

import "github.com/gophercloud/gophercloud"

const (
	rootPath = "qos/policies"

	bandwidthLimitRulesResourcePath = "bandwidth_limit_rules"
)

func bandwidthLimitRulesRootURL(c *gophercloud.ServiceClient, policyID string) string {
	return c.ServiceURL(rootPath, policyID, bandwidthLimitRulesResourcePath)
}

func listBandwidthLimitRulesURL(c *gophercloud.ServiceClient, policyID string) string {
	return bandwidthLimitRulesRootURL(c, policyID)
}
