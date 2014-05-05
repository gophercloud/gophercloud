package snapshots

import (
	"github.com/racker/perigee"
	blockstorage "github.com/rackspace/gophercloud/openstack/blockstorage/v1"
)

func Create(c *blockstorage.Client, opts CreateOpts) (Snapshot, error) {
	var ss Snapshot
	h, err := c.GetHeaders()
	if err != nil {
		return ss, err
	}
	url := c.GetSnapshotsURL()
	_, err = perigee.Request("POST", url, perigee.Options{
		Results: &struct {
			Snapshot *Snapshot `json:"snapshot"`
		}{&ss},
		ReqBody: map[string]interface{}{
			"snapshot": opts,
		},
		MoreHeaders: h,
	})
	return ss, err
}
