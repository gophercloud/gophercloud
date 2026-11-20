package websso

import "github.com/gophercloud/gophercloud/v2"

func authURL(client *gophercloud.ServiceClient, identityProvider, protocol string) string {
	return client.ServiceURL(
		"auth",
		"OS-FEDERATION",
		"identity_providers",
		identityProvider,
		"protocols",
		protocol,
		"websso",
	)
}
