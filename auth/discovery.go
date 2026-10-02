package auth

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
)

// authenticator defers identity version discovery until the caller's
// context and configured HTTP client are available.
type authenticator struct {
	v2 AuthOptionsV2
	v3 AuthOptionsV3
}

func (a authenticator) GetAuthURL() string {
	return a.v3.AuthURL
}

func (a authenticator) Authenticate(ctx context.Context, provider *gophercloud.ProviderClient) (*AuthResult, error) {
	base, err := utils.BaseEndpoint(a.GetAuthURL())
	if err != nil {
		return nil, err
	}

	client := &gophercloud.ProviderClient{
		IdentityBase:     gophercloud.NormalizeURL(base),
		IdentityEndpoint: gophercloud.NormalizeURL(a.GetAuthURL()),
	}
	if provider != nil {
		client.HTTPClient = provider.HTTPClient
		client.UserAgent = provider.UserAgent
		client.RetryBackoffFunc = provider.RetryBackoffFunc
		client.MaxBackoffRetries = provider.MaxBackoffRetries
		client.RetryFunc = provider.RetryFunc
	}

	version, endpoint, err := utils.ChooseVersion(ctx, client, []*utils.Version{
		{ID: "v2.0", Priority: 20, Suffix: "/v2.0/"},
		{ID: "v3", Priority: 30, Suffix: "/v3/"},
	})
	if err != nil {
		return nil, err
	}

	switch version.ID {
	case "v2.0":
		return a.v2.authenticate(ctx, provider, endpoint)
	default:
		return a.v3.authenticate(ctx, provider, endpoint)
	}
}
