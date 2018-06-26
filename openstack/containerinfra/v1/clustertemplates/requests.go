package clustertemplates

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/common"
)

// Get retrieves a specific clustertemplate based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Get(getURL(c, id), &r.Body, ro)
	return
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToClusterTemplateCreateMap() (map[string]interface{}, error)
}

// CreateOpts satisfies the CreateOptsBuilder interface
type CreateOpts struct {
	Name              string            `json:"name" required:"true"`
	COE               string            `json:"coe" required:"true"`
	ImageID           string            `json:"image_id" required:"true"`
	ExternalNetworkID string            `json:"external_network_id" required:"true"`
	DiscoveryURL      string            `json:"discovery_url,omitempty"`
	DockerVolumeSize  int               `json:"docker_volume_size" required:"true"`
	NetworkDriver     string            `json:"network_driver,omiempty"`
	FlavorID          string            `json:"flavor_id,omiempty"`
	MasterFlavorID    string            `json:"master_flavor_id,omiempty"`
	KeypairID         string            `json:"keypair_id,omiempty"`
	Labels            map[string]string `json:"labels,omiempty"`
}

// ToClusterTemplateCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToClusterTemplateCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create a clustertemplate.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterTemplateCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, ro)
	return
}
