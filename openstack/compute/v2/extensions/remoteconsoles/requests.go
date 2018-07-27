package remoteconsoles

import (
	"github.com/gophercloud/gophercloud"
)

// RemoteConsoleProtocol represents valid remote console protocol.
// It can be used to create a remote console with one of the pre-defined protocol.
type RemoteConsoleProtocol string

const (
	// RemoteConsoleVNCProtocol represents the VNC console protocol.
	RemoteConsoleVNCProtocol RemoteConsoleProtocol = "vnc"

	// RemoteConsoleSPICEProtocol represents the SPICE console protocol.
	RemoteConsoleSPICEProtocol RemoteConsoleProtocol = "spice"

	// RemoteConsoleRDPProtocol represents the RDP console protocol.
	RemoteConsoleRDPProtocol RemoteConsoleProtocol = "rdp"

	// RemoteConsoleSerialProtocol represents the serial console protocol.
	RemoteConsoleSerialProtocol RemoteConsoleProtocol = "serial"

	// RemoteConsoleMKSProtocol represents the MKS console protocol.
	RemoteConsoleMKSProtocol RemoteConsoleProtocol = "mks"
)

// RemoteConsoleType represents valid remote console type.
// It can be used to create a remote console with one of the pre-defined type.
type RemoteConsoleType string

const (
	// RemoteConsoleNoVNCType represents the VNC console type.
	RemoteConsoleNoVNCType RemoteConsoleType = "novnc"

	// RemoteConsoleXVPVNCType represents the XVP VNC console type.
	RemoteConsoleXVPVNCType RemoteConsoleType = "xvpvnc"

	// RemoteConsoleRDPHTML5Type represents the RDP HTML5 console type.
	RemoteConsoleRDPHTML5Type RemoteConsoleType = "rdp-html5"

	// RemoteConsoleSPICEHTML5Type represents the SPICE HTML5 console type.
	RemoteConsoleSPICEHTML5Type RemoteConsoleType = "spice-html5"

	// RemoteConsoleSerialType represents the serial console type.
	RemoteConsoleSerialType RemoteConsoleType = "serial"

	// RemoteConsoleWebMKSType represents the web MKS console type.
	RemoteConsoleWebMKSType RemoteConsoleType = "webmks"
)

// CreateOptsBuilder allows to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToRemoteConsoleCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies parameters to the Create request.
type CreateOpts struct {
	// Protocol specifies the protocol of a new remote console.
	Protocol string `json:"protocol" required:"true"`

	// Type specifies the type of a new remote console.
	Type string `json:"type" required:"true"`
}

// ToRemoteConsoleCreateMap builds a request body from the CreateOpts.
func (opts CreateOpts) ToRemoteConsoleCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "remote_console")
}

// Create requests the creation of a new remote console on the specified server.
func Create(client *gophercloud.ServiceClient, serverID string, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToRemoteConsoleCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client, serverID), reqBody, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
