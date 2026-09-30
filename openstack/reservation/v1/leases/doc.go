/*
Package leases manages leases in the OpenStack Reservation service.

A lease reserves a set of resources over a period of time. Each lease holds one
or more reservations, and Blazar schedules events to start and end it. Dates are
truncated to the minute, and a lease cannot start in the past.

Example to list leases

	listOpts := leases.ListOpts{}

	allPages, err := leases.List(reservationClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allLeases, err := leases.ExtractLeases(allPages)
	if err != nil {
		panic(err)
	}

	for _, l := range allLeases {
		fmt.Printf("%+v\n", l)
	}

Example to create a lease reserving two hosts

	createOpts := leases.CreateOpts{
		Name:      "my_lease",
		StartDate: time.Now().UTC().Add(time.Hour),
		EndDate:   time.Now().UTC().Add(24 * time.Hour),
		Reservations: []leases.ReservationOptsBuilder{
			leases.HostReservationOpts{
				Min:                  1,
				Max:                  2,
				HypervisorProperties: `[">=", "$vcpus", "4"]`,
			},
		},
	}

	lease, err := leases.Create(context.TODO(), reservationClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to create a lease reserving instance capacity

	affinity := false

	createOpts := leases.CreateOpts{
		Name:      "lease_bar",
		StartDate: time.Now().UTC().Add(time.Hour),
		EndDate:   time.Now().UTC().Add(24 * time.Hour),
		Reservations: []leases.ReservationOptsBuilder{
			leases.InstanceReservationOpts{
				Amount:   2,
				VCPUs:    1,
				MemoryMB: 2048,
				DiskGB:   20,
				Affinity: &affinity,
			},
		},
	}

	lease, err := leases.Create(context.TODO(), reservationClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to create a lease reserving instances of an existing flavor

	createOpts := leases.CreateOpts{
		Name:      "lease_baz",
		StartDate: time.Now().UTC().Add(time.Hour),
		EndDate:   time.Now().UTC().Add(24 * time.Hour),
		Reservations: []leases.ReservationOptsBuilder{
			leases.FlavorInstanceReservationOpts{
				Amount:   2,
				FlavorID: "1e1a9b1e-1f0a-4d1e-9f1a-0b1c2d3e4f5a",
			},
		},
	}

	lease, err := leases.Create(context.TODO(), reservationClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to get a lease

	lease, err := leases.Get(context.TODO(), reservationClient, "b179d3b5-6014-44a7-969b-0d333a969631").Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", lease)

Example to extend a lease and grow its host reservation

	endDate := time.Now().UTC().Add(48 * time.Hour)

	updateOpts := leases.UpdateOpts{
		EndDate: &endDate,
		Reservations: []leases.UpdateReservationOpts{
			{
				ID:  "b4675fff-ec59-480b-9399-9a74dffbc5c1",
				Min: 2,
				Max: 3,
			},
		},
	}

	lease, err := leases.Update(context.TODO(), reservationClient, "b179d3b5-6014-44a7-969b-0d333a969631", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to delete a lease

	err := leases.Delete(context.TODO(), reservationClient, "b179d3b5-6014-44a7-969b-0d333a969631").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package leases
