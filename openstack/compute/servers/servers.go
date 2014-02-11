package servers

import (
	"github.com/mitchellh/mapstructure"
	"fmt"
)

// ErrNotImplemented indicates a failure to discover a feature of the response from the API.
// E.g., a missing server field, a missing extension, etc.
var ErrNotImplemented = fmt.Errorf("Compute Servers feature not implemented.")

// Server exposes only the standard OpenStack fields corresponding to a given server on the user's account.
//
// Id uniquely identifies this server amongst all other servers, including those not accessible to the current tenant.
//
// TenantId identifies the tenant owning this server resource.
//
// UserId uniquely identifies the user account owning the tenant.
//
// Name contains the human-readable name for the server.
//
// Updated and Created contain ISO-8601 timestamps of when the state of the server last changed, and when it was created.
//
// Status contains the current operational status of the server, such as IN_PROGRESS or ACTIVE.
//
// Progress ranges from 0..100.  A request made against the server completes only once Progress reaches 100.
//
// AccessIPv4 and AccessIPv6 contain the IP addresses of the server, suitable for remote access for administration.
//
// Image refers to a JSON object, which itself indicates the OS image used to deploy the server.
//
// Flavor refers to a JSON object, which itself indicates the hardware configuration of the deployed server.
//
// Addresses includes a list of all IP addresses assigned to the server, keyed by pool.
//
// Metadata includes a list of all user-specified key-value pairs attached to the server.
//
// Links includes HTTP references to the itself, useful for passing along to other APIs that might want a server reference.
type Server struct {
	Id string
	TenantId string `mapstructure:tenant_id`
	UserId string `mapstructure:user_id`
	Name string
	Updated string
	Created string
	HostId string
	Status string
	Progress int
	AccessIPv4 string
	AccessIPv6 string
	Image map[string]interface{}
	Flavor map[string]interface{}
	Addresses map[string]interface{}
	Metadata map[string]interface{}
	Links []interface{}
}

// GetServers interprets the result of a List() call, producing a slice of Server entities.
func GetServers(lr ListResult) ([]Server, error) {
	sa, ok := lr["servers"]
	if !ok {
		return nil, ErrNotImplemented
	}
	serversArray := sa.([]interface{})

	servers := make([]Server, len(serversArray))
	for i, so := range serversArray {
		serverObj := so.(map[string]interface{})
		err := mapstructure.Decode(serverObj, &servers[i])
		if err != nil {
			return servers, err
		}
	}

	return servers, nil
}

