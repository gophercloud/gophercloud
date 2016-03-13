package gophercloud

import "fmt"

// BaseError is an error type that all other error types embed.
type BaseError struct {
	Info string
	//Function string
}

func (e BaseError) Error() string {
	return "An error occurred while executing a Gophercloud request."
}

// ErrMissingInput is the error when input is required in a particular
// situation but not provided by the user
type ErrMissingInput struct {
	BaseError
	Argument string
}

func (e ErrMissingInput) Error() string {
	return fmt.Sprintf("Missing input for argument [%s]", e.Argument)
}

// ErrInvalidInput is an error type used for most non-HTTP Gophercloud errors.
type ErrInvalidInput struct {
	ErrMissingInput
	Value interface{}
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("Invalid input provided for argument [%s]: [%+v]", e.Argument, e.Value)
}

// ErrUnexpectedResponseCode is returned by the Request method when a response code other than
// those listed in OkCodes is encountered.
type ErrUnexpectedResponseCode struct {
	BaseError
	URL      string
	Method   string
	Expected []int
	Actual   int
	Body     []byte
}

func (err ErrUnexpectedResponseCode) Error() string {
	return fmt.Sprintf(
		"Expected HTTP response code %v when accessing [%s %s], but got %d instead\n%s",
		err.Expected, err.Method, err.URL, err.Actual, err.Body,
	)
}

// ErrDefault400 is the default error type returned on a 400 HTTP response code.
type ErrDefault400 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault401 is the default error type returned on a 401 HTTP response code.
type ErrDefault401 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault404 is the default error type returned on a 404 HTTP response code.
type ErrDefault404 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault405 is the default error type returned on a 405 HTTP response code.
type ErrDefault405 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault408 is the default error type returned on a 408 HTTP response code.
type ErrDefault408 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault429 is the default error type returned on a 429 HTTP response code.
type ErrDefault429 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault500 is the default error type returned on a 500 HTTP response code.
type ErrDefault500 struct {
	ErrUnexpectedResponseCode
}

// ErrDefault503 is the default error type returned on a 503 HTTP response code.
type ErrDefault503 struct {
	ErrUnexpectedResponseCode
}

func (e ErrDefault400) Error() string {
	return "Invalid request due to incorrect syntax or missing required parameters."
}
func (e ErrDefault401) Error() string {
	return "Authentication failed"
}
func (e ErrDefault404) Error() string {
	return "Resource not found"
}
func (e ErrDefault405) Error() string {
	return "Method not allowed"
}
func (e ErrDefault408) Error() string {
	return "The server timed out waiting for the request"
}
func (e ErrDefault429) Error() string {
	return "Too many requests have been sent in a given amount of time. Pause" +
		" requests, wait up to one minute, and try again."
}
func (e ErrDefault500) Error() string {
	return "Internal Server Error"
}
func (e ErrDefault503) Error() string {
	return "The service is currently unable to handle the request due to a temporary" +
		" overloading or maintenance. This is a temporary condition. Try again later."
}

// Err400er is the interface resource error types implement to override the error message
// from a 400 error.
type Err400er interface {
	Error400(ErrUnexpectedResponseCode) error
}

// Err401er is the interface resource error types implement to override the error message
// from a 401 error.
type Err401er interface {
	Error401(ErrUnexpectedResponseCode) error
}

// Err404er is the interface resource error types implement to override the error message
// from a 404 error.
type Err404er interface {
	Error404(ErrUnexpectedResponseCode) error
}

// Err405er is the interface resource error types implement to override the error message
// from a 405 error.
type Err405er interface {
	Error405(ErrUnexpectedResponseCode) error
}

// Err408er is the interface resource error types implement to override the error message
// from a 408 error.
type Err408er interface {
	Error408(ErrUnexpectedResponseCode) error
}

// Err429er is the interface resource error types implement to override the error message
// from a 429 error.
type Err429er interface {
	Error429(ErrUnexpectedResponseCode) error
}

// Err500er is the interface resource error types implement to override the error message
// from a 500 error.
type Err500er interface {
	Error500(ErrUnexpectedResponseCode) error
}

// Err503er is the interface resource error types implement to override the error message
// from a 503 error.
type Err503er interface {
	Error503(ErrUnexpectedResponseCode) error
}

// ErrTimeOut is the error type returned when an operations times out.
type ErrTimeOut struct {
	BaseError
}

func (e ErrTimeOut) Error() string {
	return "A time out occurred"
}

// ErrUnableToReauthenticate is the error type returned when reauthentication fails.
type ErrUnableToReauthenticate struct {
	BaseError
	ErrOriginal error
}

func (e ErrUnableToReauthenticate) Error() string {
	return fmt.Sprintf("Unable to re-authenticate: %s", e.ErrOriginal)
}

// ErrErrorAfterReauthentication is the error type returned when reauthentication
// succeeds, but an error occurs afterword (usually an HTTP error).
type ErrErrorAfterReauthentication struct {
	BaseError
	ErrOriginal error
}

func (e ErrErrorAfterReauthentication) Error() string {
	return fmt.Sprintf("Successfully re-authenticated, but got error executing request: %s", e.ErrOriginal)
}

// ErrServiceNotFound is returned when no service in a service catalog matches
// the provided EndpointOpts. This is generally returned by provider service
// factory methods like "NewComputeV2()" and can mean that a service is not
// enabled for your account.
type ErrServiceNotFound struct {
	BaseError
}

func (e ErrServiceNotFound) Error() string {
	return "No suitable service could be found in the service catalog."
}

// ErrEndpointNotFound is returned when no available endpoints match the
// provided EndpointOpts. This is also generally returned by provider service
// factory methods, and usually indicates that a region was specified
// incorrectly.
type ErrEndpointNotFound struct {
	BaseError
}

func (e ErrEndpointNotFound) Error() string {
	return "No suitable endpoint could be found in the service catalog."
}

// ErrResourceNotFound is the error when trying to retrieve a resource's
// ID by name and the resource doesn't exist.
type ErrResourceNotFound struct {
	BaseError
	Name         string
	ResourceType string
}

func (e ErrResourceNotFound) Error() string {
	return fmt.Sprintf("Unable to find %s with name %s", e.ResourceType, e.Name)
}

// ErrMultipleResourcesFound is the error when trying to retrieve a resource's
// ID by name and multiple resources have the user-provided name.
type ErrMultipleResourcesFound struct {
	BaseError
	Name         string
	Count        int
	ResourceType string
}

func (e ErrMultipleResourcesFound) Error() string {
	return fmt.Sprintf("Found %d %ss matching %s", e.Count, e.ResourceType, e.Name)
}

// ErrUnexpectedType is the error when an unexpected type is encountered
type ErrUnexpectedType struct {
	BaseError
	Expected string
	Actual   string
}

func (e ErrUnexpectedType) Error() string {
	return fmt.Sprintf("Expected %s but got %s", e.Expected, e.Actual)
}
