package clustertemplates

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an Update operation. Call its Extract
// method to interpret it as a cluster template.
type UpdateResult struct {
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

// ClusterTemplatePage is the page returned by a pager when traversing over a
// collection of clustertemplates.
type ClusterTemplatePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of clustertemplates has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterTemplatePage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Next, nil
}

// IsEmpty checks whether a ClusterTemplatePage struct is empty.
func (r ClusterTemplatePage) IsEmpty() (bool, error) {
	is, err := ExtractClusterTemplates(r)
	return len(is) == 0, err
}

// ExtractClusterTemplates accepts a Page struct, specifically a ClusterTemplatePage struct,
// and extracts the elements into a slice of ClusterTemplate structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusterTemplates(r pagination.Page) ([]ClusterTemplate, error) {
	var s struct {
		ClusterTemplates []ClusterTemplate `json:"clustertemplates"`
	}
	err := (r.(ClusterTemplatePage)).ExtractInto(&s)
	return s.ClusterTemplates, err
}
