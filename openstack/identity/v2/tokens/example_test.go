package tokens_test

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
)

type tokenRequest struct {
	// Body is the complete JSON request, including its outer envelope.
	Body map[string]any
	// Headers contains additional request headers.
	Headers map[string]string
}

func (opts tokenRequest) ToTokenCreateMap() (map[string]any, error) {
	return opts.Body, nil
}

func (opts tokenRequest) ToTokenCreateHeaders(method, url string) (map[string]string, error) {
	return opts.Headers, nil
}

func createToken(ctx context.Context, identityClient *gophercloud.ServiceClient) (*tokens.Token, error) {
	request := tokenRequest{
		Body: map[string]any{"auth": map[string]any{
			"passwordCredentials": map[string]any{"username": "username", "password": "password"},
			"tenantId":            "tenant-id",
		}},
		Headers: map[string]string{"Accept": "application/json"},
	}
	return tokens.Create(ctx, identityClient, request).ExtractToken()
}

func ExampleCreate() {
	identityClient := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "https://identity.example.com/v2.0/",
	}
	_, err := createToken(context.Background(), identityClient)
	if err != nil {
		panic(err)
	}
}
