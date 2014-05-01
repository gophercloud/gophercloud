package volumes

import (
	"github.com/racker/perigee"
	blockstorage "github.com/rackspace/gophercloud/openstack/blockstorage/v1"
)

func Create(c *blockstorage.Client, opts CreateOpts) (Volume, error) {
	var v Volume
	h, err := c.GetHeaders()
	if err != nil {
		return v, err
	}
	url := c.GetVolumesURL()
	_, err = perigee.Request("POST", url, perigee.Options{
		Results: &v,
		ReqBody: map[string]interface{}{
			"volume": opts,
		},
		MoreHeaders: h,
	})
	return v, err
}
