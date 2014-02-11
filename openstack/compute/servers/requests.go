package servers

import (
	"github.com/racker/perigee"
)

// ListResult abstracts the raw results of making a List() request against the
// API.  As OpenStack extensions may freely alter the response bodies of
// structures returned to the client, you may only safely access the data
// provided through separate, type-safe accessors or methods. 
type ListResult map[string]interface{}

// List makes a request against the API to list servers accessible to you.
func List(c *Client) (ListResult, error) {
	var lr ListResult

	h, err := c.getListHeaders()
	if err != nil {
		return nil, err
	}

	err = perigee.Get(c.getListUrl(), perigee.Options{
		Results: &lr,
		MoreHeaders: h,
	})
	return lr, err
}

