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

type UUID struct {
	UUID string `json:"uuid"`
}

func (r CreateResult) Extract() (clusterID string, err error) {
	var s *UUID
	err = r.ExtractInto(&s)
	if err == nil {
		clusterID = s.UUID
	}
	return clusterID, err

}

type Cluster struct {
	APIAddress        string             `json:"api_address"`
	COEVersion        string             `json:"coe_version"`
	ClusterTemplateID string             `json:"cluster_template_id"`
	CreateTimeout     int                `json:"create_timeout"`
	CreatedAt         time.Time          `json:"created_at"`
	DiscoveryURL      string             `json:"discovery_url"`
	KeyPair           string             `json:"keypair"`
	Links             []gophercloud.Link `json:"links"`
	MasterAddresses   []string           `json:"master_addresses"`
	MasterCount       int                `json:"master_count"`
	Name              string             `json:"name"`
	NodeAddresses     []string           `json:"node_addresses"`
	NodeCount         int                `json:"node_count"`
	StackID           string             `json:"stack_id"`
	Status            string             `json:"status"`
	StatusReason      string             `json:"status_reason"`
	UUID              string             `json:"uuid"`
	UpdatedAt         time.Time          `json:"updated_at"`
}
