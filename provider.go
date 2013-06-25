package gophercloud

import (
)

type Provider struct {
	// empty.
}

var providerMap = make(map[string]*Provider)

func (c *Context) RegisterProvider(name string, p *Provider) error {
	c.providerMap[name] = p
	return nil
}

func (c *Context) ProviderByName(name string) (p *Provider, err error) {
	for provider, descriptor := range c.providerMap {
		if name == provider {
			return descriptor, nil
		}
	}
	return nil, ErrProvider
}