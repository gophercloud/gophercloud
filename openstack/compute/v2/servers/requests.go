package servers

import (
	"encoding/base64"
	"fmt"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// List makes a request against the API to list servers accessible to you.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ServerPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	}

	return pagination.NewPager(client, listDetailURL(client), createPage)
}

// CreateOptsBuilder describes struct types that can be accepted by the Create call.
// The CreateOpts struct in this package does.
type CreateOptsBuilder interface {
	ToServerCreateMap() map[string]interface{}
}

// Network is used within CreateOpts to control a new server's network attachments.
type Network struct {
	// UUID of a nova-network to attach to the newly provisioned server.
	// Required unless Port is provided.
	UUID string

	// Port of a neutron network to attach to the newly provisioned server.
	// Required unless UUID is provided.
	Port string

	// FixedIP [optional] specifies a fixed IPv4 address to be used on this network.
	FixedIP string
}

// CreateOpts specifies server creation parameters.
type CreateOpts struct {
	// Name [required] is the name to assign to the newly launched server.
	Name string

	// ImageRef [required] is the ID or full URL to the image that contains the server's OS and initial state.
	// Optional if using the boot-from-volume extension.
	ImageRef string

	// FlavorRef [required] is the ID or full URL to the flavor that describes the server's specs.
	FlavorRef string

	// SecurityGroups [optional] lists the names of the security groups to which this server should belong.
	SecurityGroups []string

	// UserData [optional] contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you.
	UserData []byte

	// AvailabilityZone [optional] in which to launch the server.
	AvailabilityZone string

	// Networks [optional] dictates how this server will be attached to available networks.
	// By default, the server will be attached to all isolated networks for the tenant.
	Networks []Network

	// Metadata [optional] contains key-value pairs (up to 255 bytes each) to attach to the server.
	Metadata map[string]string

	// Personality [optional] includes the path and contents of a file to inject into the server at launch.
	// The maximum size of the file is 255 bytes (decoded).
	Personality []byte

	// ConfigDrive [optional] enables metadata injection through a configuration drive.
	ConfigDrive bool
}

// ToServerCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToServerCreateMap() map[string]interface{} {
	server := make(map[string]interface{})

	server["name"] = opts.Name
	server["imageRef"] = opts.ImageRef
	server["flavorRef"] = opts.FlavorRef

	if opts.UserData != nil {
		encoded := base64.StdEncoding.EncodeToString(opts.UserData)
		server["user_data"] = &encoded
	}
	if opts.Personality != nil {
		encoded := base64.StdEncoding.EncodeToString(opts.Personality)
		server["personality"] = &encoded
	}
	if opts.ConfigDrive {
		server["config_drive"] = "true"
	}
	if opts.AvailabilityZone != "" {
		server["availability_zone"] = opts.AvailabilityZone
	}
	if opts.Metadata != nil {
		server["metadata"] = opts.Metadata
	}

	if len(opts.SecurityGroups) > 0 {
		securityGroups := make([]map[string]interface{}, len(opts.SecurityGroups))
		for i, groupName := range opts.SecurityGroups {
			securityGroups[i] = map[string]interface{}{"name": groupName}
		}
	}
	if len(opts.Networks) > 0 {
		networks := make([]map[string]interface{}, len(opts.Networks))
		for i, net := range opts.Networks {
			networks[i] = make(map[string]interface{})
			if net.UUID != "" {
				networks[i]["uuid"] = net.UUID
			}
			if net.Port != "" {
				networks[i]["port"] = net.Port
			}
			if net.FixedIP != "" {
				networks[i]["fixed_ip"] = net.FixedIP
			}
		}
	}

	return map[string]interface{}{"server": server}
}

// Create requests a server to be provisioned to the user in the current tenant.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var result CreateResult
	_, result.Err = perigee.Request("POST", listURL(client), perigee.Options{
		Results:     &result.Resp,
		ReqBody:     opts.ToServerCreateMap(),
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return result
}

// Delete requests that a server previously provisioned be removed from your account.
func Delete(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", deleteURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}

// Get requests details on a single server, by ID.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var result GetResult
	_, result.Err = perigee.Request("GET", getURL(client, id), perigee.Options{
		Results:     &result.Resp,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return result
}

// UpdateOptsBuilder allows extentions to add additional attributes to the Update request.
type UpdateOptsBuilder interface {
	ToServerUpdateMap() map[string]interface{}
}

// UpdateOpts specifies the base attributes that may be updated on an existing server.
type UpdateOpts struct {
	// Name [optional] changes the displayed name of the server.
	// The server host name will *not* change.
	// Server names are not constrained to be unique, even within the same tenant.
	Name string

	// AccessIPv4 [optional] provides a new IPv4 address for the instance.
	AccessIPv4 string

	// AccessIPv6 [optional] provides a new IPv6 address for the instance.
	AccessIPv6 string
}

// ToServerUpdateMap formats an UpdateOpts structure into a request body.
func (opts UpdateOpts) ToServerUpdateMap() map[string]interface{} {
	server := make(map[string]string)
	if opts.Name != "" {
		server["name"] = opts.Name
	}
	if opts.AccessIPv4 != "" {
		server["accessIPv4"] = opts.AccessIPv4
	}
	if opts.AccessIPv6 != "" {
		server["accessIPv6"] = opts.AccessIPv6
	}
	return map[string]interface{}{"server": server}
}

// Update requests that various attributes of the indicated server be changed.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) UpdateResult {
	var result UpdateResult
	_, result.Err = perigee.Request("PUT", updateURL(client, id), perigee.Options{
		Results:     &result.Resp,
		ReqBody:     opts.ToServerUpdateMap(),
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return result
}

// ChangeAdminPassword alters the administrator or root password for a specified server.
func ChangeAdminPassword(client *gophercloud.ServiceClient, id, newPassword string) error {
	var req struct {
		ChangePassword struct {
			AdminPass string `json:"adminPass"`
		} `json:"changePassword"`
	}

	req.ChangePassword.AdminPass = newPassword

	_, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		ReqBody:     req,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return err
}

// ErrArgument errors occur when an argument supplied to a package function
// fails to fall within acceptable values.  For example, the Reboot() function
// expects the "how" parameter to be one of HardReboot or SoftReboot.  These
// constants are (currently) strings, leading someone to wonder if they can pass
// other string values instead, perhaps in an effort to break the API of their
// provider.  Reboot() returns this error in this situation.
//
// Function identifies which function was called/which function is generating
// the error.
// Argument identifies which formal argument was responsible for producing the
// error.
// Value provides the value as it was passed into the function.
type ErrArgument struct {
	Function, Argument string
	Value              interface{}
}

// Error yields a useful diagnostic for debugging purposes.
func (e *ErrArgument) Error() string {
	return fmt.Sprintf("Bad argument in call to %s, formal parameter %s, value %#v", e.Function, e.Argument, e.Value)
}

func (e *ErrArgument) String() string {
	return e.Error()
}

// RebootMethod describes the mechanisms by which a server reboot can be requested.
type RebootMethod string

// These constants determine how a server should be rebooted.
// See the Reboot() function for further details.
const (
	SoftReboot RebootMethod = "SOFT"
	HardReboot RebootMethod = "HARD"
	OSReboot                = SoftReboot
	PowerCycle              = HardReboot
)

// Reboot requests that a given server reboot.
// Two methods exist for rebooting a server:
//
// HardReboot (aka PowerCycle) restarts the server instance by physically cutting power to the machine, or if a VM,
// terminating it at the hypervisor level.
// It's done. Caput. Full stop.
// Then, after a brief while, power is restored or the VM instance restarted.
//
// SoftReboot (aka OSReboot) simply tells the OS to restart under its own procedures.
// E.g., in Linux, asking it to enter runlevel 6, or executing "sudo shutdown -r now", or by asking Windows to restart the machine.
func Reboot(client *gophercloud.ServiceClient, id string, how RebootMethod) error {
	if (how != SoftReboot) && (how != HardReboot) {
		return &ErrArgument{
			Function: "Reboot",
			Argument: "how",
			Value:    how,
		}
	}

	_, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		ReqBody: struct {
			C map[string]string `json:"reboot"`
		}{
			map[string]string{"type": string(how)},
		},
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return err
}

// Rebuild requests that the Openstack provider reprovision the server.
// The rebuild will need to know the server's name and new image reference or ID.
// In addition, and unlike building a server with Create(), you must provide an administrator password.
//
// Additional options may be specified with the additional map.
// This function treats a nil map the same as an empty map.
//
// Rebuild returns a server result as though you had called GetDetail() on the server's ID.
// The information, however, refers to the new server, not the old.
func Rebuild(client *gophercloud.ServiceClient, id, name, password, imageRef string, additional map[string]interface{}) RebuildResult {
	var result RebuildResult

	if id == "" {
		result.Err = &ErrArgument{
			Function: "Rebuild",
			Argument: "id",
			Value:    "",
		}
		return result
	}

	if name == "" {
		result.Err = &ErrArgument{
			Function: "Rebuild",
			Argument: "name",
			Value:    "",
		}
		return result
	}

	if password == "" {
		result.Err = &ErrArgument{
			Function: "Rebuild",
			Argument: "password",
			Value:    "",
		}
		return result
	}

	if imageRef == "" {
		result.Err = &ErrArgument{
			Function: "Rebuild",
			Argument: "imageRef",
			Value:    "",
		}
		return result
	}

	if additional == nil {
		additional = make(map[string]interface{}, 0)
	}

	additional["name"] = name
	additional["imageRef"] = imageRef
	additional["adminPass"] = password

	_, result.Err = perigee.Request("POST", actionURL(client, id), perigee.Options{
		ReqBody: struct {
			R map[string]interface{} `json:"rebuild"`
		}{
			additional,
		},
		Results:     &result.Resp,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return result
}

// Resize instructs the provider to change the flavor of the server.
// Note that this implies rebuilding it.
// Unfortunately, one cannot pass rebuild parameters to the resize function.
// When the resize completes, the server will be in RESIZE_VERIFY state.
// While in this state, you can explore the use of the new server's configuration.
// If you like it, call ConfirmResize() to commit the resize permanently.
// Otherwise, call RevertResize() to restore the old configuration.
func Resize(client *gophercloud.ServiceClient, id, flavorRef string) error {
	_, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		ReqBody: struct {
			R map[string]interface{} `json:"resize"`
		}{
			map[string]interface{}{"flavorRef": flavorRef},
		},
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return err
}

// ConfirmResize confirms a previous resize operation on a server.
// See Resize() for more details.
func ConfirmResize(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		ReqBody:     map[string]interface{}{"confirmResize": nil},
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}

// RevertResize cancels a previous resize operation on a server.
// See Resize() for more details.
func RevertResize(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("POST", actionURL(client, id), perigee.Options{
		ReqBody:     map[string]interface{}{"revertResize": nil},
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return err
}
