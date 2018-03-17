package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"net/http"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	// Name is the name of the cluster.
	Name string `json:"name" required:"true"`

	// When required is set to true, must contain non-zero value
	DesiredCapacity int `json:"desired_capacity" required:"false"`

	ProfileID string `json:"profile_id" required:"true"`

	MinSize int `json:"min_size,omitempty"`

	Timeout int `json:"timeout,omitempty"`

	MaxSize int `json:"max_size,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Config The structured config associated with the cluster.
	Config map[string]interface{} `json:"config,omitempty"`
}

// ToClusterCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "cluster")
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
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header

	return
}

// ActionOpts params
type ActionOpts struct {
	// Name is the name of the cluster.
	Name string `json:"name" required:"true"`

	DesiredCapacity int `json:"desired_capacity" required:"true"`

	ProfileID string `json:"profile_id" required:"true"`

	MinSize *int `json:"min_size"`

	Timeout *int `json:"timeout"`

	MaxSize *int `json:"max_size"`

	//Metadata *int `json:"metadata"`
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToClusterListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	// Limit instructs List to refrain from sending excessively large lists of clusters
	Limit int `q:"limit,omitempty"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Marker string `q:"marker,omitempty"`

	// Marker and Limit control paging. Marker instructs List where to start listing from.
	Sort string `q:"sort,omitempty"`

	// GlobalProject indicates whether to include resources for all projects or resources for the current project
	GlobalProject string `q:"global_project,omitempty"`

	// Name to filter the response by the specified name property of the object
	Name string `q:"name,omitempty"`

	// Status to filter the resource collection by the specified status property
	Status string `q:"status,omitempty"`
}

// ToClusterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail instructs OpenStack to provide a list of cluster.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToClusterListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOpts implements UpdateOpts
type UpdateOpts struct {
	// Profile A structured description of a profile object.
	Cluster map[string]interface{} `json:"-"`

	// Config The structured config associated with the cluster.
	Config string `json:"config,omitempty"`

	// Name is the name of the new profile.
	Name string `json:"name,omitempty"`

	ProfileID string `json:"profile_id,omitempty"`

	Timeout int `json:"timeout,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`

	ProfileOnly bool `json:"profile_only,omitempty"`
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToClusterUpdateMap() (map[string]interface{}, error)
}

// ToClusterUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "cluster")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Update implements profile updated request.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}

	var result *http.Response
	result, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	r.Header = result.Header
	return
}

// Delete deletes the specified cluster ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(deleteURL(client, id), nil)
	r.Header = result.Header
	return
}

// ToClusterActionMap constructs a request body from CreateOpts.
func (opts ActionOpts) ToClusterActionMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "action_name")
}

// ResizeOpts params
type ResizeOpts struct {
	// Name is the name of the cluster.
	AdjustmentType string `json:"adjustment_type,omitempty"`

	// This can be either float/int. When adjustment_type is CHANGE_IN_PERCENTAGE it is float, otherwise int
	Number float32 `json:"number,omitempty"`

	MinSize int `json:"min_size,omitempty"`

	MaxSize int `json:"max_size,omitempty"`

	MinStep int `json:"min_step,omitempty"`

	Strict bool `json:"strict,omitempty"`
}

// ToClusterResizeMap constructs a request body from ResizeOpts.
func (opts ResizeOpts) ToClusterResizeMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "resize")
}

// Resize CLI
func Resize(client *gophercloud.ServiceClient, id string, opts ResizeOpts) (r GetResult) {
	b, err := opts.ToClusterResizeMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Post(actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Get retrieves details of a single cluster. Use ExtractCluster to convert its
// result into a Cluster.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	r.Header = result.Header
	return
}

// ScaleOpts params
type ScaleOpts struct {
	Count int `json:"count,omitempty"`
}

// ToClusterScaleMap constructs a request body from ScaleInOpts.
func (opts ScaleOpts) ToClusterScaleMap(scaleAction string) (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, scaleAction)
}

// ScaleIn CLI
func ScaleIn(client *gophercloud.ServiceClient, id string, opts ScaleOpts) (r PostResult) {
	b, err := opts.ToClusterScaleMap("scale_in")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(scaleURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// ScaleOut CLI
func ScaleOut(client *gophercloud.ServiceClient, id string, opts ScaleOpts) (r PostResult) {
	b, err := opts.ToClusterScaleMap("scale_out")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(scaleURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// ToClusterNodeMap constructs a request body from ScaleInOpts.
func (opts NodeOpts) ToClusterNodeMap(nodeAction string) (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, nodeAction)
}

// NodeOpts params
type NodeOpts struct {
	Nodes                []string `json:"nodes" required:"true"`
	DestroyAfterDeletion bool     `json:"destroy_after_deletion"`
}

// Add nodes
func AddNodes(client *gophercloud.ServiceClient, id string, opts NodeOpts) (r PostResult) {
	b, err := opts.ToClusterNodeMap("add_nodes")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(nodeURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Remove nodes
func DeleteNodes(client *gophercloud.ServiceClient, id string, opts NodeOpts) (r PostResult) {
	b, err := opts.ToClusterNodeMap("del_nodes")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(nodeURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Replace nodes
func ReplaceNodes(client *gophercloud.ServiceClient, id string, opts NodeOpts) (r PostResult) {
	b, err := opts.ToClusterNodeMap("replace_nodes")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(nodeURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// ToClusterPolicyMap constructs a request body from ScaleInOpts.
func (opts PolicyOpts) ToClusterPolicyMap(policyAction string) (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, policyAction)
}

// PolicyOpts params
type PolicyOpts struct {
	PolicyID string `json:"policy_id" required:"true"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

// Attach Policy
func AttachPolicy(client *gophercloud.ServiceClient, id string, opts PolicyOpts) (r PostResult) {
	b, err := opts.ToClusterPolicyMap("policy_attach")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(policyURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Detach Policy
func DetachPolicy(client *gophercloud.ServiceClient, id string, opts PolicyOpts) (r PostResult) {
	b, err := opts.ToClusterPolicyMap("policy_detach")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(policyURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Update Policy
func UpdatePolicy(client *gophercloud.ServiceClient, id string, opts PolicyOpts) (r PostResult) {
	b, err := opts.ToClusterPolicyMap("policy_update")
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(policyURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Collect Attributes

// Recover cluster

// Perform operation

// ToClusterHealthCheckMap constructs a request body from HealthCheckOpts.
func (opts HealthCheckOpts) ToClusterHealthCheckMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// PolicyOpts params
type HealthCheckOpts struct {
	PolicyID string `json:"policy_id" required:"true"`
	Enabled  bool   `json:"enabled"`
}

// Check a Cluster's Health Status
func (opts HealthCheckOpts) CheckHealth(client *gophercloud.ServiceClient, id string) (r PostResult) {
	b, err := opts.ToClusterHealthCheckMap()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(policyURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// ToClusterHealthCheckMap constructs a request body from HealthCheckOpts.
func (opts RecoverHealthOpts) ToClusterRecoverHealthMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "recover")
}

// PolicyOpts params
type RecoverHealthOpts struct {
	Operation     string `json:"operation" required:"true"`
	Check         bool   `json:"check"`
	CheckCapacity bool   `json:"check_capacity"`
}

// Recover a Cluster to a Healthy Status
func RecoverHealth(client *gophercloud.ServiceClient, id string, opts RecoverHealthOpts) (r PostResult) {
	b, err := opts.ToClusterRecoverHealthMap()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(healthURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	r.Header = result.Header
	return
}

// Perform an Operation on a Cluster

// ToClusterCompleteLifecycleMap constructs a request body from CompleteLifecycleOpts.
func (opts CompleteLifecycleOpts) ToClusterCompleteLifecycleMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "complete_lifecycle")
}

// Complete lifecycle params
type CompleteLifecycleOpts struct {
	LifecycleActionTokenID string `json:"lifecycle_action_token" required:"true"`
}

// Complete lifecycle
func CompleteLifecycle(client *gophercloud.ServiceClient, id string, opts CompleteLifecycleOpts) (r PostResult) {
	b, err := opts.ToClusterCompleteLifecycleMap()
	if err != nil {
		r.Err = err
		return
	}

	client.Microversion = "1.9"
	var result *http.Response
	result, r.Err = client.Post(actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	client.Microversion = ""
	r.Header = result.Header

	return
}
