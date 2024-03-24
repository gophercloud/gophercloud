package peers

/*
Package peers contains the functionality for working with Neutron bgp peers.

1. List BGP Peers, a.k.a. GET /bgp-peers

Example:

        pages, err := peers.List(c).AllPages(context.TODO())
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

3. Create BGP Peer, a.k.a. POST /bgp-peers

Example:
        var opts peers.CreateOpts
        opts.AuthType = "md5"
        opts.Password = "notSoStrong"
        opts.RemoteAS = 20000
        opts.Name = "gophercloud-testing-bgp-peer"
        opts.PeerIP = "192.168.0.1"
        r, err := peers.Create(c, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *r)

4. Delete BGP Peer, a.k.a. DELETE /bgp-peers/{id}

Example:

        err := peers.Delete(c, bgpPeerID).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("BGP Peer deleted")


5. Update BGP Peer, a.k.a. PUT /bgp-peers/{id}

Example:

        var opt peers.UpdateOpts
        opt.Name = "peer-name-updated"
        opt.Password = "superStrong"
        p, err := peers.Update(c, id, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", p)
*/
