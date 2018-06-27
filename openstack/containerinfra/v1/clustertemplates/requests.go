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
	Name                string            `json:"name,omitempty"`
	COE                 string            `json:"coe" required:"true"`
	ImageID             string            `json:"image_id" required:"true"`
	FlavorID            string            `json:"flavor_id,omitempty"`
	MasterFlavorID      string            `json:"master_flavor_id,omitempty"`
	DNSNameServer       string            `json:"dns_nameserver,omitempty"`
	KeyPairID           string            `json:"keypair_id,omitempty"`
	ExternalNetworkID   string            `json:"external_network_id,omitempty"`
	FixedNetwork        string            `json:"fixed_network,omitempty"`
	FixedSubnet         string            `json:"fixed_subnet,omitempty"`
	NetworkDriver       string            `json:"network_driver,omitempty"`
	APIServerPort       int               `json:"apiserver_port,omitempty"`
	DockerVolumeSize    int               `json:"docker_volume_size,omitempty"`
	ClusterDistro       string            `json:"cluster_distro,omitempty"`
	HTTPProxy           string            `json:"http_proxy,omitempty"`
	HTTPSProxy          string            `json:"https_proxy,omitempty"`
	NoProxy             string            `json:"no_proxy,omitempty"`
	VolumeDriver        string            `json:"volume_driver,omitempty"`
	RegistryEnabled     bool              `json:"registry_enabled,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
	TLSDisabled         bool              `json:"tls_disabled,omitempty"`
	Public              bool              `json:"public,omitempty"`
	ServerType          string            `json:"server_type,omitempty"`
	InsecureRegistry    string            `json:"insecure_registry,omitempty"`
	DockerStorageDriver string            `json:"docker_storage_driver,omitempty"`
	MasterLBEnabled     bool              `json:"master_lb_enabled,omitempty"`
	FloatingIPEnabled   bool              `json:"floating_ip_enabled,omitempty"`
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
