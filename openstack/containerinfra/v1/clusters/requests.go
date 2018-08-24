package clusters

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
	ClusterTemplateID string            `json:"cluster_template_id" required:"true"`
	CreateTimeout     *int              `json:"create_timeout"`
	DiscoveryURL      string            `json:"discovery_url,omitempty"`
	DockerVolumeSize  *int              `json:"docker_volume_size,omitempty"`
	FlavorID          string            `json:"flavor_id,omitempty"`
	Keypair           string            `json:"keypair,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	MasterCount       *int              `json:"master_count,omitempty"`
	MasterFlavorID    string            `json:"master_flavor_id,omitempty"`
	Name              string            `json:"name"`
	NodeCount         *int              `json:"node_count,omitempty"`
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
		OkCodes: []int{202},
	})

	if r.Err == nil {
		r.Header = result.Header
	}

	return
}
