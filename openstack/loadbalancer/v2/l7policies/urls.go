package l7policies

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "lbaas"
	resourcePath = "l7policies"
	rulePath     = "rules"
)

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func ruleRootURL(c gophercloud.Client, policyID string) string {
	return c.ServiceURL(rootPath, resourcePath, policyID, rulePath)
}

func ruleResourceURL(c gophercloud.Client, policyID string, ruleID string) string {
	return c.ServiceURL(rootPath, resourcePath, policyID, rulePath, ruleID)
}
