/*
Package vlantransparent provides the ability to retrieve and manage networks
with the vlan-transparent extension through the Neutron API.

Example of Listing Networks with the vlan-transparent extension

    iTrue := true
    networkListOpts := networks.ListOpts{}
    listOpts := vlantransparent.ListOptsExt{
        ListOptsBuilder: networkListOpts,
        VLANTransparent: &iTrue,
    }

    type NetworkWithVLANTransparentExt struct {
        networks.Network
        vlantransparent.NetworkVLANTransparentExt
    }

    var allNetworks []NetworkWithVLANTransparentExt

    allPages, err := networks.List(networkClient, listOpts).AllPages()
    if err != nil {
        panic(err)
    }

    err = networks.ExtractNetworksInto(allPages, &allNetworks)
    if err != nil {
        panic(err)
    }

    for _, network := range allNetworks {
        fmt.Println("%+v\n", network)
    }
*/
package vlantransparent
