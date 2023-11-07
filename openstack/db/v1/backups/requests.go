package backups

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder is the top-level interface for creating JSON maps.
type CreateOptsBuilder interface {
	ToBackupCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the struct responsible for configuring a new user; often in the
// context of an instance.
type CreateOpts struct {
	// Specifies a name for the user. Valid names can be composed
	// of the following characters: letters (either case); numbers; these
	// characters '@', '?', '#', ' ' but NEVER beginning a name string; '_' is
	// permitted anywhere. Prohibited characters that are forbidden include:
	// single quotes, double quotes, back quotes, semicolons, commas, backslashes,
	// and forward slashes. Spaces at the front or end of a user name are also
	// not permitted.
	Name string `json:"name" required:"true"`
	// Specifies the ID of the instance to create backup for.
	Instance string `json:"instance" required:"true"`
	// Specifies the ID of the parent backup to perform an incremental backup from.
	// Optional
	ParentId string `json:"parentId,omitempty"`
	// Create an incremental backup based on the last full backup or not.
	// The value must be 0 or 1. Optional.
	Incremental int `json:"incremental,omitempty"`
	// Specifies a description for the backup. Optional.
	Description string `json:"description,omitempty"`
	// Specifies the swift container name. When creating backups, the swift
	// container is created automatically if does not exist.
	// The backup data is stored as object in the Swift container. Optional.
	SwiftContainer string `json:"swiftContainer,omitempty"`
}

// ToBackupCreateMap will generate a JSON map.
func (opts CreateOpts) ToBackupCreateMap() (map[string]interface{}, error) {
	if len(opts.Name) > 64 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "backups.CreateOpts.Name"
		err.Value = opts.Name
		err.Info = "Must be less than 64 chars long"
		return nil, err
	}
	if opts.Instance == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "backups.CreateOpts.Instance"}
	}
	if opts.Incremental > 1 || opts.Incremental < 0 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "backups.CreateOpts.Incremental"
		err.Value = opts.Incremental
		err.Info = "Incremental parameter must be 1 or 0"
		return nil, err
	}
	backup := map[string]interface{}{
		"name":     opts.Name,
		"instance": opts.Instance,
	}
	backup["incremental"] = opts.Incremental
	if opts.ParentId != "" {
		backup["parent_id"] = opts.ParentId
	}
	if opts.Description != "" {
		backup["description"] = opts.Description
	}
	if opts.SwiftContainer != "" {
		backup["swift_container"] = opts.SwiftContainer
	}
	return map[string]interface{}{"backup": backup}, nil
}

// Create asynchronously provisions a new backup for the specified database
// instance based on the configuration defined in CreateOpts.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBackupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// List will list all database backups information for a project.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, baseURL(client), func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get will show detailes of a backup for a project.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(resourceURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete will permanently delete a database backup.
// All the child backups are deleted automatically when a parent backup is deleted.
func Delete(client *gophercloud.ServiceClient, backupId string) (r DeleteResult) {
	resp, err := client.Delete(resourceURL(client, backupId), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOfInstance will get all the backups of an instance
func ListOfInstance(client *gophercloud.ServiceClient, instanceId string) pagination.Pager {
	return pagination.NewPager(client, listURL(client, instanceId), func(r pagination.PageResult) pagination.Page {
		return BackupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
