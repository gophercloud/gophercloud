package roles

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// RoleAssignmentsOpts allows you to query the RoleAssignments method.
type RoleAssignmentsOpts struct {
	GroupId        string `q:"group.id"`
	RoleId         string `q:"role.id"`
	ScopeDomainId  string `q:"scope.domain.id"`
	ScopeProjectId string `q:"scope.project.id"`
	UserId         string `q:"user.id"`
	Effective      bool   `q:"effective"`
}

// RoleAssignments enumerates the roles assigned to a specified resource.
func RoleAssignments(client *gophercloud.ServiceClient, opts RoleAssignmentsOpts) pagination.Pager {
	u := roleAssignmentsURL(client)
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	u += q.String()
	createPage := func(r pagination.PageResult) pagination.Page {
		return RoleAssignmentsPage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, u, createPage)
}
