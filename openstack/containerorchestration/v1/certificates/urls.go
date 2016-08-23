package certificates

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("certificates")
}

func getCertificateAuthorityURL(c *gophercloud.ServiceClient, bayID string) string {
	return c.ServiceURL("certificates", bayID)
}
