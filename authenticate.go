package gophercloud

type AuthOptions struct {
	Username, Password, TenantId string
}

func (c *Context) Authenticate(provider string, options AuthOptions) (*int, error) {
	_, err := c.ProviderByName(provider)
	if err != nil {
		return nil, err
	}

	if (options.Username == "") || (options.Password == "") {
		return nil, ErrCredentials
	}
	return nil, nil
}
