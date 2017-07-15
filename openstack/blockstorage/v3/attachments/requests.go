package attachments

import (
	"github.com/gophercloud/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// CreateOpts are the valid API json options for creating a Cinder attachment
// I don't know that these need to be public?  Just use args to the method maybe?
type CreateOpts struct {
	ProjectID string            `json:project_id`
	ServerID  string            `json:instance_uuid`
	VolumeID  string            `json:"volume_id"`
	Connector map[string]string `json:"connector,omitempty"`
}

// UpdateOpts are the valid API json options for updating a Cinder attachment
type UpdateOpts struct {
	ProjectID    string            `json:project_id`
	AttachmentID string            `json:attachment_id`
	Connector    map[string]string `json:"connector"`
}

// there's only two options to list, we probably don't even need a struct for them, but it's convenient
// for the json part, we'll just use vars in the function for now and use the opts struct internally
type listOpts struct {
	serverID   string `json:instance_uuid`
	allTenants bool   `json:all_tenants`
}

/* XXXDelete(project_id, attachment_id)
*  XXXGet(project_id, attachment_id)
*  XXXList(project_id, sort_key, sort_dir, limit, offset, marker) //Add a withDetail flag (only changes response)
*  Create
*  Update
 */

func (opts CreateOpts) toAttachmentCreateMap() (map[string]interface{}, error) {
	return gophercloud.RequestBody(opts, "attachment")
}

// Create attachment creates an attachment record.  This can be a simple reserve (no connector info) or a full
// attach process if the optional connector info is included.  NOTE if you don't provide the connector and the
// volume will be reserved and you can't connect it until you finish the process with an `Update` call that
// provides the required connector info
func Create(client *gophercloud.ServiceClient, opts CreateOpts) (result UpdateResult) {
	client.Microversion = "3.27"
	b, err := opts.toAttachmentCreateMap()
	if err != nil {
		result.Err = err
		return result
	}
	_, result.Err = client.Post(createURL(client), b, &result.Body, &gophercloud.RequestOpts{
		OKCodes: []int{202},
	})
	return result
}

func Update(client *gohpercoud.ServiceClient, attachmentID string, connector map[string]string) (result CreateResult) {
	client.Microversion = "3.27"
	//TODO(j-griffith) need the whole body construction etc etc
	if err != nil {
		result.Err = err
		return result
	}
	_, result.Err = client.Post(createURL(client), b, &result.Body, &gophercloud.RequestOpts{
		OKCodes: []int{202},
	})
	return result
}

// Delete destroys the attachment record specified by ID, not that this does not deal with cleaning up the actual
// connection on the Host side, but only the target and database.  For things like iSCSI you should make sure you
// have performed the iscsiadm logout elsewhere NOTE: Cinder attachment API's require micro-version >= 3.27
func Delete(client *gophercloud.ServiceClient, id string) (result DeleteResult) {
	client.Microversion = "3.27"
	_, result.Err = client.Delete(deleteURL(client, id), nil)
	return result
}

// Get fetches the attachment record specified by ID
func Get(client *gophercloud.ServiceClient, id string) (result GetResult) {
	client.Microversion = "3.27"
	_, result.Err = client.Get(getURL(client, id), &result.Body, nil)
	return result
}

// List returns a list of attachments known by Cinder
func List(client *gophercloud.ServiceClient, serverID string, allTenants bool) pagination.Pager {
	client.Microversion = "3.27"
	opts := new(listOpts)
	if serverID != "" {
		opts.serverID = serverID
	}
	if allTenants == true {
		opts.allTenants = true
	}

	url := listURL(client)
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return pagination.Pager{Err: err}
	}
	url += q.String()
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AttachmentPage{pagination.SinglePageBase(r)}
	})
}
