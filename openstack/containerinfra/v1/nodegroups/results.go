package nodegroups

import (
	"time"

	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*NodeGroup, error) {
	var s NodeGroup
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult is the response from a Get request.
// Use the Extract method to retrieve the NodeGroup itself.
type GetResult struct {
	commonResult
}

// NodeGroup is the API representation of a Magnum node group.
type NodeGroup struct {
	ID               int                `json:"id"`
	UUID             string             `json:"uuid"`
	Name             string             `json:"name"`
	ClusterID        string             `json:"cluster_id"`
	ProjectID        string             `json:"project_id"`
	DockerVolumeSize *int               `json:"docker_volume_size"`
	Labels           map[string]string  `json:"labels"`
	Links            []gophercloud.Link `json:"links"`
	FlavorID         string             `json:"flavor_id"`
	ImageID          string             `json:"image_id"`
	NodeAddresses    []string           `json:"node_addresses"`
	NodeCount        int                `json:"node_count"`
	Role             string             `json:"role"`
	MinNodeCount     int                `json:"min_node_count"`
	MaxNodeCount     *int               `json:"max_node_count"`
	IsDefault        bool               `json:"is_default"`
	StackID          string             `json:"stack_id"`
	Status           string             `json:"status"`
	StatusReason     string             `json:"status_reason"`
	Version          string             `json:"version"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}
