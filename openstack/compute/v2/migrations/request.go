package migrations

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"net/url"
	"time"
)

type ListOptsBuilder interface {
	ToMigrationsListQuery() (string, error)
}

type ListOpts struct {
	// The source/destination compute node of migration to filter
	Host *string `q:"host"`
	// The uuid of the instance that migration is operated on to filter
	InstanceID *string `q:"instance_uuid"`
	// The type of migration to filter. Valid values are: evacuation, live-migration, migration, resize
	MigrationType *string `q:"migration_type"`
	// The source compute node of migration to filter
	SourceCompute *string `q:"source_compute"`
	// The status of migration to filter
	Status *string `q:"status"`
	// Requests a page size of items
	Limit *int `q:"limit"`
	// The UUID of the last-seen migration
	Marker *string `q:"marker"`
	// Filters the response by a date and time stamp when the migration last changed
	ChangesSince *time.Time `q:"changes-since"`
	// Filters the response by a date and time stamp when the migration last changed
	ChangesBefore *time.Time `q:"changes-before"`
	// Filter the migrations by the given user ID
	UserID *string `q:"user_id"`
	// Filter the migrations by the given project ID
	ProjectID *string `q:"project_id"`
}

func (opts ListOpts) ToMigrationsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()

	if opts.ChangesSince != nil {
		params.Add("changes-since", opts.ChangesSince.Format(time.RFC3339))
	}

	if opts.ChangesBefore != nil {
		params.Add("changes-before", opts.ChangesBefore.Format(time.RFC3339))
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), nil
}

func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	reqUrl := listURL(client)
	if opts != nil {
		query, err := opts.ToMigrationsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		reqUrl += query
	}

	return pagination.NewPager(client, reqUrl, func(r pagination.PageResult) pagination.Page {
		return MigrationPage{pagination.SinglePageBase(r)}
	})
}
