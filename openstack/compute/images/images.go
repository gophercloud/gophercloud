package images

import "github.com/mitchellh/mapstructure"

// Image is used for JSON (un)marshalling.
// It provides a description of an OS image.
//
// The Id field contains the image's unique identifier.
// For example, this identifier will be useful for specifying which operating system to install on a new server instance.
//
// The MinDisk and MinRam fields specify the minimum resources a server must provide to be able to install the image.
//
// The Name field provides a human-readable moniker for the OS image.
//
// The Progress and Status fields indicate image-creation status.
// Any usable image will have 100% progress.
//
// The Updated field indicates the last time this image was changed.
type Image struct {
	Created         string
	Id              string
	MinDisk         int
	MinRam          int
	Name            string
	Progress        int
	Status          string
	Updated         string
}

func GetImages(lr ListResults) ([]Image, error) {
	ia, ok := lr["images"]
	if !ok {
		return nil, ErrNotImplemented
	}
	ims := ia.([]interface{})

	images := make([]Image, len(ims))
	for i, im := range ims {
		imageObj := im.(map[string]interface{})
		err := mapstructure.Decode(imageObj, &images[i])
		if err != nil {
			return images, err
		}
	}
	return images, nil
}

