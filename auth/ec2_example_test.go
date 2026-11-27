package auth_test

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

func ExampleAuthOptionsEC2_Authenticate() {
	var serviceProvider *gophercloud.ProviderClient // Configure service credentials separately.
	opts := auth.AuthOptionsEC2{
		AuthURL:         "https://identity.example.com/v3",
		ServiceProvider: serviceProvider,
		Auth: auth.EC2TokenOpts{
			Access: "access-key",
			Secret: "secret-key",
			Host:   "ec2.example.com",
			Path:   "/",
			Verb:   "GET",
			Params: map[string]string{
				"SignatureVersion": "2",
				"SignatureMethod":  auth.EC2CredentialsHmacSha256V2,
			},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = result
}
