package base

import "github.com/gophercloud/gophercloud"

// Get retrieves the home document, allowing the user to discover the
// entire API.
func Get(c *gophercloud.ServiceClient) GetResult {
	var r GetResult
	_, r.Err = c.Get(getURL(c), &r.Body, nil)
	return r
}

// Ping retrieves a ping to the server.
func Ping(c *gophercloud.ServiceClient) PingResult {
	var r PingResult
	_, r.Err = c.Get(pingURL(c), nil, &gophercloud.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: map[string]string{"Accept": ""},
	})
	return r
}
