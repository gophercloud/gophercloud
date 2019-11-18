package backup

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder builds create options
type CreateOptsBuilder interface {
	ToBackupCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the struct responsible for configuring a new database backup.
type CreateOpts struct {
	// Name of the instance backup to create. The length of the name is limited to
	// 255 characters and any characters are permitted. Optional.
	Name        string `json:"name" required:"true"`
	// The ID of the instance to create backup for.
	Instance    string `json:"instance" required:"true"`
	// ID of the parent backup to perform an incremental backup from.
	ParentID    string `json:"parent_id,omitempty"`
	// Create an incremental backup based on the last full backup by setting this
	// parameter to 1 or 0. It will create a full backup if no existing backup found.
	Incremental int    `json:"incremental"`
	// An optional description for the backup.
	Description string `json:"description,omitempty"`
	// The account ID of the owner of the instance.
	AccountID   string `json:"account_id,omitempty"`
}

// ToMap is a helper function to convert individual DB create opt structures
// into sub-maps.
func (opts CreateOpts) ToMap() (map[string]interface{}, error) {
	if len(opts.Name) > 64 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "backups.CreateOpts.Name"
		err.Value = opts.Name
		err.Info = "Must be less than 64 chars long"
		return nil, err
	}
	return gophercloud.BuildRequestBody(opts, "")
}

// ToInstanceCreateMap will render a JSON map.
func (opts CreateOpts) ToBackupCreateMap() (map[string]interface{}, error) {
	BackupMap, err := opts.ToMap()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"backup": BackupMap}, nil
}

// Create will create a new database backup within the specified instance. If the
// specified instance does not exist, a 404 error will be returned.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBackupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(baseURL(client), &b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{202}})
	return
}

// List retrieves the status and information for all database backup.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, baseURL(client), func(r pagination.PageResult) pagination.Page {
		return BackupsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// List all database backups for the specified instance.
func ListByInstance(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {
	return pagination.NewPager(client,
		instanceBackupURL(client, instanceID),
		func(r pagination.PageResult) pagination.Page {
			return BackupsPage{pagination.LinkedPageBase{PageResult: r}}
		})
}

// Get detailes of a backup.
func Get(client *gophercloud.ServiceClient, backupID string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, backupID), &r.Body, nil)
	return
}

// Delete will permanently delete the database backup
func Delete(client *gophercloud.ServiceClient, backupID string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, backupID), nil)
	return
}
