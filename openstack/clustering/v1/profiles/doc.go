/*
Example to Create a profile

	networks := []map[string]interface{} {
		{"network": "test-network"},
	}

	props := map[string]interface{}{
		"name":            "test_gophercloud_profile",
		"flavor":          "t2.micro",
		"image":           "centos7.3-latest",
		"networks":        networks,
		"security_groups": "",
	}

	createOpts := profiles.CreateOpts {
		Name: "test_profile",
		Spec: profiles.Spec{
			Type:       "os.nova.server",
			Version:    "1.0",
			Properties: props,
		},
	}

	profile, err := profiles.Create(serviceClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println("Profile", profile)

Example to Get profile

	profile, err := profiles.Get(serviceClient, "profile-name").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Print("profile", profile)


Example to List profiles

	profiles.List(serviceClient, profiles.ListOpts{Limit: 2}).EachPage(func(page pagination.Page) (bool, error) {
		allProfiles, err := profiles.ExtractProfiles(page)
		if err != nil {
			panic(err)
		}

		for _, profile := range allProfiles {
			fmt.Printf("%+v\n", profile)
		}
		return true, nil
	})

Example to Update profile

	profile, err := profiles.Update(serviceClient, profileName, profiles.UpdateOpts{Name: newProfileName}).Extract()
    if err != nil {
		panic(err)
    }

    fmt.Print("profile", profile)

Example to Delete profile

	profileID := "6dc6d336e3fc4c0a951b5698cd1236ee"
	err := profiles.Delete(serviceClient, profileID).ExtractErr()
	if err != nil {
		panic(err)
	}

*/
package profiles
