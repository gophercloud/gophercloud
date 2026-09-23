package tokens_test

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
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
			"identity": map[string]any{
				"methods":  []string{"password"},
				"password": map[string]any{"user": map[string]any{"id": "user-id", "password": "password"}},
			},
			"scope": map[string]any{"project": map[string]any{"id": "project-id"}},
		}},
		Headers: map[string]string{"Accept": "application/json"},
	}
	return tokens.Create(ctx, identityClient, request).ExtractToken()
}

func ExampleCreate() {
	identityClient := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       "https://identity.example.com/v3/",
	}
	_, err := createToken(context.Background(), identityClient)
	if err != nil {
		panic(err)
	}
}
