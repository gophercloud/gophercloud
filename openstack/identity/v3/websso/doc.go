/*
Package websso authenticates through Keystone's browser-based WebSSO flow.

Authenticate starts a loopback HTTP listener, opens Keystone's WebSSO
endpoint, accepts Keystone's form POST, and returns a standard Identity v3
token result. The callback origin must be registered in Keystone's
trusted_dashboard configuration. RedirectHost is restricted to loopback
addresses so the token listener cannot be exposed on an external interface.

Keystone's WebSSO callback has no request correlation value. Authenticate
validates the received token with Keystone, but cannot prove that the POST
came from the browser flow it opened. Another valid token submitted during the
callback window can therefore cause login CSRF.

Token caching is disabled by default. When TokenCache is set, the validated
unscoped token is reused across scopes. CacheNamespace must identify the local
login profile so unrelated browser identities cannot share a token. Changing
the active browser account does not invalidate a cached token; use a different
namespace, clear the cache, or disable caching when switching accounts.

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
