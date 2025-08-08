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

Example to Show the details of a specific Tap Mirror by ID

	tapMirror, err := tapmirrors.Get(context.TODO(), networkClient, "f2b08c1e-aa81-4668-8ae1-1401bcb0576c").Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Tap Mirror

	err = tapmirrors.Delete(context.TODO(), networkClient, "5291b189-fd84-46e5-84bd-78f40c05d69c").ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Update an Tap Mirror

	name := "updated name"
	description := "updated description"
	updateOps := tapmirrors.UpdateOpts{
		Description: &description,
		Name:        &name,
	}

	updatedTapMirror, err := tapmirrors.Update(context.TODO(), networkClient, "5c561d9d-eaea-45f6-ae3e-08d1a7080828", updateOps).Extract()
	if err != nil {
		panic(err)
	}

Example to List Tap Mirrors

	allPages, err := tapmirrors.List(networkClient, nil).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allTapMirrors, err := tapmirrors.ExtractTapMirrors(allPages)
	if err != nil {
		panic(err)
	}
*/
package tapmirrors
