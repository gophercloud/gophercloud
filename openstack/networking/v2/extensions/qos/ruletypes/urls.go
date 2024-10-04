package ruletypes

import "github.com/gophercloud/gophercloud/v2"

func listRuleTypesURL(c gophercloud.Client) string {
	return c.ServiceURL("qos", "rule-types")
}

func getRuleTypeURL(c gophercloud.Client, name string) string {
	return c.ServiceURL("qos", "rule-types", name)
}
