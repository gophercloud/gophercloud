package peers

/*
Package peers contains the functionality for working with Neutron bgp peers.

1. List BGP Peers, a.k.a. GET /bgp-peers

Example:

        pages, err := peers.List(c).AllPages()
        if err != nil {
                log.Panic(err)
        }
        allPeers, err := peers.ExtractBGPPeers(pages)
        if err != nil {
                log.Panic(err)
        }

        for _, peer := range allPeers {
                log.Printf("%+v", peer)
        }

2. Get BGP Peer, a.k.a. GET /bgp-peers/{id}

Example:
        p, err := peers.Get(c, id).Extract()

        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *p)
*/
