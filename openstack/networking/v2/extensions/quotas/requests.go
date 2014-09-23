package quotas

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

func Get(c *gophercloud.ServiceClient) GetResult {
	var res GetResult
	_, err := perigee.Request("GET", rootURL(c), perigee.Options{
		Results:     &res.Resp,
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
	})
	res.Err = err
	return res
}

type UpdateOpts struct {
	Subnet        *int
	Router        *int
	Network       *int
	FloatingIP    *int
	Port          *int
	HealthMonitor *int
	SecGroupRule  *int
	SecGroup      *int
	VIP           *int
	Member        *int
	Pool          *int
}

func Update(c *gophercloud.ServiceClient, opts UpdateOpts) UpdateResult {
	type quota struct {
		Subnet        *int `json:"subnet,omitempty"`
		Router        *int `json:"router,omitempty"`
		Network       *int `json:"network,omitempty"`
		FloatingIP    *int `json:"floatingip,omitempty"`
		Port          *int `json:"port,omitempty"`
		HealthMonitor *int `json:"health_monitor,omitempty"`
		SecGroupRule  *int `json:"security_group_rule,omitempty"`
		VIP           *int `json:"vip,omitempty"`
		SecGroup      *int `json:"security_group,omitempty"`
		Member        *int `json:"member,omitempty"`
		Pool          *int `json:"pool,omitempty"`
	}

	type request struct {
		Quota quota `json:"quota"`
	}

	reqBody := request{Quota: quota{
		Subnet:        opts.Subnet,
		Router:        opts.Router,
		Network:       opts.Network,
		FloatingIP:    opts.FloatingIP,
		Port:          opts.Port,
		HealthMonitor: opts.HealthMonitor,
		SecGroupRule:  opts.SecGroupRule,
		VIP:           opts.VIP,
		SecGroup:      opts.SecGroup,
		Member:        opts.Member,
		Pool:          opts.Pool,
	}}

	var res UpdateResult
	_, err := perigee.Request("PUT", rootURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Resp,
		OkCodes:     []int{200},
	})
	res.Err = err
	return res
}

func Reset(c *gophercloud.ServiceClient) DeleteResult {
	var res DeleteResult
	_, err := perigee.Request("DELETE", rootURL(c), perigee.Options{
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	res.Err = err
	return res
}
