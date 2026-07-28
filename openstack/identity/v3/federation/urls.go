package federation

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath              = "OS-FEDERATION"
	identityProvidersPath = "identity_providers"
	mappingsPath          = "mappings"
	protocolsPath         = "protocols"
	serviceProvidersPath  = "service_providers"
)

func identityProvidersRootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, identityProvidersPath)
}

func identityProvidersResourceURL(c *gophercloud.ServiceClient, identityProviderID string) string {
	return c.ServiceURL(rootPath, identityProvidersPath, identityProviderID)
}

func mappingsRootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, mappingsPath)
}

func mappingsResourceURL(c *gophercloud.ServiceClient, mappingID string) string {
	return c.ServiceURL(rootPath, mappingsPath, mappingID)
}

func protocolsRootURL(c *gophercloud.ServiceClient, identityProviderID string) string {
	return c.ServiceURL(rootPath, identityProvidersPath, identityProviderID, protocolsPath)
}

func protocolsResourceURL(c *gophercloud.ServiceClient, identityProviderID, protocolID string) string {
	return c.ServiceURL(rootPath, identityProvidersPath, identityProviderID, protocolsPath, protocolID)
}

func serviceProvidersRootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, serviceProvidersPath)
}

func serviceProvidersResourceURL(c *gophercloud.ServiceClient, serviceProviderID string) string {
	return c.ServiceURL(rootPath, serviceProvidersPath, serviceProviderID)
}
