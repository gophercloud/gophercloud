package nodegroups

import (
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Get makes a request to the Magnum API to retrieve a node group
// with the given ID/name belonging to the given cluster.
// Use the Extract method of the returned GetResult to extract the
// node group from the result.
func Get(client *gophercloud.ServiceClient, clusterID, nodeGroupID string) (r GetResult) {
	var response *http.Response
	response, r.Err = client.Get(getURL(client, clusterID, nodeGroupID), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	if r.Err == nil {
		r.Header = response.Header
	}
	return
}

type ListOptsBuilder interface {
	ToNodeGroupsListQuery() (string, error)
}

// ListOpts is used to filter and sort the node groups of a cluster
// when using List.
type ListOpts struct {
	// Pagination marker for large data sets. (UUID field from node group).
	Marker int `q:"marker"`
	// Maximum number of resources to return in a single page.
	Limit int `q:"limit"`
	// Column to sort results by. Default: id.
	SortKey string `q:"sort_key"`
	// Direction to sort. "asc" or "desc". Default: asc.
	SortDir string `q:"sort_dir"`
	// List all nodegroups with the specified role.
	Role string `q:"role"`
}

func (opts ListOpts) ToNodeGroupsListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request to the Magnum API to retrieve node groups
// belonging to the given cluster. The request can be modified to
// filter or sort the list using the options available in ListOpts.
//
// Use the AllPages method of the returned Pager to ensure that
// all node groups are returned (for example when using the Limit
// option to limit the number of node groups returned per page).
//
// Not all node group fields are returned in a list request.
// Only the fields UUID, Name, FlavorID, ImageID,
// NodeCount, Role, IsDefault, Status and StackID
// are returned, all other fields are omitted
// and will have their zero value when extracted.
func List(client *gophercloud.ServiceClient, clusterID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, clusterID)
	if opts != nil {
		query, err := opts.ToNodeGroupsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NodeGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
