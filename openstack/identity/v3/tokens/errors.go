package tokens

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

func unacceptedAttributeErr(attribute string) string {
	return fmt.Sprintf("The base Identity V3 API does not accept authentication by %s", attribute)
}

func redundantWithTokenErr(attribute string) string {
	return fmt.Sprintf("%s may not be provided when authenticating with a TokenID", attribute)
}

func redundantWithUserID(attribute string) string {
	return fmt.Sprintf("%s may not be provided when authenticating with a UserID", attribute)
}

// ErrAPIKeyProvided indicates that an APIKey was provided but can't be used.
type ErrAPIKeyProvided struct{ gophercloud.BaseError }

func (e ErrAPIKeyProvided) Error() string {
	return unacceptedAttributeErr("APIKey")
}

// ErrTenantIDProvided indicates that a TenantID was provided but can't be used.
type ErrTenantIDProvided struct{ gophercloud.BaseError }

func (e ErrTenantIDProvided) Error() string {
	return unacceptedAttributeErr("TenantID")
}

// ErrTenantNameProvided indicates that a TenantName was provided but can't be used.
type ErrTenantNameProvided struct{ gophercloud.BaseError }

func (e ErrTenantNameProvided) Error() string {
	return unacceptedAttributeErr("TenantName")
}

// ErrUsernameWithToken indicates that a Username was provided, but token authentication is being used instead.
type ErrUsernameWithToken struct{ gophercloud.BaseError }

func (e ErrUsernameWithToken) Error() string {
	return redundantWithTokenErr("Username")
}

// ErrUserIDWithToken indicates that a UserID was provided, but token authentication is being used instead.
type ErrUserIDWithToken struct{ gophercloud.BaseError }

func (e ErrUserIDWithToken) Error() string {
	return redundantWithTokenErr("UserID")
}

// ErrDomainIDWithToken indicates that a DomainID was provided, but token authentication is being used instead.
type ErrDomainIDWithToken struct{ gophercloud.BaseError }

func (e ErrDomainIDWithToken) Error() string {
	return redundantWithTokenErr("DomainID")
}

// ErrDomainNameWithToken indicates that a DomainName was provided, but token authentication is being used instead.s
type ErrDomainNameWithToken struct{ gophercloud.BaseError }

func (e ErrDomainNameWithToken) Error() string {
	return redundantWithTokenErr("DomainName")
}

// ErrUsernameOrUserID indicates that neither username nor userID are specified, or both are at once.
type ErrUsernameOrUserID struct{ gophercloud.BaseError }

func (e ErrUsernameOrUserID) Error() string {
	return "Exactly one of Username and UserID must be provided for password authentication"
}

// ErrDomainIDWithUserID indicates that a DomainID was provided, but unnecessary because a UserID is being used.
type ErrDomainIDWithUserID struct{ gophercloud.BaseError }

func (e ErrDomainIDWithUserID) Error() string {
	return redundantWithUserID("DomainID")
}

// ErrDomainNameWithUserID indicates that a DomainName was provided, but unnecessary because a UserID is being used.
type ErrDomainNameWithUserID struct{ gophercloud.BaseError }

func (e ErrDomainNameWithUserID) Error() string {
	return redundantWithUserID("DomainName")
}

// ErrDomainIDOrDomainName indicates that a username was provided, but no domain to scope it.
// It may also indicate that both a DomainID and a DomainName were provided at once.
type ErrDomainIDOrDomainName struct{ gophercloud.BaseError }

func (e ErrDomainIDOrDomainName) Error() string {
	return "You must provide exactly one of DomainID or DomainName to authenticate by Username"
}

// ErrMissingPassword indicates that no password was provided and no token is available.
type ErrMissingPassword struct{ gophercloud.BaseError }

func (e ErrMissingPassword) Error() string {
	return "You must provide a password to authenticate"
}

// ErrScopeDomainIDOrDomainName indicates that a domain ID or Name was required in a Scope, but not present.
type ErrScopeDomainIDOrDomainName struct{ gophercloud.BaseError }

func (e ErrScopeDomainIDOrDomainName) Error() string {
	return "You must provide exactly one of DomainID or DomainName in a Scope with ProjectName"
}

// ErrScopeProjectIDOrProjectName indicates that both a ProjectID and a ProjectName were provided in a Scope.
type ErrScopeProjectIDOrProjectName struct{ gophercloud.BaseError }

func (e ErrScopeProjectIDOrProjectName) Error() string {
	return "You must provide at most one of ProjectID or ProjectName in a Scope"
}

// ErrScopeProjectIDAlone indicates that a ProjectID was provided with other constraints in a Scope.
type ErrScopeProjectIDAlone struct{ gophercloud.BaseError }

func (e ErrScopeProjectIDAlone) Error() string {
	return "ProjectID must be supplied alone in a Scope"
}

// ErrScopeDomainName indicates that a DomainName was provided alone in a Scope.
type ErrScopeDomainName struct{ gophercloud.BaseError }

func (e ErrScopeDomainName) Error() string {
	return "DomainName must be supplied with a ProjectName or ProjectID in a Scope"
}

// ErrScopeEmpty indicates that no credentials were provided in a Scope.
type ErrScopeEmpty struct{ gophercloud.BaseError }

func (e ErrScopeEmpty) Error() string {
	return "You must provide either a Project or Domain in a Scope"
}
