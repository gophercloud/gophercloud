package volumes

type Volume struct {
	Status              string
	Display_name        string
	Attachments         []string
	Availability_zone   string
	Bootable            bool
	Created_at          string
	Display_description string
	Volume_type         string
	Snapshot_id         string
	Source_volid        string
	Metadata            map[string]string
	Id                  string
	Size                int
}
type CreateOpts map[string]interface{}
type GetOpts map[string]string
type DeleteOpts map[string]string
