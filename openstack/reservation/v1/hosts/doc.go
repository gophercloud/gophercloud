/*
Package hosts manages compute hosts in the OpenStack Reservation service.

A host must be enrolled into Blazar's freepool before any reservation can draw
on it. Listing the freepool is administrative.

Example to list hosts

	allPages, err := hosts.List(reservationClient).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allHosts, err := hosts.ExtractHosts(allPages)
	if err != nil {
		panic(err)
	}

	for _, h := range allHosts {
		fmt.Printf("%+v\n", h)
	}
*/
package hosts
