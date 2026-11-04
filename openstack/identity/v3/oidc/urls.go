package oidc

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath              = "OS-FEDERATION"
	identityProvidersPath = "identity_providers"
	protocolsPath         = "protocols"
	authPath              = "auth"
)

func authURL(c *gophercloud.ServiceClient, idpID, protocolID string) string {
	return c.ServiceURL(rootPath, identityProvidersPath, idpID, protocolsPath, protocolID, authPath)
}
