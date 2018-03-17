/*
Package profiles provides information and interaction with the profiles through
the OpenStack Compute service.

Lists all profiles and creates, shows information for, updates, and deletes a profile.

Example to List Profiles

	listOpts := profiles.ListOpts{
		Limit: 2,
	}

	allPages, err := profiles.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allProfiles, err := profiles.ExtractProfiles(allPages)
	if err != nil {
		panic(err)
	}

	for _, profile := range allProfiles {
		fmt.Printf("%+v\n", profile)
	}

Example to Create a profile

	createOpts := profiles.CreateOpts{
    Name: "testprofile",
    Spec: map[string]interface{}{
      "type":       "os.nova.server",
      "version":    "1.0",
      "properties": props,
    },
	}
  networks := &[...]map[string]interface{}{
		{"network": "sandbox-internal-net"},
		//{"network": "sandbox-internal-net2"},
	}
	props := &map[string]interface{}{
		"name":            "test_gopher_cloud_profile",
		"flavor":          "t2.micro",
		"image":           "centos7.3-latest",
		"networks":        networks,
		"security_groups": "",
	}

	profile, err := profiles.Create(computeClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package profiles
