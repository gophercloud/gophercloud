package snapshots

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/racker/perigee"
)

type CreateOpts struct {
	Description string
	Force       bool
	Metadata    map[string]interface{}
	Name        string
	VolumeID    string
}

func Create(client *gophercloud.ServiceClient, opts *CreateOpts) CreateResult {
	type snapshot struct {
		Description *string                `json:"display_description,omitempty"`
		Force       bool                   `json:"force,omitempty"`
		Metadata    map[string]interface{} `json:"metadata,omitempty"`
		Name        *string                `json:"display_name,omitempty"`
		VolumeID    *string                `json:"volume_id,omitempty"`
	}

	type request struct {
		Snapshot snapshot `json:"snapshot"`
	}

	reqBody := request{
		Snapshot: snapshot{},
	}

	reqBody.Snapshot.Description = utils.MaybeString(opts.Description)
	reqBody.Snapshot.Name = utils.MaybeString(opts.Name)
	reqBody.Snapshot.VolumeID = utils.MaybeString(opts.VolumeID)

	reqBody.Snapshot.Force = opts.Force

	var res CreateResult
	_, res.Err = perigee.Request("POST", createURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200, 201},
		ReqBody:     &reqBody,
		Results:     &res.Resp,
	})
	return res
}

func Delete(client *gophercloud.ServiceClient, id string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(client, id), perigee.Options{
		Results:     &res.Resp,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
	})
	return res
}

type ListOpts struct {
	Name     string `q:"display_name"`
	Status   string `q:"status"`
	VolumeID string `q:"volume_id"`
}

func List(client *gophercloud.ServiceClient, opts *ListOpts) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := gophercloud.BuildQueryString(opts)
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query.String()
	}

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListResult{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(client, url, createPage)
}

type UpdateOpts struct {
	Description string
	Name        string
}

func Update(client *gophercloud.ServiceClient, id string, opts *UpdateOpts) UpdateResult {
	type update struct {
		Description *string `json:"display_description,omitempty"`
		Name        *string `json:"display_name,omitempty"`
	}

	type request struct {
		Volume update `json:"snapshot"`
	}

	reqBody := request{
		Volume: update{},
	}

	reqBody.Volume.Description = utils.MaybeString(opts.Description)
	reqBody.Volume.Name = utils.MaybeString(opts.Name)

	var res UpdateResult

	_, res.Err = perigee.Request("PUT", updateURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		ReqBody:     &reqBody,
		Results:     &res.Resp,
	})
	return res
}
