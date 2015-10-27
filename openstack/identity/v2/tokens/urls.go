package tokens

import "github.com/rackspace/gophercloud"

// CreateURL generates the URL used to create new Tokens.
func CreateURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("tokens")
}

// CreateGetURL generates the URL used to Validate Tokens.
func CreateGetURL(client *gophercloud.ServiceClient, token string) string {
    return client.ServiceURL("tokens", token)
}