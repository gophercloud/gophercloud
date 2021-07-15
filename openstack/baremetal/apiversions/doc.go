/*
Package apiversions provides information about the versions supported by a specific Ironic API.

	Example to list versions

		allVersions, err := apiversions.List(baremetalClient).Extract()
		if err != nil {
			panic("unable to get API versions: " + err.Error())
		}

		for _, version := range allVersions.Versions {
			fmt.Printf("%+v\n", version)
		}

	Example to get a specific version

		actual, err := apiversions.Get(baremetalClient).Extract()
		if err != nil {
			panic("unable to get API version: " + err.Error())
		}

*/
package apiversions
