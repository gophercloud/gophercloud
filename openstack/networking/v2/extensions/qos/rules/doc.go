/*
Package rules provides the ability to retrieve and manage QoS policy rules through the Neutron API.

Example of Listing BandwidthLimitRules

    listOpts := rules.BandwidthLimitRuleListOpts{
        MaxKBps: 3000,
    }

    policyID := "501005fa-3b56-4061-aaca-3f24995112e1"

    allPages, err := rules.BandwidthLimitRuleList(networkClient, policyID, listOpts).AllPages()
    if err != nil {
        panic(err)
    }

    allBandwidthLimitRules, err := rules.ExtractBandwidthLimitRules(allPages)
    if err != nil {
        panic(err)
    }

    for _, bandwidthLimitRule := range allBandwidthLimitRules {
        fmt.Printf("%+v\n", bandwidthLimitRule)
    }
*/
package rules
