package clustertemplates

import (
	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// CreateResult temporarily contains the response from a Create call.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a clustertemplate resource.
func (r commonResult) Extract() (*ClusterTemplate, error) {
	var s *ClusterTemplate
	err := r.ExtractInto(&s)
	return s, err
}

// Represents a template for a Cluster
type ClusterTemplate struct {
	// UUID for the clustertemplate
	ID string `json:"uuid"`

	// Human-readable name for the clustertemplate. Might not be unique.
	Name string `json:"name"`

	// The type of container orchestration engine used by the cluster.
	COE string `json:"coe"`

	// The underlying type of the host nodes, such as lxc or vm
	ServerType string `json:"server_type"`

	// The flavor used by nodes in the cluster.
	FlavorID string `json:"flavor_id"`

	// The image used by nodes in the cluster.
	ImageID string `json:"image_id"`

	// Specifies if the cluster should use TLS certificates.
	TLSDisabled bool `json:"tls_disabled"`

	// The KeyPair used by the cluster.
	KeyPairID string `json:"keypair_id"`

	// The networt driver used by the cluster
	NetworkDriver string `json:"network_driver"`
}
