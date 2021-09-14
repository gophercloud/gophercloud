package volumehost

// VolumeHostExt is an extension to the base Volume object
type VolumeHostExt struct {
	// Host is the identifier of the host holding the volume.
	Host string `json:"os-vol-host-attr:host"`
}
