package v2

import "fmt"

// ErrNotImplemented errors may occur in two contexts:
// (1) development versions of this package may return this error for endpoints which are defined but not yet completed, and,
// (2) production versions of this package may return this error when a provider fails to offer the requested Identity extension.
//
// ErrEndpoint errors occur when the authentication URL provided to Authenticate() either isn't valid
// or the endpoint provided doesn't respond like an Identity V2 API endpoint should.
//
// ErrCredentials errors occur when authentication fails due to the caller possessing insufficient access privileges.
var (
	ErrNotImplemented = fmt.Errorf("Identity feature not yet implemented")
	ErrEndpoint       = fmt.Errorf("Improper or missing Identity endpoint")
	ErrCredentials    = fmt.Errorf("Improper or missing Identity credentials")
)
