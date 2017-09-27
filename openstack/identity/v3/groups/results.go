package groups

import (
	"github.com/gophercloud/gophercloud/openstack/identity/v3/groupresults"
	"github.com/gophercloud/gophercloud/pagination"
)

// ExtractUsers returns a slice of Groups contained in a single page of results.
func ExtractGroups(r pagination.Page) ([]groupresults.Group, error) {
	return groupresults.ExtractGroups(r)
}
