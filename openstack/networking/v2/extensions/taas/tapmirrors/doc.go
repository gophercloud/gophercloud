/*
Package tapmirrors manages and retrieves Tap Mirrors in the OpenStack Networking Service.

Example to Create a Tap Mirror

	createopts := tapmirrors.CreateOpts{
		Name:        "tapmirror1",
		Description: "Description of tapmirror1",
		PortID:      "a25290e9-1a54-4c26-a5b3-34458d122acc",
		MirrorType:  tapmirrors.MirrorTypeErspanv1,
		RemoteIP:    "192.168.54.217",
		Directions: tapmirrors.Directions{
			In:  "1",
			Out: "2",
		},
	}

	mirror, err := tapmirrors.Create(context.TODO(), networkClient, createopts).Extract()
	if err != nil {
		panic(err)
	}
*/
package tapmirrors
