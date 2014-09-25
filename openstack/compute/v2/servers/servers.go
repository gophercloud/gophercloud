package servers

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// ErrCannotInterpret is returned by an Extract call if the response body doesn't have the expected structure.
var ErrCannotInterpet = errors.New("Unable to interpret a response body.")

// Server exposes only the standard OpenStack fields corresponding to a given server on the user's account.
type Server struct {
	// ID uniquely identifies this server amongst all other servers, including those not accessible to the current tenant.
	ID string

	// TenantID identifies the tenant owning this server resource.
	TenantID string `mapstructure:"tenant_id"`

	// UserID uniquely identifies the user account owning the tenant.
	UserID string `mapstructure:"user_id"`

	// Name contains the human-readable name for the server.
	Name string

	// Updated and Created contain ISO-8601 timestamps of when the state of the server last changed, and when it was created.
	Updated string
	Created string

	HostID string

	// Status contains the current operational status of the server, such as IN_PROGRESS or ACTIVE.
	Status string

	// Progress ranges from 0..100.
	// A request made against the server completes only once Progress reaches 100.
	Progress int

	// AccessIPv4 and AccessIPv6 contain the IP addresses of the server, suitable for remote access for administration.
	AccessIPv4 string
	AccessIPv6 string

	// Image refers to a JSON object, which itself indicates the OS image used to deploy the server.
	Image map[string]interface{}

	// Flavor refers to a JSON object, which itself indicates the hardware configuration of the deployed server.
	Flavor map[string]interface{}

	// Addresses includes a list of all IP addresses assigned to the server, keyed by pool.
	Addresses map[string]interface{}

	// Metadata includes a list of all user-specified key-value pairs attached to the server.
	Metadata map[string]interface{}

	// Links includes HTTP references to the itself, useful for passing along to other APIs that might want a server reference.
	Links []interface{}

	// KeyName indicates which public key was injected into the server on launch.
	KeyName string `mapstructure:"keyname"`

	// AdminPass will generally be empty ("").  However, it will contain the administrative password chosen when provisioning a new server without a set AdminPass setting in the first place.
	// Note that this is the ONLY time this field will be valid.
	AdminPass string `mapstructure:"adminPass"`
}

// ExtractServers interprets the results of a single page from a List() call, producing a slice of Server entities.
func ExtractServers(page pagination.Page) ([]Server, error) {
	casted := page.(ListPage).Body

	var response struct {
		Servers []Server `mapstructure:"servers"`
	}
	err := mapstructure.Decode(casted, &response)
	return response.Servers, err
}

// ExtractServer interprets the result of a call expected to return data on a single server.
func ExtractServer(sr ServerResult) (*Server, error) {
	so, ok := sr["server"]
	if !ok {
		return nil, ErrCannotInterpet
	}
	serverObj := so.(map[string]interface{})
	s := new(Server)
	err := mapstructure.Decode(serverObj, s)
	return s, err
}
