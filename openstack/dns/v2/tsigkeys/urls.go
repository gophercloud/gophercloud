package tsigkeys

import "github.com/gophercloud/gophercloud/v2"

// baseURL returns the base URL for TSIG keys.
func baseURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("tsigkeys")
}

// tsigkeyURL returns the URL for a specific TSIG key.
func tsigkeyURL(c *gophercloud.ServiceClient, tsigkeyID string) string {
	return c.ServiceURL("tsigkeys", tsigkeyID)
}
