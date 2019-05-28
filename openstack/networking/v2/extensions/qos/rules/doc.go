/*
Package rules provides the ability to retrieve and manage QoS policy rules through the Neutron API.

Example of Listing BandwidthLimitRules

    listOpts := rules.BandwidthLimitRulesListOpts{
        MaxKBps: 3000,
    }

    policyID := "501005fa-3b56-4061-aaca-3f24995112e1"

    allPages, err := rules.BandwidthLimitRulesList(networkClient, policyID, listOpts).AllPages()
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

Example of Getting a single BandwidthLimitRule

    policyID := "501005fa-3b56-4061-aaca-3f24995112e1"
    ruleID   := "30a57f4a-336b-4382-8275-d708babd2241"

    rule, err := rules.GetBandwidthLimitRule(networkClient, policyID, ruleID).ExtractBandwidthLimitRule()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Rule: %+v\n", rule)
*/
package rules
