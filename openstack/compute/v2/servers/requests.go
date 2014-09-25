package servers

import (
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

	return pagination.NewPager(client, detailURL(client), createPage)
}

// Create requests a server to be provisioned to the user in the current tenant.
func Create(client *gophercloud.ServiceClient, opts map[string]interface{}) CreateResult {
	var result CreateResult
	_, result.Err = perigee.Request("POST", listURL(client), perigee.Options{
		Results: &result.Resp,
		ReqBody: map[string]interface{}{
			"server": opts,
		},
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return result
}

// Delete requests that a server previously provisioned be removed from your account.
func Delete(client *gophercloud.ServiceClient, id string) error {
	_, err := perigee.Request("DELETE", serverURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return err
}

// Get requests details on a single server, by ID.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var result GetResult
	_, result.Err = perigee.Request("GET", serverURL(client, id), perigee.Options{
		Results:     &result.Resp,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return result
}

// Update requests that various attributes of the indicated server be changed.
func Update(client *gophercloud.ServiceClient, id string, opts map[string]interface{}) UpdateResult {
	var result UpdateResult
	_, result.Err = perigee.Request("PUT", serverURL(client, id), perigee.Options{
		Results: &result.Resp,
		ReqBody: map[string]interface{}{
			"server": opts,
		},
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
