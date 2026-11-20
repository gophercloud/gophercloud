/*
Package oidc authenticates to Keystone through the OS-FEDERATION API using
OpenID Connect client credentials.

Authentication obtains an access token from the identity provider, exchanges
it for an unscoped Keystone token, and optionally scopes the result. The OIDC
token endpoint can be configured directly with AccessTokenEndpoint or resolved
from DiscoveryEndpoint.

Example:

	authOptions := &oidc.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		ClientID:             "my-client-id",
		ClientSecret:         "my-client-secret",
		AccessTokenEndpoint:  "https://idp.example.com/oauth2/token",
		Scope: tokens.Scope{
			ProjectName: "my-project",
			DomainName:  "Default",
		},
	}

	result := oidc.Create(context.TODO(), identityClient, authOptions)
	token, err := result.ExtractToken()

For the reference behavior, see the OpenStack 2025.1 keystoneauth OIDC plugin:
https://opendev.org/openstack/keystoneauth/src/branch/stable/2025.1/keystoneauth1/identity/v3/oidc.py
*/
package oidc
