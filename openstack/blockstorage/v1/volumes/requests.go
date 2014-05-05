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
		Results: &struct {
			Volume *Volume `json:"volume"`
		}{&v},
		ReqBody: map[string]interface{}{
			"volume": opts,
		},
		MoreHeaders: h,
	})
	return v, err
}

func List(c *blockstorage.Client, opts ListOpts) ([]Volume, error) {
	var v []Volume
	var url string
	h, err := c.GetHeaders()
	if err != nil {
		return v, err
	}
	if full := opts["full"]; full {
		url = c.GetVolumesURL()
	} else {
		url = c.GetVolumeURL("detail")
	}
	_, err = perigee.Request("GET", url, perigee.Options{
		Results: &struct {
			Volume *[]Volume `json:"volumes"`
		}{&v},
		MoreHeaders: h,
	})
	return v, err
}

func Get(c *blockstorage.Client, opts GetOpts) (Volume, error) {
	var v Volume
	h, err := c.GetHeaders()
	if err != nil {
		return v, err
	}
	url := c.GetVolumeURL(opts["id"])
	_, err = perigee.Request("GET", url, perigee.Options{
		Results: &struct {
			Volume *Volume `json:"volume"`
		}{&v},
		MoreHeaders: h,
	})
	return v, err
}

func Delete(c *blockstorage.Client, opts DeleteOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}
	url := c.GetVolumeURL(opts["id"])
	_, err = perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: h,
	})
	return err
}
