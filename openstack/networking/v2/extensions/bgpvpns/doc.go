package bgpvpns

/*
Package bgpvpns contains the functionality for working with Neutron BGP VPNs.

1. List BGP VPNs, a.k.a. GET /bgpvpn/bgpvpns

Example:

        pages, err := bgpvpns.List(client).AllPages(context.TODO())
        if err != nil {
                log.Panic(err)
        }
        allVPNs, err := bgpvpns.ExtractBGPVPNs(pages)
        if err != nil {
                log.Panic(err)
        }

        for _, bgpvpn := range allVPNs {
                log.Printf("%+v", bgpvpn)
        }

2. Get BGP VPN, a.k.a. GET /bgpvpn/bgpvpns/{id}

Example:
        p, err := bgpvpns.Get(context.TODO(), client, id).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *p)

3. Create BGP VPN, a.k.a. POST /bgpvpn/bgpvpns

Example:
        opts := bgpvpns.CreateOpts{
                name: "gophercloud-testing-bgpvpn".
        }
        r, err := bgpvpns.Create(context.TODO(), client, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", *r)

4. Delete BGP VPN, a.k.a. DELETE /bgpvpn/bgpvpns/{id}

Example:
        err := bgpvpns.Delete(context.TODO(), client, bgpVpnID).ExtractErr()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("BGP VPN deleted")


5. Update BGP VPN, a.k.a. PUT /bgpvpn/bgpvpns/{id}

Example:
        nameUpdated := "bgpvpn-name-updated"
        opts := bgpvpns.UpdateOpts{
                name: &nameUpdated,
        }
        p, err := bgpvpns.Update(context.TODO(), client, id, opts).Extract()
        if err != nil {
                log.Panic(err)
        }
        log.Printf("%+v", p)
*/
