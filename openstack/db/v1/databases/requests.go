package databases

import (
	"fmt"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	ToDBCreateMap() (map[string]interface{}, error)
}

// DatabaseOpts is the struct responsible for configuring a database; often in
// the context of an instance.
type CreateOpts struct {
	// Specifies the name of the database. Optional.
	Name string

	// Set of symbols and encodings. Optional; the default character set is utf8.
	CharSet string

	// Set of rules for comparing characters in a character set. Optional; the
	// default value for collate is utf8_general_ci.
	Collate string
}

func (opts CreateOpts) ToMap() (map[string]string, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf("Name is a required field")
	}
	if len(opts.Name) > 64 {
		return nil, fmt.Errorf("Name must be less than 64 chars long")
	}

	db := map[string]string{"name": opts.Name}

	if opts.CharSet != "" {
		db["character_set"] = opts.CharSet
	}
	if opts.Collate != "" {
		db["collate"] = opts.Collate
	}
	return db, nil
}

type BatchCreateOpts []CreateOpts

func (opts BatchCreateOpts) ToDBCreateMap() (map[string]interface{}, error) {
	var dbs []map[string]string
	for _, db := range opts {
		dbMap, err := db.ToMap()
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, dbMap)
	}
	return map[string]interface{}{"databases": dbs}, nil
}

func Create(client *gophercloud.ServiceClient, instanceID string, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToDBCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Request("POST", baseURL(client, instanceID), gophercloud.RequestOpts{
		JSONBody:     &reqBody,
		JSONResponse: &res.Body,
		OkCodes:      []int{202},
	})

	return res
}

func List(client *gophercloud.ServiceClient, instanceID string) pagination.Pager {
	createPageFn := func(r pagination.PageResult) pagination.Page {
		return DBPage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, baseURL(client, instanceID), createPageFn)
}

func Delete(client *gophercloud.ServiceClient, instanceID, dbName string) DeleteResult {
	var res DeleteResult

	_, res.Err = client.Request("DELETE", dbURL(client, instanceID, dbName), gophercloud.RequestOpts{
		JSONBody: &res.Body,
		OkCodes:  []int{202},
	})

	return res
}
