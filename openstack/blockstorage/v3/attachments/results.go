package attachments

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Attachment representation of a Cinder Attachment record
type Attachment struct {
	// ID is a UUID for the individual attachment record
	ID string `json:"id"`
	// VolumeID is a UUID for the volume associated with this attachment record
	VolumeID string `json:"volume_id"`
	// ServerID is an identifier specifying the server (Instance) this attachment is allocated for
	// In a traditional OpenStack/Nova deployment this will be the Nova Instance/Server UUID, in
	// a Standalone environment we may or may not have that additional nested layer so we `make something up`
	ServerID string `json:"instance_uuid"`
	// AttachedHost identifier of the actual host machine we're making the attachment to.  This could be IP or network/hostname.
	AttachedHost string `json:"attached_host"`
	// Mountpoint is the mountpoint inside of the Nova Server/Instance (ie /dev/vdb).
	Mountpoint string `json:"mountpoint"`
	// AttachedAt time-stamp when the volume was moved into a reserve/attached state
	AttachedAt time.Time `json:"attach_time"`
	// DetachedAt time-stamp
	DetachedAt time.Time `json:"detach_time"`
	// Status of the attachment, ie `reserve`, `in-use`, `attaching`, `available`
	Status string `json:"attach_status"`
	// Mode of the attachment (r/w, r/o)
	Mode string `json:"attach_mode"`
	// ConnectionInfo hall al the gory details on HOW to actuallly make the connection
	ConnectionInfo map[string]string `json:"connection_info"`
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type AttachmentPage struct {
	pagination.SinglePage
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
}

type UpdateResult struct {
}

type commonResult struct {
	gophercloud.Result
}

func (a *Attachment) UnmarshalJSON(b []byte) error {
	type tmp Attachment
	var s struct {
		tmp
		AttachedAt gophercloud.JSONRFC3339MilliNoZ `json:"attached_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*a = Attachment(s.tmp)

	a.AttachedAt = time.Time(s.AttachedAt)

	return err
}
