package backups

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/db/v1/datastores"
)

type Backup struct {
	Description string
	ID          string
	InstanceID  string `json:"instance_id" mapstructure:"instance_id"`
	LocationRef string
	Name        string
	ParentID    string `json:"parent_id" mapstructure:"parent_id"`
	Size        float64
	Status      string
	Created     string
	Updated     string
	Datastore   datastores.DatastorePartial
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type commonResult struct {
	gophercloud.Result
}

func (r commonResult) Extract() (*Backup, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Backup Backup `mapstructure:"backup"`
	}

	err := mapstructure.Decode(r.Body, &response)
	return &response.Backup, err
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type BackupPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether an BackupPage struct is empty.
func (r BackupPage) IsEmpty() (bool, error) {
	is, err := ExtractBackups(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractBackups(page pagination.Page) ([]Backup, error) {
	casted := page.(BackupPage).Body

	var resp struct {
		Backups []Backup `mapstructure:"backups" json:"backups"`
	}

	err := mapstructure.Decode(casted, &resp)
	return resp.Backups, err
}
