/*
Package websso authenticates through Keystone's browser-based WebSSO flow.

Authenticate starts a loopback HTTP listener, opens Keystone's WebSSO
endpoint, accepts Keystone's form POST, and returns a standard Identity v3
token result. The callback origin must be registered in Keystone's
trusted_dashboard configuration. RedirectHost is restricted to loopback
addresses so the token listener cannot be exposed on an external interface.

Token caching is disabled by default. When TokenCache is set, CacheNamespace
must identify the local login profile so unrelated browser identities cannot
share a cached token.

Example:

	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		RedirectHost:         "localhost",
		RedirectPort:         9990,
		Scope: tokens.Scope{
			ProjectID: "project-id",
		},
	}

	result := websso.Authenticate(context.TODO(), identityClient, opts)
	if result.Err != nil {
		panic(result.Err)
	}
*/
package websso
