package peer

/*
Package peer contains the functionality for working with Neutron bgp peer.

1. List BGP Peers, a.k.a. GET /bgp-peers

Example:

        pages, err := peer.List(c).AllPages()
        if err != nil {
                log.Panic(err)
        }
        allPeers, err := peer.ExtractBGPPeers(pages)
        if err != nil {
                log.Panic(err)
        }

        for _, peer := range allPeers {
                log.Printf("%+v", peer)
        }

2. Get BGP Peer, a.k.a. GET /bgp-peers/{id}

Example:
        p, err := peer.Get(c, id).Extract()

        if err != nil {
                log.Panic(nil)
        }
        log.Printf("%+v", *p)


3. Create BGP Peer, a.k.a. POST /bgp-peers

Example:

        var opts peer.CreateOpts
        opts.AuthType = "md5"
        opts.Password = "notSoStrong"
        opts.RemoteAS = 20000
        opts.Name = "gophercloud-testing-bgp-peer"
        opts.PeerIP = "192.168.0.1"
        r, err := peer.Create(c, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *r)


4. Delete BGP Peer, a.k.a. DELETE /bgp-peers/{id}

Example:

        err := peer.Delete(c, bgpPeerID).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("BGP Peer deleted")

5. Update BGP Peer, a.k.a. PUT /bgp-peers/{id}

        var optPeer peer.UpdateOpts
        optPeer.Name = "peer-name-updated"
        optPeer.Password = "superStrong"
        p, err := peer.Update(c, id, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", p)
*/
