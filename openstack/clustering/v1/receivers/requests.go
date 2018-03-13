package receivers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToReceiverCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	Name      string                 `json:"name" required:"true"`
	ClusterID string                 `json:"cluster_id,omitempty"`
	Type      string                 `json:"type" required:"true"`
	Action    string                 `json:"action,omitempty"`
	Actor     map[string]interface{} `json:"actor,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

// ToReceiverCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToReceiverCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "receiver")
}

// Create requests the creation of a new receiver.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToReceiverCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	return
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

	// Type to filter the response by the specified type property of the object
	Type string `q:"type,omitempty"`

	// ClusterID The name, short-ID or UUID of the cluster object
	ClusterID string `q:"cluster_id,omitempty"`

	// Action Filters the response by the action targeted by the receiver.
	Action string `q:"action,omitempty"`

	// User Filters the response by the user targeted by the receiver.
	User string `q:"user,omitempty"`
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToReceiverListQuery() (string, error)
}

// ToReceiverListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToReceiverListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail instructs OpenStack to provide a list of cluster.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToReceiverListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ReceiverPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details of a single receiver. Use ExtractReceiver to convert its
// result into a Receiver.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// UpdateOpts params
type UpdateOpts struct {
	Receiver map[string]interface{} `json:"-"`

	// Name is the name of the receiver.
	Name string `json:"name,omitempty"`

	Action string `json:"action,omitempty"`

	Params map[string]interface{} `json:"params,omitempty"`
}

// UpdateOptsBuilder params
type UpdateOptsBuilder interface {
	ToReceiverUpdateMap() (map[string]interface{}, error)
}

// ToReceiverUpdateMap constructs a request body from UpdateOpts.
func (opts UpdateOpts) ToReceiverUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "receiver")
}

// Update requests the update of a receiver.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToReceiverUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete deletes the specified receiver ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// NotifyReceiver Notifies message type receiver
func NotifyReceiver(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	result, err := client.Post(notifyURL(client, id), "", &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202, 204},
	})

	if err != nil {
		if err.Error() == "EOF" {
			// valid 204
			r.Err = nil
		} else {
			r.Err = err
		}
	} else {
		r.Header = result.Header
	}

	return
}
