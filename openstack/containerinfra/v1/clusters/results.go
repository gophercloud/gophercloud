package clusters

import (
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

func (r CreateResult) Extract() (clusterID string, err error) {
	var s struct {
		UUID string
	}
	err = r.ExtractInto(&s)
	return s.UUID, err
}

type Cluster struct {
	APIAddress        string             `json:"api_address"`
	COEVersion        string             `json:"coe_version"`
	ClusterTemplateID string             `json:"cluster_template_id"`
	ContainerVersion  string             `json:"container_version"`
	CreateTimeout     int                `json:"create_timeout"`
	CreatedAt         time.Time          `json:"created_at"`
	DiscoveryURL      string             `json:"discovery_url"`
	DockerVolumeSize  int                `json:"docker_volume_size"`
	Faults            map[string]string  `json:"faults"`
	FlavorID          string             `json:"flavor_id"`
	KeyPair           string             `json:"keypair"`
	Labels            map[string]string  `json:"labels"`
	Links             []gophercloud.Link `json:"links"`
	MasterFlavorID    string             `json:"master_flavor_id"`
	MasterAddresses   []string           `json:"master_addresses"`
	MasterCount       int                `json:"master_count"`
	Name              string             `json:"name"`
	NodeAddresses     []string           `json:"node_addresses"`
	NodeCount         int                `json:"node_count"`
	ProjectID         string             `json:"project_id"`
	StackID           string             `json:"stack_id"`
	Status            string             `json:"status"`
	StatusReason      string             `json:"status_reason"`
	UUID              string             `json:"uuid"`
	UpdatedAt         time.Time          `json:"updated_at"`
	UserID            string             `json:"user_id"`
}
