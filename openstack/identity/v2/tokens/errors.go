package tokens

import "fmt"

var (
	// ErrUserIDProvided is returned if you attempt to authenticate with a UserID.
	ErrUserIDProvided = unacceptedAttributeErr("UserID")

	// ErrAPIKeyProvided is returned if you attempt to authenticate with an APIKey.
	ErrAPIKeyProvided = unacceptedAttributeErr("APIKey")

	// ErrDomainIDProvided is returned if you attempt to authenticate with a DomainID.
	ErrDomainIDProvided = unacceptedAttributeErr("DomainID")

	// ErrDomainNameProvided is returned if you attempt to authenticate with a DomainName.
	ErrDomainNameProvided = unacceptedAttributeErr("DomainName")
)

func unacceptedAttributeErr(attribute string) error {
	return fmt.Errorf("The base Identity V2 API does not accept authentication by %s", attribute)
}
