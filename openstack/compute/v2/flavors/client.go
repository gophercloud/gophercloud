package flavors

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/rackspace/gophercloud"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
)

type Client struct {
	endpoint  string
	authority identity.AuthResults
	options   gophercloud.AuthOptions
}

func NewClient(e string, a identity.AuthResults, ao gophercloud.AuthOptions) *Client {
	return &Client{
		endpoint:  e,
		authority: a,
		options:   ao,
	}
}

func (c *Client) getListUrl(lfo ListFilterOptions) string {
	v := url.Values{}
	if lfo.ChangesSince != "" {
		v.Set("changes-since", lfo.ChangesSince)
	}
	if lfo.MinDisk != 0 {
		v.Set("minDisk", strconv.Itoa(lfo.MinDisk))
	}
	if lfo.MinRam != 0 {
		v.Set("minRam", strconv.Itoa(lfo.MinRam))
	}
	if lfo.Marker != "" {
		v.Set("marker", lfo.Marker)
	}
	if lfo.Limit != 0 {
		v.Set("limit", strconv.Itoa(lfo.Limit))
	}
	tail := ""
	if len(v) > 0 {
		tail = fmt.Sprintf("?%s", v.Encode())
	}
	return fmt.Sprintf("%s/flavors/detail%s", c.endpoint, tail)
}

func (c *Client) getGetUrl(id string) string {
	return fmt.Sprintf("%s/flavors/%s", c.endpoint, id)
}

func (c *Client) getListHeaders() (map[string]string, error) {
	t, err := identity.GetToken(c.authority)
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		"X-Auth-Token": t.ID,
	}, nil
}
