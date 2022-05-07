package speakers

/*
Package speakers contains the functionality for working with Neutron bgp speakers.


1. List BGP Speakers, e.g. GET /bgp-speakers

Example:

        pages, err := speakers.List(c).AllPages()
        if err != nil {
                log.Panic(err)
        }
        allSpeakers, err := speakers.ExtractBGPSpeakers(pages)
        if err != nil {
                log.Panic(err)
        }

        for _, speaker := range allSpeakers {
                log.Printf("%+v", speaker)
        }


2. Get BGP speakers, e.g. GET /bgp-speakers/{id}

Example:

        speaker, err := speakers.Get(c, id).Extract()
        if err != nil {
                log.Panic(nil)
        }
        log.Printf("%+v", *speaker)


3. Create BGP Speaker, a.k.a. POST /bgp-speakers

Example:

	opts := speakers.CreateOpts{
		IPVersion:                     6,
		AdvertiseFloatingIPHostRoutes: false,
		AdvertiseTenantNetworks:       true,
		Name:                          "gophercloud-testing-bgp-speaker",
		LocalAS:                       "2000",
		Networks:                      []string{},
	}
        r, err := speaker.Create(c, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *r)


5. Delete BGP Speaker, a.k.a. DELETE /bgp-speakers/{id}

Example:

        err := speaker.Delete(auth, speakerID).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("Speaker Deleted")
*/
