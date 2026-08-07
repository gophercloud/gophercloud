package websso

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath       = "auth"
	federationPath = "OS-FEDERATION"
	idpPath        = "identity_providers"
	protocolsPath  = "protocols"
	webssoPath     = "websso"
	tokensPath     = "tokens"
)

// webssoURL builds the URL for initiating WebSSO authentication
func webssoURL(c *gophercloud.ServiceClient, idp, protocol string) string {
	return c.ServiceURL(rootPath, federationPath, idpPath, idp, protocolsPath, protocol, webssoPath)
}

// tokenURL builds the URL for token operations
func tokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, tokensPath)
}
