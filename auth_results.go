package gophercloud

// AuthResults encapsulates the raw results from an authentication request. As OpenStack allows
// extensions to influence the structure returned in ways that Gophercloud cannot predict at
// compile-time, you should use type-safe accessors to work with the data represented by this type,
// such as ServiceCatalog() and TokenID().
type AuthResults interface {

	// Retrieve the authentication token's value from the authentication response.
	GetTokenID() (string, error)
}
