package tokens

import (
	"errors"
	"fmt"
)

var (
	// ErrUserIDProvided is returned if you attempt to authenticate with a UserID.
	ErrUserIDProvided = unacceptedAttributeErr("UserID")

	// ErrDomainIDProvided is returned if you attempt to authenticate with a DomainID.
	ErrDomainIDProvided = unacceptedAttributeErr("DomainID")

	// ErrDomainNameProvided is returned if you attempt to authenticate with a DomainName.
	ErrDomainNameProvided = unacceptedAttributeErr("DomainName")

	// ErrUsernameRequired is returned if you attempt ot authenticate without a Username.
	ErrUsernameRequired = errors.New("You must supply a Username in your AuthOptions.")

	// ErrPasswordOrAPIKey is returned if you provide both a password and an API key.
	ErrPasswordOrAPIKey = errors.New("Please supply exactly one of Password or APIKey in your AuthOptions.")
)

func unacceptedAttributeErr(attribute string) error {
	return fmt.Errorf("The base Identity V2 API does not accept authentication by %s", attribute)
}
