package gophercloud

import (
	"fmt"
	"net/http"
	"strings"
)

type AuthError struct{ error }
type CatalogError struct{ error }
type RequestError struct{ error }
type ValidationError struct{ error }

func (e AuthError) Unwrap() error       { return e.error }
func (e CatalogError) Unwrap() error    { return e.error }
func (e RequestError) Unwrap() error    { return e.error }
func (e ValidationError) Unwrap() error { return e.error }

// ErrMissingInput is the error when input is required in a particular
// situation but not provided by the user
type ErrMissingInput struct {
	Argument string
}

func (e ErrMissingInput) Error() string {
	return fmt.Sprintf("missing input for argument %q", e.Argument)
}

func NewMissingInputError(argument string) ValidationError {
	return ValidationError{ErrMissingInput{Argument: argument}}
}

// ErrInvalidInput is an error type used for most non-HTTP Gophercloud errors.
type ErrInvalidInput struct {
	Argument string
	Value    any
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input provided for argument %q: %+v", e.Argument, e.Value)
}

func NewInvalidInputError(argument string, value any) ValidationError {
	return ValidationError{ErrInvalidInput{Argument: argument, Value: value}}
}

// ErrConflictingInput is returned when two mutually exclusive arguments were
// passed.
type ErrConflictingInput struct {
	Argument1, Argument2 string
}

func (e ErrConflictingInput) Error() string {
	return fmt.Sprintf("exactly one of %s and %s must be provided", e.Argument1, e.Argument2)
}

func NewConflictingInputError(argument1, argument2 string) ValidationError {
	return ValidationError{ErrConflictingInput{Argument1: argument1, Argument2: argument2}}
}

// ErrMissingEnvironmentVariables is the error when environment variables are
// required in a particular situation but not provided by the user
type ErrMissingEnvironmentVariables struct {
	EnvironmentVariables []string
}

func (e ErrMissingEnvironmentVariables) Error() string {
	return fmt.Sprintf("missing required environment variables %s", strings.Join(e.EnvironmentVariables, ", "))
}

func NewMissingEnvironmentVariables(names ...string) ValidationError {
	return ValidationError{ErrMissingEnvironmentVariables{names}}
}

// ErrUnexpectedResponseCode is returned by the Request method when a response code other than
// those listed in OkCodes is encountered.
type ErrUnexpectedResponseCode struct {
	URL            string
	Method         string
	Expected       []int
	Actual         int
	Body           []byte
	ResponseHeader http.Header
	wrapped        ErrHTTPCode
}

func (e ErrUnexpectedResponseCode) Error() string {
	return fmt.Sprintf(
		"expected HTTP response code %v when accessing %s %q, but got %d instead\n%s",
		e.Expected, e.Method, e.URL, e.Actual, e.Body,
	)
}

func (e ErrUnexpectedResponseCode) Unwrap() error {
	return e.wrapped
}

func NewUnexpectedResponseCodeError(url string, method string, expected []int, actual int, body []byte, responseHeader http.Header) RequestError {
	return RequestError{ErrUnexpectedResponseCode{
		URL:            url,
		Method:         method,
		Expected:       expected,
		Actual:         actual,
		Body:           body,
		ResponseHeader: responseHeader,
		wrapped:        ErrHTTPCode(actual),
	}}
}

type ErrHTTPCode int

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "not found"
}

func (e ErrHTTPCode) Error() string {
	return fmt.Sprintf("unexpected HTTP response code %d %s", e, http.StatusText(int(e)))
}

func (e ErrHTTPCode) Unwrap() error {
	if e == http.StatusNotFound {
		return ErrNotFound{}
	}
	return nil
}

// ErrTimeout is the error type returned when an operation times out.
type ErrTimeout struct{}

func (e ErrTimeout) Error() string {
	return "timed out"
}

func NewTimeoutError() RequestError {
	return RequestError{ErrTimeout{}}
}

// ErrUnableToReauthenticate is the error type returned when reauthentication fails.
type ErrUnableToReauthenticate struct {
	wrapped []error
}

func (e ErrUnableToReauthenticate) Error() string {
	return fmt.Sprintf("unable to re-authenticate: %v", e.wrapped)
}

func (e ErrUnableToReauthenticate) Unwrap() []error {
	return e.wrapped
}

func NewReauthenticationError(errs ...error) AuthError {
	return AuthError{ErrUnableToReauthenticate{errs}}
}

// ErrServiceNotFound is returned when no service in a service catalog matches
// the provided EndpointOpts. This is generally returned by provider service
// factory methods like "NewComputeV2()" and can mean that a service is not
// enabled for your account.
type ErrServiceNotFound struct{}

func (e ErrServiceNotFound) Error() string {
	return "no suitable service could be found in the service catalog."
}

func NewServiceNotFoundError() CatalogError {
	return CatalogError{ErrServiceNotFound{}}
}

// ErrEndpointNotFound is returned when no available endpoints match the
// provided EndpointOpts. This is also generally returned by provider service
// factory methods, and usually indicates that a region was specified
// incorrectly.
type ErrEndpointNotFound struct{}

func (e ErrEndpointNotFound) Error() string {
	return "no suitable endpoint could be found in the service catalog."
}

func NewEndpointNotFoundError() CatalogError {
	return CatalogError{ErrEndpointNotFound{}}
}

// ErrMultipleResourcesFound is the error when trying to retrieve a resource's
// ID by name and multiple resources have the user-provided name.
type ErrMultipleResourcesFound struct {
	Name         string
	Count        int
	ResourceType string
}

func (e ErrMultipleResourcesFound) Error() string {
	return fmt.Sprintf("found more than one (%d) %s with name %q", e.Count, e.ResourceType, e.Name)
}

func NewMultipleResourcesFoundError(name, resourceType string, count int) RequestError {
	return RequestError{ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: resourceType}}
}

// ErrUnexpectedType is the error when an unexpected type is encountered
type ErrUnexpectedType struct {
	Expected string
	Actual   string
}

func (e ErrUnexpectedType) Error() string {
	return fmt.Sprintf("expected %s but got %s", e.Expected, e.Actual)
}

func NewUnexpectedTypeError(expected, actual string) RequestError {
	return RequestError{ErrUnexpectedType{Expected: expected, Actual: actual}}
}

type ErrInvalidAuthenticationMethod struct {
	was string
}

func (e ErrInvalidAuthenticationMethod) Error() string {
	return fmt.Sprintf("the Identity V3 API does not accept authentication by %s", e.was)
}

func NewInvalidAuthenticationMethodError(was string) AuthError {
	return AuthError{ErrInvalidAuthenticationMethod{was: was}}
}

type ErrFieldsIncompatibleWithToken struct {
	Redundant []string
}

func (e ErrFieldsIncompatibleWithToken) Error() string {
	return fmt.Sprintf("these fields may not be provided when authenticating with a TokenID: %v", e.Redundant)
}

func NewFieldsIncompatibleWithTokenError(redundant ...string) AuthError {
	return AuthError{ErrFieldsIncompatibleWithToken{Redundant: redundant}}
}

// // ErrAPIKeyProvided indicates that an APIKey was provided but can't be used.
// type ErrAPIKeyProvided struct{ RequestError }

// func (e ErrAPIKeyProvided) Error() string {
// 	return unacceptedAttributeErr("APIKey")
// }

// // ErrTenantIDProvided indicates that a TenantID was provided but can't be used.
// type ErrTenantIDProvided struct{ RequestError }

// func (e ErrTenantIDProvided) Error() string {
// 	return unacceptedAttributeErr("TenantID")
// }

// // ErrTenantNameProvided indicates that a TenantName was provided but can't be used.
// type ErrTenantNameProvided struct{ RequestError }

// func (e ErrTenantNameProvided) Error() string {
// 	return unacceptedAttributeErr("TenantName")
// }

// ErrUsernameOrUserID indicates that neither username nor userID are specified, or both are at once.
type ErrUsernameOrUserID struct{}

func (e ErrUsernameOrUserID) Error() string {
	return "exactly one of Username and UserID must be provided for password authentication"
}

func NewUsernameOrUserIDError() AuthError {
	return AuthError{ErrUsernameOrUserID{}}
}

// ErrDomainIDWithUserID indicates that a DomainID was provided, but unnecessary because a UserID is being used.
type ErrFieldsIncompatibleWithUserID struct {
	Redundant []string
}

func (e ErrFieldsIncompatibleWithUserID) Error() string {
	return fmt.Sprintf("these fields may not be provided when authenticating with a UserID: %v", e.Redundant)
}

func NewFieldsIncompatibleWithUserIDError(redundant ...string) AuthError {
	return AuthError{ErrFieldsIncompatibleWithUserID{Redundant: redundant}}
}

// ErrDomainIDOrDomainName indicates that a username was provided, but no
// domain to scope it. It may also indicate that both a DomainID and a
// DomainName were provided at once.
type ErrDomainNameOrDomainID struct{}

func (e ErrDomainNameOrDomainID) Error() string {
	return "exactly one of DomainID or DomainName must be provided to authenticate by Username"
}

func NewDomainNameOrDomainIDError() AuthError {
	return AuthError{ErrDomainNameOrDomainID{}}
}

// // ErrMissingPassword indicates that no password was provided and no token is available.
type ErrMissingPassword struct{}

func (e ErrMissingPassword) Error() string {
	return "you must provide a password to authenticate"
}

func NewMissingPasswordError() AuthError {
	return AuthError{ErrMissingPassword{}}
}

// ErrScopeDomainIDOrDomainName indicates that a domain ID or Name was required in a Scope, but not present.
type ErrScopeDomainIDOrDomainName struct{}

func (e ErrScopeDomainIDOrDomainName) Error() string {
	return "You must provide exactly one of DomainID or DomainName in a Scope with ProjectName"
}

func NewScopeDomainIDOrDomainNameError() AuthError {
	return AuthError{ErrScopeDomainIDOrDomainName{}}
}

// ErrScopeProjectIDOrProjectName indicates that both a ProjectID and a ProjectName were provided in a Scope.
type ErrScopeProjectIDOrProjectName struct{ RequestError }

func (e ErrScopeProjectIDOrProjectName) Error() string {
	return "You must provide at most one of ProjectID or ProjectName in a Scope"
}

func NewScopeProjectIDOrProjectNameError() AuthError {
	return AuthError{ErrScopeProjectIDOrProjectName{}}
}

// ErrScopeProjectIDAlone indicates that a ProjectID was provided with other constraints in a Scope.
type ErrScopeProjectIDAlone struct{}

func (e ErrScopeProjectIDAlone) Error() string {
	return "ProjectID must be supplied alone in a Scope"
}

func NewScopeProjectIDAloneError() AuthError {
	return AuthError{ErrScopeProjectIDAlone{}}
}

// // ErrScopeEmpty indicates that no credentials were provided in a Scope.
// type ErrScopeEmpty struct{ RequestError }

// func (e ErrScopeEmpty) Error() string {
// 	return "You must provide either a Project or Domain in a Scope"
// }

// // ErrAppCredMissingSecret indicates that no Application Credential Secret was provided with Application Credential ID or Name
type ErrAppCredMissingSecret struct{}

func (e ErrAppCredMissingSecret) Error() string {
	return "you must provide an application credential secret"
}

func NewAppCredMissingSecretError() AuthError {
	return AuthError{ErrAppCredMissingSecret{}}
}
