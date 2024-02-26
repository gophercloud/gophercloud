package osinherit

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
)

// AssignOpts provides options to assign an inherited role
type AssignOpts struct {
	// UserID is the ID of a user to assign an inherited role
	// Note: exactly one of UserID or GroupID must be provided
	UserID string `xor:"GroupID"`

	// GroupID is the ID of a group to assign an inherited role
	// Note: exactly one of UserID or GroupID must be provided
	GroupID string `xor:"UserID"`

	// ProjectID is the ID of a project to assign an inherited role on
	// Note: exactly one of ProjectID or DomainID must be provided
	ProjectID string `xor:"DomainID"`

	// DomainID is the ID of a domain to assign an inherited role on
	// Note: exactly one of ProjectID or DomainID must be provided
	DomainID string `xor:"ProjectID"`
}

// ValidateOpts provides options to which role to validate
type ValidateOpts struct {
	// UserID is the ID of a user to validate an inherited role
	// Note: exactly one of UserID or GroupID must be provided
	UserID string `xor:"GroupID"`

	// GroupID is the ID of a group to validate an inherited role
	// Note: exactly one of UserID or GroupID must be provided
	GroupID string `xor:"UserID"`

	// ProjectID is the ID of a project to validate an inherited role on
	// Note: exactly one of ProjectID or DomainID must be provided
	ProjectID string `xor:"DomainID"`

	// DomainID is the ID of a domain to validate an inherited role on
	// Note: exactly one of ProjectID or DomainID must be provided
	DomainID string `xor:"ProjectID"`
}

// UnassignOpts provides options to unassign an inherited role
type UnassignOpts struct {
	// UserID is the ID of a user to unassign an inherited role
	// Note: exactly one of UserID or GroupID must be provided
	UserID string `xor:"GroupID"`

	// GroupID is the ID of a group to unassign an inherited role
	// Note: exactly one of UserID or GroupID must be provided
	GroupID string `xor:"UserID"`

	// ProjectID is the ID of a project to assign an inherited role on
	// Note: exactly one of ProjectID or DomainID must be provided
	ProjectID string `xor:"DomainID"`

	// DomainID is the ID of a domain to assign an inherited role on
	// Note: exactly one of ProjectID or DomainID must be provided
	DomainID string `xor:"ProjectID"`
}

// Assign is the operation responsible for assigning an inherited role
// to a user/group on a project/domain.
func Assign(ctx context.Context, client *gophercloud.ServiceClient, roleID string, opts AssignOpts) (r AssignmentResult) {
	// Check xor conditions
	_, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	// Get corresponding URL
	var targetID string
	var targetType string
	if opts.ProjectID != "" {
		targetID = opts.ProjectID
		targetType = "projects"
	} else {
		targetID = opts.DomainID
		targetType = "domains"
	}

	var actorID string
	var actorType string
	if opts.UserID != "" {
		actorID = opts.UserID
		actorType = "users"
	} else {
		actorID = opts.GroupID
		actorType = "groups"
	}

	resp, err := client.Put(ctx, assignURL(client, targetType, targetID, actorType, actorID, roleID), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Validate is the operation responsible for validating an inherited role
// of a user/group on a project/domain.
func Validate(ctx context.Context, client *gophercloud.ServiceClient, roleID string, opts ValidateOpts) (r ValidateResult) {
	// Check xor conditions
	_, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	// Get corresponding URL
	var targetID string
	var targetType string
	if opts.ProjectID != "" {
		targetID = opts.ProjectID
		targetType = "projects"
	} else {
		targetID = opts.DomainID
		targetType = "domains"
	}

	var actorID string
	var actorType string
	if opts.UserID != "" {
		actorID = opts.UserID
		actorType = "users"
	} else {
		actorID = opts.GroupID
		actorType = "groups"
	}

	resp, err := client.Head(ctx, assignURL(client, targetType, targetID, actorType, actorID, roleID), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Unassign is the operation responsible for unassigning an inherited
// role to a user/group on a project/domain.
func Unassign(ctx context.Context, client *gophercloud.ServiceClient, roleID string, opts UnassignOpts) (r UnassignmentResult) {
	// Check xor conditions
	_, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	// Get corresponding URL
	var targetID string
	var targetType string
	if opts.ProjectID != "" {
		targetID = opts.ProjectID
		targetType = "projects"
	} else {
		targetID = opts.DomainID
		targetType = "domains"
	}

	var actorID string
	var actorType string
	if opts.UserID != "" {
		actorID = opts.UserID
		actorType = "users"
	} else {
		actorID = opts.GroupID
		actorType = "groups"
	}

	resp, err := client.Delete(ctx, assignURL(client, targetType, targetID, actorType, actorID, roleID), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
