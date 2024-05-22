package replicas

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToReplicaCreateMap() (map[string]any, error)
}

// CreateOpts contains the options for create a Share Replica. This object is
// passed to replicas.Create function. For more information about these parameters,
// please refer to the Replica object, or the shared file systems API v2
// documentation.
type CreateOpts struct {
	// The UUID of the share from which to create a share replica.
	ShareID string `json:"share_id" required:"true"`
	// The UUID of the share network to which the share replica should
	// belong to.
	ShareNetworkID string `json:"share_network_id,omitempty"`
	// The availability zone of the share replica.
	AvailabilityZone string `json:"availability_zone,omitempty"`
	// One or more scheduler hints key and value pairs as a dictionary of
	// strings. Minimum supported microversion for SchedulerHints is 2.67.
	SchedulerHints map[string]string `json:"scheduler_hints,omitempty"`
}

// ToReplicaCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToReplicaCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "share_replica")
}

// Create will create a new Share Replica based on the values in CreateOpts. To extract
// the Replica object from the response, call the Extract method on the
// CreateResult.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToReplicaCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOpts holds options for listing Share Replicas. This object is passed to the
// replicas.List or replicas.ListDetail functions.
type ListOpts struct {
	// The UUID of the share.
	ShareID string `q:"share_id"`
	// Per page limit for share replicas
	Limit int `q:"limit"`
	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`
	// The ID of the last-seen item.
	Marker string `q:"marker"`
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToReplicaListQuery() (string, error)
}

// ToReplicaListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToReplicaListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns []Replica optionally limited by the conditions provided in ListOpts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToReplicaListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := ReplicaPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// ListDetail returns []Replica optionally limited by the conditions provided in ListOpts.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToReplicaListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := ReplicaPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// Delete will delete an existing Replica with the given UUID.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get will get a single share with given UUID
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListExportLocations will list replicaID's export locations.
// Minimum supported microversion for ListExportLocations is 2.47.
func ListExportLocations(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ListExportLocationsResult) {
	resp, err := client.Get(ctx, listExportLocationsURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetExportLocation will get replicaID's export location by an ID.
// Minimum supported microversion for GetExportLocation is 2.47.
func GetExportLocation(ctx context.Context, client *gophercloud.ServiceClient, replicaID string, id string) (r GetExportLocationResult) {
	resp, err := client.Get(ctx, getExportLocationURL(client, replicaID, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// PromoteOptsBuilder allows extensions to add additional parameters to the
// Promote request.
type PromoteOptsBuilder interface {
	ToReplicaPromoteMap() (map[string]any, error)
}

// PromoteOpts contains options for promoteing a Replica to active replica state.
// This object is passed to the replicas.Promote function.
type PromoteOpts struct {
	// The quiesce wait time in seconds used during replica promote.
	// Minimum supported microversion for QuiesceWaitTime is 2.75.
	QuiesceWaitTime int `json:"quiesce_wait_time,omitempty"`
}

// ToReplicaPromoteMap assembles a request body based on the contents of a
// PromoteOpts.
func (opts PromoteOpts) ToReplicaPromoteMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "promote")
}

// Promote will promote an existing Replica to active state. PromoteResult contains only the error.
// To extract it, call the ExtractErr method on the PromoteResult.
func Promote(ctx context.Context, client *gophercloud.ServiceClient, id string, opts PromoteOptsBuilder) (r PromoteResult) {
	b, err := opts.ToReplicaPromoteMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Resync a replica with its active mirror. ResyncResult contains only the error.
// To extract it, call the ExtractErr method on the ResyncResult.
func Resync(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ResyncResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"resync": nil}, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResetStatusOptsBuilder allows extensions to add additional parameters to the
// ResetStatus request.
type ResetStatusOptsBuilder interface {
	ToReplicaResetStatusMap() (map[string]any, error)
}

// ResetStatusOpts contain options for updating a Share Replica status. This object is passed
// to the replicas.ResetStatus function. Administrator only.
type ResetStatusOpts struct {
	// The status of a share replica. List of possible values: "available",
	// "error", "creating", "deleting" or "error_deleting".
	Status string `json:"status" required:"true"`
}

// ToReplicaResetStatusMap assembles a request body based on the contents of an
// ResetStatusOpts.
func (opts ResetStatusOpts) ToReplicaResetStatusMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "reset_status")
}

// ResetStatus will reset the Share Replica status with provided information.
// ResetStatusResult contains only the error. To extract it, call the ExtractErr
// method on the ResetStatusResult.
func ResetStatus(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ResetStatusOptsBuilder) (r ResetStatusResult) {
	b, err := opts.ToReplicaResetStatusMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ResetStateOptsBuilder allows extensions to add additional parameters to the
// ResetState request.
type ResetStateOptsBuilder interface {
	ToReplicaResetStateMap() (map[string]any, error)
}

// ResetStateOpts contain options for updating a Share Replica state. This object is passed
// to the replicas.ResetState function. Administrator only.
type ResetStateOpts struct {
	// The state of a share replica. List of possible values: "active",
	// "in_sync", "out_of_sync" or "error".
	State string `json:"replica_state" required:"true"`
}

// ToReplicaResetStateMap assembles a request body based on the contents of an
// ResetStateOpts.
func (opts ResetStateOpts) ToReplicaResetStateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "reset_replica_state")
}

// ResetState will reset the Share Replica state with provided information.
// ResetStateResult contains only the error. To extract it, call the ExtractErr
// method on the ResetStateResult.
func ResetState(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ResetStateOptsBuilder) (r ResetStateResult) {
	b, err := opts.ToReplicaResetStateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, actionURL(client, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ForceDelete force-deletes a Share Replica in any state. ForceDeleteResult
// contains only the error. To extract it, call the ExtractErr method on the
// ForceDeleteResult. Administrator only.
func ForceDelete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ForceDeleteResult) {
	resp, err := client.Post(ctx, actionURL(client, id), map[string]any{"force_delete": nil}, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
