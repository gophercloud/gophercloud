package clustertemplates

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/common"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToClusterTemplateListQuery() (string, error)
}

// ListOpts allows the sorting of paginated collections through
// the API. SortKey allows you to sort by a particular clustertemplate attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for pagination.
type ListOpts struct {
	Marker  string `q:"marker"`
	Limit   int    `q:"limit"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
}

// ToClusterTemplateListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterTemplateListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// clustertemplates. It accepts a ListOptsBuilder, which allows you to sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToClusterTemplateListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterTemplatePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific clustertemplate based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Get(getURL(c, id), &r.Body, ro)
	return
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToClusterTemplateCreateMap() (map[string]interface{}, error)
}

// CreateOpts satisfies the CreateOptsBuilder interface
type CreateOpts struct {
	Name              string            `json:"name" required:"true"`
	COE               string            `json:"coe" required:"true"`
	ImageID           string            `json:"image_id" required:"true"`
	ExternalNetworkID string            `json:"external_network_id" required:"true"`
	DiscoveryURL      string            `json:"discovery_url,omitempty"`
	DockerVolumeSize  int               `json:"docker_volume_size" required:"true"`
	NetworkDriver     string            `json:"network_driver,omiempty"`
	FlavorID          string            `json:"flavor_id,omiempty"`
	MasterFlavorID    string            `json:"master_flavor_id,omiempty"`
	KeypairID         string            `json:"keypair_id,omiempty"`
	Labels            map[string]string `json:"labels,omiempty"`
}

// ToClusterTemplateCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToClusterTemplateCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create a clustertemplate.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterTemplateCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, ro)
	return
}

// Delete accepts a unique ID and deletes the cluster template associated with it.
func Delete(c *gophercloud.ServiceClient, clusterTemplateID string) (r DeleteResult) {
	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Delete(deleteURL(c, clusterTemplateID), ro)
	return
}

// Update implements cluster template updated request.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterTemplateUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}
	_, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
