package rules

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath = "qos/policies"

	bandwidthLimitRulesResourcePath   = "bandwidth_limit_rules"
	dscpMarkingRulesResourcePath      = "dscp_marking_rules"
	minimumBandwidthRulesResourcePath = "minimum_bandwidth_rules"
)

func bandwidthLimitRulesRootURL(c gophercloud.Client, policyID string) string {
	return c.ServiceURL(rootPath, policyID, bandwidthLimitRulesResourcePath)
}

func bandwidthLimitRulesResourceURL(c gophercloud.Client, policyID, ruleID string) string {
	return c.ServiceURL(rootPath, policyID, bandwidthLimitRulesResourcePath, ruleID)
}

func listBandwidthLimitRulesURL(c gophercloud.Client, policyID string) string {
	return bandwidthLimitRulesRootURL(c, policyID)
}

func getBandwidthLimitRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return bandwidthLimitRulesResourceURL(c, policyID, ruleID)
}

func createBandwidthLimitRuleURL(c gophercloud.Client, policyID string) string {
	return bandwidthLimitRulesRootURL(c, policyID)
}

func updateBandwidthLimitRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return bandwidthLimitRulesResourceURL(c, policyID, ruleID)
}

func deleteBandwidthLimitRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return bandwidthLimitRulesResourceURL(c, policyID, ruleID)
}

func dscpMarkingRulesRootURL(c gophercloud.Client, policyID string) string {
	return c.ServiceURL(rootPath, policyID, dscpMarkingRulesResourcePath)
}

func dscpMarkingRulesResourceURL(c gophercloud.Client, policyID, ruleID string) string {
	return c.ServiceURL(rootPath, policyID, dscpMarkingRulesResourcePath, ruleID)
}

func listDSCPMarkingRulesURL(c gophercloud.Client, policyID string) string {
	return dscpMarkingRulesRootURL(c, policyID)
}

func getDSCPMarkingRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return dscpMarkingRulesResourceURL(c, policyID, ruleID)
}

func createDSCPMarkingRuleURL(c gophercloud.Client, policyID string) string {
	return dscpMarkingRulesRootURL(c, policyID)
}

func updateDSCPMarkingRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return dscpMarkingRulesResourceURL(c, policyID, ruleID)
}

func deleteDSCPMarkingRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return dscpMarkingRulesResourceURL(c, policyID, ruleID)
}

func minimumBandwidthRulesRootURL(c gophercloud.Client, policyID string) string {
	return c.ServiceURL(rootPath, policyID, minimumBandwidthRulesResourcePath)
}

func minimumBandwidthRulesResourceURL(c gophercloud.Client, policyID, ruleID string) string {
	return c.ServiceURL(rootPath, policyID, minimumBandwidthRulesResourcePath, ruleID)
}

func listMinimumBandwidthRulesURL(c gophercloud.Client, policyID string) string {
	return minimumBandwidthRulesRootURL(c, policyID)
}

func getMinimumBandwidthRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return minimumBandwidthRulesResourceURL(c, policyID, ruleID)
}

func createMinimumBandwidthRuleURL(c gophercloud.Client, policyID string) string {
	return minimumBandwidthRulesRootURL(c, policyID)
}

func updateMinimumBandwidthRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return minimumBandwidthRulesResourceURL(c, policyID, ruleID)
}

func deleteMinimumBandwidthRuleURL(c gophercloud.Client, policyID, ruleID string) string {
	return minimumBandwidthRulesResourceURL(c, policyID, ruleID)
}
