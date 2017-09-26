package users

import (
	"github.com/gophercloud/gophercloud/openstack/identity/v3/userresults"
	"github.com/gophercloud/gophercloud/pagination"
)

// ExtractUsers returns a slice of Users contained in a single page of results.
func ExtractUsers(r pagination.Page) ([]userresults.User, error) {
	return userresults.ExtractUsers(r)
}
