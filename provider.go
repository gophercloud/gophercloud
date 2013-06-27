package gophercloud

// Provider structures exist for each tangible provider of OpenStack service.
// For example, Rackspace, Hewlett-Packard, and NASA might have their own instance of this structure.
//
// At a minimum, a provider must expose an authentication endpoint.
type Provider struct {
	AuthEndpoint string
}

// RegisterProvider allows a unit test to register a mythical provider convenient for testing.
// If the provider structure lacks adequate configuration, or the configuration given has some
// detectable error, an ErrConfiguration error will result.
func (c *Context) RegisterProvider(name string, p Provider) error {
	if p.AuthEndpoint == "" {
		return ErrConfiguration
	}

	c.providerMap[name] = p
	return nil
}

// ProviderByName will locate a provider amongst those previously registered, if it exists.
// If the named provider has not been registered, an ErrProvider error will result.
func (c *Context) ProviderByName(name string) (p Provider, err error) {
	for provider, descriptor := range c.providerMap {
		if name == provider {
			return descriptor, nil
		}
	}
	return Provider{}, ErrProvider
}
