package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

// ClientOption customizes a *gophercloud.ProviderClient built by
// NewClientFromAuthResult.
type ClientOption func(*gophercloud.ProviderClient)

// NewClientFromAuthResult reconstructs a *gophercloud.ProviderClient from a
// previously obtained auth.AuthResult (e.g. one returned by
// auth.AuthOptionsV2/V3.Authenticate and cached across processes) without
// making any network call. The returned client has no ReauthFunc wired -
// there are no credentials here to reauthenticate with - unless WithReauth
// is passed.
func NewClientFromAuthResult(result auth.AuthResult, opts ...ClientOption) (*gophercloud.ProviderClient, error) {
	client := new(gophercloud.ProviderClient)
	client.UseTokenLock()
	if err := client.SetTokenAndAuthResult(result); err != nil {
		return nil, err
	}
	client.EndpointLocator = result.EndpointLocator()
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

// WithReauth wires client.ReauthFunc using options, for callers who kept
// both the auth.AuthResult and the auth.AuthOptionsBuilder that produced
// it. Does nothing if options.CanReauth() is false.
func WithReauth(options auth.Authenticator, authResult auth.AuthResult) ClientOption {
	return func(client *gophercloud.ProviderClient) {
		if !authResult.CanReauth {
			return
		}
		client.ReauthFunc = func(ctx context.Context) error {
			result, err := options.Authenticate(ctx, &client.HTTPClient)
			if err != nil {
				return err
			}
			if err := client.SetTokenAndAuthResult(result); err != nil {
				return err
			}
			client.EndpointLocator = result.EndpointLocator()
			return nil
		}
	}
}
