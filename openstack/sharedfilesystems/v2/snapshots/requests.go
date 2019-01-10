package snapshots

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOpts holds options for listing Snapshots. It is passed to the
// snapshots.List function.
type ListOpts struct {
	// (Admin only). Defines whether to list the requested resources for all projects.
	AllTenants bool `q:"all_tenants"`
	// The snapshot name.
	Name string `q:"name"`
	// Filter  by a snapshot description.
	Description string `q:"description"`
	// Filters by a share from which the snapshot was created.
	ShareID string `q:"share_id"`
	// Filters by a snapshot size in GB.
	Size int `q:"size"`
	// Filters by a snapshot status.
	Status string `q:"status"`
	// The maximum number of snapshots to return.
	Limit int `q:"limit"`
	// The offset to define start point of snapshot or snapshot group listing.
	Offset int `q:"offset"`
	// The key to sort a list of snapshots.
	SortKey string `q:"sort_key"`
	// The direction to sort a list of snapshots.
	SortDir string `q:"sort_dir"`
	// The UUID of the project in which the snapshot was created. Useful with all_tenants parameter.
	ProjectID string `q:"project_id"`
	// The name pattern that can be used to filter snapshots, snapshot snapshots, snapshot networks or snapshot groups.
	NamePattern string `q:"name~"`
	// The description pattern that can be used to filter snapshots, snapshot snapshots, snapshot networks or snapshot groups.
	DescriptionPattern string `q:"description~"`
}

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToSnapshotListQuery() (string, error)
}

// ToSnapshotListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSnapshotListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail returns []Snapshot optionally limited by the conditions provided in ListOpts.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToSnapshotListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		p := SnapshotPage{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	})
}

// Get will get a single snapshot with given UUID
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}
