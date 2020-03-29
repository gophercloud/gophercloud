/*
Package amphorae provides information and interaction with Amphorae
of OpenStack Load-balancing service.

Example to List Amphorae

	listOpts := amphorae.ListOpts{
		LoadbalancerID: "6bd55cd3-802e-447e-a518-1e74e23bb106",
	}

	allPages, err := amphorae.List(octaviaClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allAmphorae, err := amphorae.ExtractAmphorae(allPages)
	if err != nil {
		panic(err)
	}

	for _, amphora := range allAmphorae {
		fmt.Printf("%+v\n", amphora)
	}

Example to Failover an amphora

	ampID := "d67d56a6-4a86-4688-a282-f46444705c64"

	err := amphorae.Failover(octaviaClient, ampID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package amphorae
