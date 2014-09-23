package snapshots

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/utils"
)

type CreateOpts struct {
	Description string
	Force       bool
	Metadata    map[string]interface{}
	Name        string
	VolumeID    string
}

func Create(client *gophercloud.ServiceClient, opts CreateOpts) (*Snapshot, error) {
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

	type response struct {
		Snapshot Snapshot `json:"snapshot"`
	}

	var respBody response

	_, err := perigee.Request("POST", snapshotsURL(client), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{201},
		ReqBody:     &reqBody,
		Results:     &respBody,
	})
	if err != nil {
		return nil, err
	}

	return &respBody.Snapshot, nil
}

func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var gr GetResult
	_, err := perigee.Request("GET", snapshotURL(client, id), perigee.Options{
		Results:     &gr.r,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	gr.err = err
	return gr
}
