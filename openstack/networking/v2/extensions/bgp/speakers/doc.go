package speakers

/*
Package speakers contains the functionality for working with Neutron bgp speakers.


1. List BGP Speakers, e.g. GET /bgp-speakers

Example:

        pages, err := speakers.List(client).AllPages(context.TODO())
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

        speaker, err := speakers.Get(context.TODO(), client, id).Extract()
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
        r, err := speakers.Create(context.TODO(), client, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *r)


5. Delete BGP Speaker, a.k.a. DELETE /bgp-speakers/{id}

Example:

        err := speakers.Delete(auth, speakerID).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("Speaker Deleted")


6. Update BGP Speaker

Example:

        opts := speakers.UpdateOpts{
                Name:                          "testing-bgp-speaker",
                AdvertiseTenantNetworks:       false,
                AdvertiseFloatingIPHostRoutes: true,
        }
        spk, err := speakers.Update(context.TODO(), client, bgpSpeakerID, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", spk)


7. Add BGP Peer, a.k.a. PUT /bgp-speakers/{id}/add_bgp_peer

Example:

        opts := speakers.AddBGPPeerOpts{BGPPeerID: bgpPeerID}
        r, err := speakers.AddBGPPeer(context.TODO(), client, bgpSpeakerID, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", r)


8. Remove BGP Peer, a.k.a. PUT /bgp-speakers/{id}/remove_bgp_peer

Example:

        opts := speakers.RemoveBGPPeerOpts{BGPPeerID: bgpPeerID}
        err := speakers.RemoveBGPPeer(context.TODO(), client, bgpSpeakerID, opts).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("Successfully removed BGP Peer")


9. Get advertised routes, a.k.a. GET /bgp-speakers/{id}/get_advertised_routes

Example:

        pages, err := speakers.GetAdvertisedRoutes(client, speakerID).AllPages(context.TODO())
        if err != nil {
                log.Panic(err)
        }
        routes, err := speakers.ExtractAdvertisedRoutes(pages)
        if err != nil {
                log.Panic(err)
        }
        for _, r := range routes {
                log.Printf("%+v", r)
        }


10. Add geteway network to BGP Speaker, a.k.a. PUT /bgp-speakers/{id}/add_gateway_network

Example:


        opts := speakers.AddGatewayNetworkOpts{NetworkID: networkID}
        r, err := speakers.AddGatewayNetwork(context.TODO(), client, speakerID, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", r)


11. Remove gateway network to BGP Speaker, a.k.a. PUT /bgp-speakers/{id}/remove_gateway_network

Example:

        opts := speakers.RemoveGatewayNetworkOpts{NetworkID: networkID}
        err := speakers.RemoveGatewayNetwork(context.TODO(), client, speakerID, opts).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("Successfully removed gateway network")
*/
