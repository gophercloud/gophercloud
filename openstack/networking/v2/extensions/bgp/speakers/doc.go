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
*/
