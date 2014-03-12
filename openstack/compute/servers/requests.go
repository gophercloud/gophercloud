package servers

import (
	"github.com/racker/perigee"
	"fmt"
)

// ListResult abstracts the raw results of making a List() request against the
// API.  As OpenStack extensions may freely alter the response bodies of
// structures returned to the client, you may only safely access the data
// provided through separate, type-safe accessors or methods.
type ListResult map[string]interface{}

// ServerResult abstracts a single server description,
// as returned by the OpenStack provider.
// As OpenStack extensions may freely alter the response bodies of the
// structures returned to the client,
// you may only safely access the data provided through
// separate, type-safe accessors or methods.
type ServerResult map[string]interface{}

// List makes a request against the API to list servers accessible to you.
func List(c *Client) (ListResult, error) {
	var lr ListResult

	h, err := c.getListHeaders()
	if err != nil {
		return nil, err
	}

	err = perigee.Get(c.getListUrl(), perigee.Options{
		Results:     &lr,
		MoreHeaders: h,
	})
	return lr, err
}

// Create requests a server to be provisioned to the user in the current tenant.
func Create(c *Client, opts map[string]interface{}) (ServerResult, error) {
	var sr ServerResult

	h, err := c.getCreateHeaders()
	if err != nil {
		return nil, err
	}

	err = perigee.Post(c.getCreateUrl(), perigee.Options{
		Results: &sr,
		ReqBody: map[string]interface{}{
			"server": opts,
		},
		MoreHeaders: h,
		OkCodes:     []int{202},
	})
	return sr, err
}

// Delete requests that a server previously provisioned be removed from your account.
func Delete(c *Client, id string) error {
	h, err := c.getDeleteHeaders()
	if err != nil {
		return err
	}

	err = perigee.Delete(c.getDeleteUrl(id), perigee.Options{
		MoreHeaders: h,
		OkCodes:     []int{204},
	})
	return err
}

// GetDetail requests details on a single server, by ID.
func GetDetail(c *Client, id string) (ServerResult, error) {
	var sr ServerResult

	h, err := c.getDetailHeaders()
	if err != nil {
		return nil, err
	}

	err = perigee.Get(c.getDetailUrl(id), perigee.Options{
		Results:     &sr,
		MoreHeaders: h,
	})
	return sr, err
}

// Update requests that various attributes of the indicated server be changed.
func Update(c *Client, id string, opts map[string]interface{}) (ServerResult, error) {
	var sr ServerResult

	h, err := c.getUpdateHeaders()
	if err != nil {
		return nil, err
	}

	err = perigee.Put(c.getUpdateUrl(id), perigee.Options{
		Results: &sr,
		ReqBody: map[string]interface{}{
			"server": opts,
		},
		MoreHeaders: h,
	})
	return sr, err
}

// ChangeAdminPassword alters the administrator or root password for a specified
// server.
func ChangeAdminPassword(c *Client, id, newPassword string) error {
	h, err := c.getActionHeaders()
	if err != nil {
		return err
	}

	err = perigee.Post(c.getActionUrl(id), perigee.Options{
		ReqBody: struct{C map[string]string `json:"changePassword"`}{
			map[string]string{"adminPass": newPassword},
		},
		MoreHeaders: h,
		OkCodes: []int{202},
	})
	return err
}

// ArgumentError errors occur when an argument supplied to a package function
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
type ArgumentError struct {
	Function, Argument string
	Value interface{}
}

// Error yields a useful diagnostic for debugging purposes.
func (e *ArgumentError) Error() string {
	return fmt.Sprintf("Bad argument in call to %s, formal parameter %s, value %#v", e.Function, e.Argument, e.Value)
}

func (e *ArgumentError) String() string {
	return e.Error()
}

// These constants determine how a server should be rebooted.
// See the Reboot() function for further details.
const (
	SoftReboot = "SOFT"
	HardReboot = "HARD"
	OSReboot = SoftReboot
	PowerCycle = HardReboot
)

// Reboot requests that a given server reboot.
// Two methods exist for rebooting a server:
//
// HardReboot (aka PowerCycle) -- restarts the server instance by physically
// cutting power to the machine, or if a VM, terminating it at the hypervisor
// level.  It's done.  Caput.  Full stop.  Then, after a brief while, power is
// restored or the VM instance restarted.
//
// SoftReboot (aka OSReboot).  This approach simply tells the OS to restart
// under its own procedures.  E.g., in Linux, asking it to enter runlevel 6,
// or executing "sudo shutdown -r now", or by wasking Windows to restart the
// machine.
func Reboot(c *Client, id, how string) error {
	if (how != SoftReboot) && (how != HardReboot) {
		return &ArgumentError{
			Function: "Reboot",
			Argument: "how",
			Value: how,
		}
	}
	
	h, err := c.getActionHeaders()
	if err != nil {
		return err
	}

	err = perigee.Post(c.getActionUrl(id), perigee.Options{
		ReqBody: struct{C map[string]string `json:"reboot"`}{
			map[string]string{"type": how},
		},
		MoreHeaders: h,
		OkCodes: []int{202},
	})
	return err
}
