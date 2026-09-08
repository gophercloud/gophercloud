package backupstrategies

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToBackupStrategyListQuery() (string, error)
}

// ListOpts represents options for listing backup strategies.
type ListOpts struct {
	// Return the strategy for a particular database instance.
	InstanceID string `q:"instance_id"`
	// Return strategies for a particular project. Admin only.
	ProjectID string `q:"project_id"`
}

// ToBackupStrategyListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToBackupStrategyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List lists backup strategies.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToBackupStrategyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return BackupStrategyPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder is the top-level interface for create backup strategy
// options.
type CreateOptsBuilder interface {
	ToBackupStrategyCreateMap() (map[string]any, error)
}

// CreateOpts represents options for creating a backup strategy.
type CreateOpts struct {
	// The database instance ID. When omitted, the strategy applies at project
	// scope.
	InstanceID string `json:"instance_id,omitempty"`
	// User-defined Swift container name.
	SwiftContainer string `json:"swift_container" required:"true"`
}

// ToBackupStrategyCreateMap converts a CreateOpts struct into a request body.
func (opts CreateOpts) ToBackupStrategyCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "backup_strategy")
}

// Create creates a backup strategy.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBackupStrategyCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteOptsBuilder allows extensions to add additional parameters to the
// Delete request.
type DeleteOptsBuilder interface {
	ToBackupStrategyDeleteQuery() (string, error)
}

// DeleteOpts represents options for deleting backup strategies.
type DeleteOpts struct {
	// Delete the strategy for a particular database instance.
	InstanceID string `q:"instance_id"`
	// Delete strategies for a particular project. Admin only.
	ProjectID string `q:"project_id"`
}

// ToBackupStrategyDeleteQuery formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToBackupStrategyDeleteQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Delete deletes a backup strategy.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, opts DeleteOptsBuilder) (r DeleteResult) {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToBackupStrategyDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
