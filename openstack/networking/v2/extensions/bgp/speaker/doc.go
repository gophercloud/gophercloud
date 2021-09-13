package speaker

/*
Package speaker contains the functionality for working with Neutron bgp speaker.


1. List BGP Speakers, a.k.a. GET /bgp-speakers

Example:

        pages, err := speaker.List(c).AllPages()
        if err != nil {
                log.Panic(err)
        }
        allSpeakers, err := speaker.ExtractBGPSpeakers(pages)
        if err != nil {
                log.Panic(err)
        }

        for _, spk := range allSpeakers {
                log.Printf("%+v", spk)
        }


2. Get BGP Speaker, a.k.a. 94cfb61e-8614-4aa1-8191-ade11c0fc2c3/{id}

Example:

        p, err := speaker.Get(c, id).Extract()
        if err != nil {
                log.Panic(nil)
        }
        log.Printf("%+v", *p)


3. Create BGP Speaker, a.k.a. POST /bgp-speakers

Example:

        name := "gophercloud-testing-bgp-speaker"
        localas := "2000"
        networks := []string{}

        m := make(map[string]string)
        // IPVersion should be either "4" or "6", by defualt it is "4"
        m["IPVersion"] = "6"
        // uncomment the following line to set "AdvertiseFloatingIPHostRoutes" to true"
        m["AdvertiseFloatingIPHostRoutes"] = "false"
        // uncomment the following line to set "AdvertiseTenantNetworks" to   true
        m["AdvertiseTenantNetworks"] = "false"

        opts := speaker.BuildCreateOpts(name, localas, networks, m)
        r, err := speaker.Create(c, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *r)


4. Update BGP Speaker, a.k.a. PUT /bgp-speakers/{id}

Example:

	m := make(map[string]string)
        m["Name"] = "testing-bgp-speaker"
	// The values for AdvertiseTenantNetworks and AdvertiseFloatingIPHostRoutes
	// are case insensitive
        m["AdvertiseFloatingIPHostRoutes"] = "TRUE"
        m["AdvertiseTenantNetworks"] = "False"
	opts := speaker.BuildUpdateOpts(c, speakerID, m)
	shhgs.UpdateBGPSpeaker(c, speakerID, opts)


5. Delete BGP Speaker, a.k.a. DELETE /bgp-speakers/{id}

Example:

        err := speaker.Delete(auth, speakerID).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("Speaker Deleted")

6. Add BGP Peer, a.k.a. PUT /bgp-speakers/{id}/add_bgp_peer

Example:

        r, err := speaker.AddBGPPeer(c, bgpSpeakerID, bgpPeerID).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", r)


7. Remove BGP Peer, a.k.a. PUT /bgp-speakers/{id}/remove_bgp_peer

Example:

        r, err := speaker.RemoveBGPPeer(c, bgpSpeakerID, bgpPeerID).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", r)

8. Add geteway network to BGP Speaker, a.k.a. PUT /bgp-speakers/{id}/add_gateway_network

Example:

        r, err := speaker.AddGatewayNetwork(c, speakerID, networkID).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", r)


9. Remove gateway network to BGP Speaker, a.k.a. PUT /bgp-speakers/{id}/remove_gateway_network

Example:

        r, err := speaker.RemoveGatewayNetwork(c, speakerID, networkID).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", r)


10. Get advertised routes, a.k.a. GET /bgp-speakers/{id}/get_advertised_routes

Example:

        pages, err := speaker.GetAdvertisedRoutes(c, speakerID).AllPages()
        if err != nil {
                log.Panic(err)
        }
        allPages, err := speaker.ExtractAdvertisedRoutes(pages)
        if err != nil {
                log.Panic(err)
        }
        for _, r := range allPages {
                log.Printf("%+v", r)
        }
*/
