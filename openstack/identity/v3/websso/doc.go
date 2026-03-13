/*
Package websso provides WebSSO (Web Single Sign-On) authentication support
for federated identity providers using the OS-FEDERATION API.

WebSSO authentication enables browser-based OpenID Connect authentication
with OpenStack Keystone. This is commonly used with identity providers like
Keycloak, Okta, or other OIDC-compliant providers.

# Authentication Flow

The WebSSO authentication flow works as follows:

 1. The client opens a browser to Keystone's WebSSO endpoint
 2. The user authenticates via the identity provider
 3. The identity provider redirects back to a local callback server
 4. The callback server extracts the authentication token
 5. The token is used to configure the ProviderClient

# Example Usage

Basic authentication with an identity provider:

	import (
		"context"

		"github.com/gophercloud/gophercloud/v2"
		"github.com/gophercloud/gophercloud/v2/openstack"
		"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/websso"
	)

	// Create an unauthenticated provider client
	provider, err := openstack.NewClient("https://keystone.example.org:5000/v3")
	if err != nil {
		panic(err)
	}

	// Configure WebSSO authentication options
	opts := websso.AuthOptions{
		IdentityEndpoint: "https://keystone.example.org:5000/v3",
		IdentityProvider: "my-idp",
		Protocol:         "openid",
		AllowReauth:      true,
	}

	// Authenticate - this will open a browser for the user to log in
	err = websso.Authenticate(context.TODO(), provider, opts)
	if err != nil {
		panic(err)
	}

	// The provider is now authenticated and can be used to create service clients
	computeClient, err := openstack.NewComputeV2(context.TODO(), provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

# Scoped Authentication

To authenticate with a specific project or domain scope:

	opts := websso.AuthOptions{
		IdentityEndpoint: "https://keystone.example.org:5000/v3",
		IdentityProvider: "my-idp",
		Protocol:         "openid",
		Scope: gophercloud.AuthScope{
			ProjectName: "my-project",
			DomainName:  "Default",
		},
		AllowReauth: true,
	}

	err := websso.Authenticate(context.TODO(), provider, opts)

# Keystone Configuration

The Keystone server must be configured to support WebSSO. In keystone.conf:

	[federation]
	trusted_dashboard = http://localhost:9990/auth/websso/

The callback URL (localhost:9990 by default) must be listed as a trusted dashboard.

# Customizing the Callback Server

The callback server's host and port can be customized:

	opts := websso.AuthOptions{
		IdentityEndpoint: "https://keystone.example.org:5000/v3",
		IdentityProvider: "my-idp",
		Protocol:         "openid",
		RedirectHost:     "127.0.0.1",
		RedirectPort:     8080,
	}

# Token Caching

Tokens are cached locally to avoid repeated browser authentication. The cache
is stored in ~/.cache/gophercloud/ by default. To customize the cache location:

	opts := websso.AuthOptions{
		IdentityEndpoint: "https://keystone.example.org:5000/v3",
		IdentityProvider: "my-idp",
		Protocol:         "openid",
		CachePath:        "/custom/cache/path",
	}

Cached tokens are automatically validated for expiration before reuse.
*/
package websso
