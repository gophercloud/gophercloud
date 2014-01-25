package notificationPlans

import (
	"fmt"
	"github.com/rackspace/gophercloud/rackspace/monitoring"
	"github.com/racker/perigee"
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

	url := fmt.Sprintf("%s/notification_plans/%s", c.options.Endpoint, id)
	err := perigee.Delete(url, perigee.Options{
		Results: &dr,
		OkCodes: []int{204},
	})
	return dr, err
}
