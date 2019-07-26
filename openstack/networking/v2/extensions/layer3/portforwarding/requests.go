package portforwarding

import "github.com/gophercloud/gophercloud"

// CreateOpts contains all the values needed to create a new port forwarding
// resource. All attributes are required.
type CreateOpts struct {
	InternalPortID    string `json:"internal_port_id"`
	InternalIPAddress string `json:"internal_ip_address"`
	InternalPort      int    `json:"internal_port"`
	ExternalPort      int    `json:"external_port"`
	Protocol          string `json:"protocol"`
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPortForwardingCreateMap() (map[string]interface{}, error)
}

// ToPortForwardingCreateMap allows CreateOpts to satisfy the CreateOptsBuilder
// interface
func (opts CreateOpts) ToPortForwardingCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "port_forwarding")
}

// Create accepts a CreateOpts struct and uses the values provided to create a
// new port forwarding for an existing floating IP.
func Create(c *gophercloud.ServiceClient, floatingIpId string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPortForwardingCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(portForwardingUrl(c, floatingIpId), b, &r.Body, nil)
	return
}

// Delete will permanently delete a particular port forwarding for a given floating ID.
func Delete(c *gophercloud.ServiceClient, floatingIpId string, pfId string) (r DeleteResult) {
	_, r.Err = c.Delete(singlePortForwardingUrl(c, floatingIpId, pfId), nil)
	return
}
