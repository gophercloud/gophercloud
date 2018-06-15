package clustertemplates

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a cluster-template resource.
func (r commonResult) Extract() (*ClusterTemplate, error) {
	var s *ClusterTemplate
	err := r.ExtractInto(&s)
	return s, err
}

// Represents a template for a Cluster Template
type ClusterTemplate struct {
	InsecureRegistry    string             `json:"insecure_registry"`
	Links               []gophercloud.Link `json:"links"`
	HTTPProxy           string             `json:"http_proxy"`
	UpdatedAt           time.Time          `json:"-"`
	FloatingIPEnabled   bool               `json:"floating_ip_enabled"`
	FixedSubnet         string             `json:"fixed_subnet"`
	MasterFlavorID      string             `json:"master_flavor_id"`
	UUID                string             `json:"uuid"`
	NoProxy             string             `json:"no_proxy"`
	HTTPSProxy          string             `json:"https_proxy"`
	TLSDisabled         bool               `json:"tls_disabled"`
	KeyPairID           string             `json:"keypair_id"`
	Public              bool               `json:"public"`
	Labels              map[string]string  `json:"labels"`
	DockerVolumeSize    int                `json:"docker_volume_size"`
	ServerType          string             `json:"server_type"`
	ExternalNetworkID   string             `json:"external_network_id"`
	ClusterDistro       string             `json:"cluster_distro"`
	ImageID             string             `json:"image_id"`
	VolumeDriver        string             `json:"volume_driver"`
	RegistryEnabled     bool               `json:"registry_enabled"`
	DockerStorageDriver string             `json:"docker_storage_driver"`
	APIServerPort       string             `json:"apiserver_port"`
	Name                string             `json:"name"`
	CreatedAt           time.Time          `json:"-"`
	NetworkDriver       string             `json:"network_driver"`
	FixedNetwork        string             `json:"fixed_network"`
	COE                 string             `json:"coe"`
	FlavorID            string             `json:"flavor_id"`
	MasterLBEnabled     bool               `json:"master_lb_enabled"`
	DNSNameServer       string             `json:"dns_nameserver"`
}

func (r *ClusterTemplate) UnmarshalJSON(b []byte) error {
	type tmp ClusterTemplate
	var s struct {
		tmp
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ClusterTemplate(s.tmp)

	if s.CreatedAt != "" {
		r.CreatedAt, err = time.Parse(time.RFC3339, s.CreatedAt)
		if err != nil {
			return err
		}
	}

	if s.UpdatedAt != "" {
		r.UpdatedAt, err = time.Parse(time.RFC3339, s.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}
