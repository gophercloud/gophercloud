package sharenetworks

import "github.com/gophercloud/gophercloud"

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToShareNetworkCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a ShareNetwork. This object is
// passed to the sharenetworks.Create function. For more information about
// these parameters, see the ShareNetwork object.
type CreateOpts struct {
	// The UUID of the Neutron network to set up for share servers
	NeutronNetID string `json:"neutron_net_id,omitempty"`
	// The UUID of the Neutron subnet to set up for share servers
	NeutronSubnetID string `json:"neutron_subnet_id,omitempty"`
	// The UUID of the nova network to set up for share servers
	NovaNetID string `json:"nova_net_id,omitempty"`
	// The share network name
	Name string `json:"name"`
	// The share network description
	Description string `json:"description"`
}

// ToShareNetworkCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToShareNetworkCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "share_network")
}

// Create will create a new ShareNetwork based on the values in CreateOpts. To
// extract the ShareNetwork object from the response, call the Extract method
// on the CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToShareNetworkCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})
	return
}

// Delete will delete the existing ShareNetwork with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
