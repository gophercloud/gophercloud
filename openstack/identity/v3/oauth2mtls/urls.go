package oauth2mtls

import "github.com/gophercloud/gophercloud/v2"

func tokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("OS-OAUTH2", "token")
}

func validateURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
