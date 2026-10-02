package auth_test

import (
	"context"

	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
)

func ExampleAuthOptionsFromEnv() {
	opts, err := auth.AuthOptionsFromEnv()
	if err != nil {
		panic(err)
	}
	provider, err := openstack.AuthenticatedClient(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	_ = provider
}

func ExampleAuthOptionsFromCloud() {
	cloud, endpointOpts, tlsConfig, err := clouds.Parse()
	if err != nil {
		panic(err)
	}
	opts, err := auth.AuthOptionsFromCloud(cloud, auth.WithTLSConfig(tlsConfig))
	if err != nil {
		panic(err)
	}
	provider, err := config.NewProviderClient(context.Background(), opts, config.WithTLSConfig(tlsConfig))
	if err != nil {
		panic(err)
	}
	network, err := openstack.NewNetworkV2(context.Background(), provider, endpointOpts)
	if err != nil {
		panic(err)
	}

	_ = network
}

func ExampleV3OAuth1Opts() {
	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3OAuth1Opts{
			ConsumerKey:     "consumer-key",
			ConsumerSecret:  "consumer-secret",
			Token:           "access-token",
			TokenSecret:     "access-secret",
			SignatureMethod: auth.OAuth1HMACSHA1,
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = result
}

func ExampleAuthResult() {
	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3PasswordOpts{
			UserID: "user-id", Password: "password", AllowReauth: true,
			Scope: &auth.Scope{ProjectID: "project-id"},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	provider, err := openstack.NewClientFromAuthResult(
		*result, openstack.WithReauth(opts, *result),
	)
	if err != nil {
		panic(err)
	}
	_ = provider
}
