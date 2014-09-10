package notificationPlans

import (
	"fmt"

	"github.com/racker/perigee"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
	"github.com/rackspace/gophercloud/rackspace/monitoring"
)

var ErrNotImplemented = fmt.Errorf("notificationPlans feature not yet implemented")

type Client struct {
	options monitoring.Options
}

type DeleteResults map[string]interface{}

func NewClient(mo monitoring.Options) *Client {
	return &Client{
		options: mo,
	}
}

func (c *Client) Delete(id string) (DeleteResults, error) {
	var dr DeleteResults

	tok, err := identity.GetToken(c.options.Authentication)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/notification_plans/%s", c.options.Endpoint, id)
	err = perigee.Delete(url, perigee.Options{
		Results: &dr,
		OkCodes: []int{204},
		MoreHeaders: map[string]string{
			"X-Auth-Token": tok.ID,
		},
	})
	return dr, err
}
