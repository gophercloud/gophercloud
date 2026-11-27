package gophercloud

// AuthResult is the result of acquiring a provider client's Keystone token.
// ProviderClient.GetAuthResult returns this interface.
//
// The normal authentication flow stores an auth.AuthResult value or pointer
// from github.com/gophercloud/gophercloud/v2/auth. It contains the token,
// expiry, user, scope, and service catalog. Low-level Identity token results
// also satisfy this interface, but are not used by the auth package.
//
// Implementations expose ExtractTokenID so that ProviderClient can store a
// token without depending on a particular authentication implementation.
type AuthResult interface {
	ExtractTokenID() (string, error)
}
