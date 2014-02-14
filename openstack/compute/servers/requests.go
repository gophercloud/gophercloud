package servers

import (
	"github.com/racker/perigee"
)

// ListResult abstracts the raw results of making a List() request against the
// API.  As OpenStack extensions may freely alter the response bodies of
// structures returned to the client, you may only safely access the data
// provided through separate, type-safe accessors or methods.
type ListResult map[string]interface{}

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
