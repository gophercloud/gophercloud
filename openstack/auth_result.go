package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

// ClientOption customizes a *gophercloud.ProviderClient built by
// NewClientFromAuthResult.
type ProviderClientOption func(*gophercloud.ProviderClient)

// Manually creates a Provider Client by directly using auth.AuthResult
func NewClientFromAuthResult(result auth.AuthResult, opts ...ProviderClientOption) (*gophercloud.ProviderClient, error) {
	client := new(gophercloud.ProviderClient)
	client.UseTokenLock()
	if err := client.SetTokenAndAuthResult(result); err != nil {
		return nil, err
	}
	client.EndpointLocator = result.EndpointLocator(client)
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

// Convenience function to help implement ReauthFunc for users who configure provider
// clients manually with an auth.AuthResult
func WithReauth(options auth.Authenticator, authResult auth.AuthResult) ProviderClientOption {
	return func(client *gophercloud.ProviderClient) {
		if !authResult.CanReauth {
			return
		}
		throwaway := *client
		throwaway.SetThrowaway(true)
		throwaway.ReauthFunc = nil
		if err := throwaway.SetTokenAndAuthResult(nil); err != nil {
			return
		}
		client.ReauthFunc = func(ctx context.Context) error {
			result, err := options.Authenticate(ctx, &throwaway)
			if err != nil {
				return err
			}
			if err := client.SetTokenAndAuthResult(result); err != nil {
				return err
			}
			client.EndpointLocator = result.EndpointLocator(client)
			return nil
		}
	}
}
