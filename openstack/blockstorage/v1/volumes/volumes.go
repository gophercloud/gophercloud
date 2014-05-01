package volumes

type Volume map[string]interface{}
type CreateOpts map[string]interface{}

/*
type CreateOpts struct {
	Availability_zone   string            `json:"size"`
	Source_volid        string            `json:"source_volid"`
	Display_description string            `json:"display_description"`
	Snapshot_id         string            `json:"snapshot_id"`
	Size                int               `json:"size"`
	Display_name        string            `json:"display_name"`
	ImageRef            string            `json:"imageRef"`
	Volume_type         string            `json:"volume_type"`
	Bootable            bool              `json:"bootable"`
	Metadata            map[string]string `json:"metadata"`
}
*/
