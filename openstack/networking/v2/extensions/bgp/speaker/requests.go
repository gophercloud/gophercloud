package speaker

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// For the sake of consistency, this function will return a pagination.Pager
func List(c *gophercloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BGPSpeakerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieve the specific bgp speaker by its uuid
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
