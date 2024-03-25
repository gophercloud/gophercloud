/*
Package volumehost provides the ability to extend a volume result with
information about the Openstack host holding the volume. Example:

	type VolumeWithHost struct {
		volumes.Volume
		volumehost.VolumeHostExt
	}

	var allVolumes []VolumeWithHost

	allPages, err := volumes.List(client, nil).AllPages(context.TODO())
	if err != nil {
		panic("Unable to retrieve volumes: %s", err)
	}

	err = volumes.ExtractVolumesInto(allPages, &allVolumes)
	if err != nil {
		panic("Unable to extract volumes: %s", err)
	}

	for _, volume := range allVolumes {
		fmt.Println(volume.Host)
	}
*/
package volumehost
