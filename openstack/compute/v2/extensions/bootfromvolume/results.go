package bootfromvolume

type BlockDeviceMapping struct {
	BootIndex           int    `json:"boot_index"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
	DestinationType     string `json:"destination_type"`
	SourceType          string `json:"source_type"`
	UUID                string `json:"uuid"`
	VolumeSize          int    `json:"volume_size"`
}
