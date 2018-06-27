package clustertemplates

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	Labels              map[string]string `json:"labels,omitempty"`
	FixedSubnet         string            `json:"fixed_subnet,omitempty"`
	MasterFlavorID      string            `json:"master_flavor_id,omitempty"`
	NoProxy             string            `json:"no_proxy,omitempty"`
	HTTPSProxy          string            `json:"https_proxy,omitempty"`
	HTTPProxy           string            `json:"http_proxy,omitempty"`
	TLSDisabled         bool              `json:"tls_disabled,omitempty"`
	KeyPairID           string            `json:"keypair_id,omitempty"`
	Public              bool              `json:"public,omitempty"`
	DockerVolumeSize    int               `json:"docker_volume_size"`
	ServerType          string            `json:"server_type,omitempty"`
	ExternalNetworkID   string            `json:"external_network_id,omitempty"`
	ImageID             string            `json:"image_id" required:"true"`
	VolumeDriver        string            `json:"volume_driver,omitempty"`
	RegistryEnabled     bool              `json:"registry_enabled,omitempty"`
	DockerStorageDriver string            `json:"docker_storage_driver,omitempty"`
	Name                string            `json:"name,omitempty"`
	NetworkDriver       string            `json:"network_driver,omitempty"`
	FixedNetwork        string            `json:"fixed_network,omitempty"`
	COE                 string            `json:"coe" required:"true"`
	FlavorID            string            `json:"flavor_id,omitempty"`
	MasterLBEnabled     bool              `json:"master_lb_enabled,omitempty"`
	DNSNameServer       string            `json:"dns_nameserver,omitempty"`
}

// ToClusterCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create requests the creation of a new cluster.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})

	if r.Err == nil {
		r.Header = result.Header
	}

	return
}
