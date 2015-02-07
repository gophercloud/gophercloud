package volumeattach

import (
	"errors"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

// CreateOptsExt adds a VolumeAttachment option to the base CreateOpts.
type CreateOptsExt struct {
	servers.CreateOptsBuilder
	Device   string `json:"device"`
	VolumeID string `json:"volumeId"`
	ServerID string `json:"serverId"`
}

// ToServerCreateMap adds the volume_id, device, and optionally server_id to
// the base server creation options.
func (opts CreateOptsExt) ToServerCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToServerCreateMap()
	if err != nil {
		return nil, err
	}

	serverMap := base["server"].(map[string]interface{})
	serverMap["device"] = opts.Device
	serverMap["volume_id"] = opts.VolumeID
	serverMap["server_id"] = opts.ServerID

	return base, nil
}

// List returns a Pager that allows you to iterate over a collection of VolumeAttachments.
func List(client *gophercloud.ServiceClient, serverId string) pagination.Pager {
	return pagination.NewPager(client, listURL(client, serverId), func(r pagination.PageResult) pagination.Page {
		return VolumeAttachmentsPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder describes struct types that can be accepted by the Create call. Notable, the
// CreateOpts struct in this package does.
type CreateOptsBuilder interface {
	ToVolumeAttachmentCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies volume attachment creation or import parameters.
type CreateOpts struct {
	// Device is the device that the volume will attach to the instance as. Omit for "auto"
	Device string

	// VolumeID is the ID of the volume to attach to the instance
	VolumeID string

	// ServerID is the ID of the server that the volume will be attached to
	ServerID string
}

// ToVolumeAttachmentCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToVolumeAttachmentCreateMap() (map[string]interface{}, error) {
	if opts.VolumeID == "" {
		return nil, errors.New("Missing field required for volume attachment creation: VolumeID")
	}

	if opts.ServerID == "" {
		return nil, errors.New("Missing field required for volume attachment creation: ServerID")
	}

	volumeAttachment := make(map[string]interface{})
	volumeAttachment["volumeId"] = opts.VolumeID
	volumeAttachment["serverId"] = opts.ServerID
	if opts.Device != "" {
		volumeAttachment["device"] = opts.Device
	}

	return map[string]interface{}{"volumeAttachment": volumeAttachment}, nil
}

// Create requests the creation of a new volume attachment on the server
func Create(client *gophercloud.ServiceClient, serverId string, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToVolumeAttachmentCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = perigee.Request("POST", createURL(client, serverId), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		ReqBody:     reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// Get returns public data about a previously created VolumeAttachment.
func Get(client *gophercloud.ServiceClient, serverId, aId string) GetResult {
	var res GetResult
	_, res.Err = perigee.Request("GET", getURL(client, serverId, aId), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// Delete requests the deletion of a previous stored VolumeAttachment from the server.
func Delete(client *gophercloud.ServiceClient, serverId, aId string) DeleteResult {
	var res DeleteResult
	_, res.Err = perigee.Request("DELETE", deleteURL(client, serverId, aId), perigee.Options{
		MoreHeaders: client.AuthenticatedHeaders(),
		OkCodes:     []int{202},
	})
	return res
}
