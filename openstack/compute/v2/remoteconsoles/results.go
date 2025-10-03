package remoteconsoles

import "github.com/gophercloud/gophercloud/v2"

type commonResult struct {
	gophercloud.Result
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a RemoteConsole.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Console
type GetResult struct {
	commonResult
}

// RemoteConsole represents the Compute service remote console object.
type RemoteConsole struct {
	// Protocol contains remote console protocol.
	// You can use the RemoteConsoleProtocol custom type to unmarshal raw JSON
	// response into the pre-defined valid console protocol.
	Protocol string `json:"protocol"`

	// Type contains remote console type.
	// You can use the RemoteConsoleType custom type to unmarshal raw JSON
	// response into the pre-defined valid console type.
	Type string `json:"type"`

	// URL can be used to connect to the remote console.
	URL string `json:"url"`
}

// Console represents the Compute service console object.
type Console struct {
	// InstanceUUID contains the UUID of the server.
	InstanceUUID string `json:"instance_uuid"`

	// Host contains the name or ID of the host.
	// (Optional)
	Host string `json:"host"`

	// Port contains the port number
	Port int `json:"port"`

	// TLSPort contains the port number of a port requiring the TLS connection.
	// (Optional)
	TLSPort int `json:"tls_port"`

	// InternalAccessPath contains the id representing the internal access path.
	//(Optional)
	InternalAccessPath string `json:"internal_access_path"`
}

// Extract interprets any commonResult as a RemoteConsole.
func (r commonResult) Extract() (*RemoteConsole, error) {
	var s struct {
		RemoteConsole *RemoteConsole `json:"remote_console"`
	}
	err := r.ExtractInto(&s)
	return s.RemoteConsole, err
}

// Extract interprets any commonResult as a Console.
func (r GetResult) Extract() (*Console, error) {
	var s struct {
		Console *Console `json:"console"`
	}
	err := r.ExtractInto(&s)
	return s.Console, err
}
