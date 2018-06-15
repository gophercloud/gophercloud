package clustertemplates

import (
	"fmt"
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

// DeleteResult is the result from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract is a function that accepts a result and extracts a cluster-template resource.
func (r commonResult) Extract() (*ClusterTemplate, error) {
	var s *ClusterTemplate
	err := r.ExtractInto(&s)
	return s, err
}

// Represents a template for a Cluster Template
type ClusterTemplate struct {
	APIServerPort       string             `json:"apiserver_port"`
	COE                 string             `json:"coe"`
	ClusterDistro       string             `json:"cluster_distro"`
	CreatedAt           time.Time          `json:"created_at"`
	DNSNameServer       string             `json:"dns_nameserver"`
	DockerStorageDriver string             `json:"docker_storage_driver"`
	DockerVolumeSize    int                `json:"docker_volume_size"`
	ExternalNetworkID   string             `json:"external_network_id"`
	FixedNetwork        string             `json:"fixed_network"`
	FixedSubnet         string             `json:"fixed_subnet"`
	FlavorID            string             `json:"flavor_id"`
	FloatingIPEnabled   bool               `json:"floating_ip_enabled"`
	HTTPProxy           string             `json:"http_proxy"`
	HTTPSProxy          string             `json:"https_proxy"`
	ImageID             string             `json:"image_id"`
	InsecureRegistry    string             `json:"insecure_registry"`
	KeyPairID           string             `json:"keypair_id"`
	Labels              map[string]string  `json:"labels"`
	Links               []gophercloud.Link `json:"links"`
	MasterFlavorID      string             `json:"master_flavor_id"`
	MasterLBEnabled     bool               `json:"master_lb_enabled"`
	Name                string             `json:"name"`
	NetworkDriver       string             `json:"network_driver"`
	NoProxy             string             `json:"no_proxy"`
	ProjectID           string             `json:"project_id"`
	Public              bool               `json:"public"`
	RegistryEnabled     bool               `json:"registry_enabled"`
	ServerType          string             `json:"server_type"`
	TLSDisabled         bool               `json:"tls_disabled"`
	UUID                string             `json:"uuid"`
	UpdatedAt           time.Time          `json:"updated_at"`
	UserID              string             `json:"user_id"`
	VolumeDriver        string             `json:"volume_driver"`
}

func (r DeleteResult) Extract() (string, error) {
	uuid := ""
	idKey := "X-Openstack-Request-Id"
	if len(r.Header[idKey]) > 0 {
		uuid = r.Header[idKey][0]
		if uuid == "" {
			return "", fmt.Errorf("No uuid value in response header %s", idKey)
		}
	} else {
		return "", fmt.Errorf("Missing [%s] in response header", idKey)
	}
	return uuid, r.ExtractErr()
}
