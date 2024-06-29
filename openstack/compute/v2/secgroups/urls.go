package secgroups

import "github.com/gophercloud/gophercloud/v2"

const (
	secgrouppath = "os-security-groups"
	rulepath     = "os-security-group-rules"
)

func resourceURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(secgrouppath, id)
}

func rootURL(c gophercloud.Client) string {
	return c.ServiceURL(secgrouppath)
}

func listByServerURL(c gophercloud.Client, serverID string) string {
	return c.ServiceURL("servers", serverID, secgrouppath)
}

func rootRuleURL(c gophercloud.Client) string {
	return c.ServiceURL(rulepath)
}

func resourceRuleURL(c gophercloud.Client, id string) string {
	return c.ServiceURL(rulepath, id)
}

func serverActionURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("servers", id, "action")
}
