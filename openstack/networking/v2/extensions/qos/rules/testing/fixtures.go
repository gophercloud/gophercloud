package testing

// BandwidthLimitRulesListResult represents a raw result of a List call to BandwidthLimitRules.
const BandwidthLimitRulesListResult = `
{
    "bandwidth_limit_rules": [
        {
            "max_kbps": 3000,
            "direction": "egress",
            "id": "30a57f4a-336b-4382-8275-d708babd2241",
            "max_burst_kbps": 300
        }
    ]
}
`

// BandwidthLimitRulesGetResult represents a raw result of a Get call to a specific BandwidthLimitRule.
const BandwidthLimitRulesGetResult = `
{
    "bandwidth_limit_rule": {
        "max_kbps": 3000,
        "direction": "egress",
        "id": "30a57f4a-336b-4382-8275-d708babd2241",
        "max_burst_kbps": 300
    }
}
`

// BandwidthLimitRulesCreateRequest represents a raw body of a Create BandwidthLimitRule call.
const BandwidthLimitRulesCreateRequest = `
{
    "bandwidth_limit_rule": {
        "max_kbps": 2000,
        "max_burst_kbps": 200
    }
}
`

// BandwidthLimitRulesCreateResult represents a raw result of a Create BandwidthLimitRule call.
const BandwidthLimitRulesCreateResult = `
{
    "bandwidth_limit_rule": {
        "max_kbps": 2000,
        "id": "30a57f4a-336b-4382-8275-d708babd2241",
        "max_burst_kbps": 200
    }
}
`

// BandwidthLimitRulesUpdateRequest represents a raw body of a Update BandwidthLimitRule call.
const BandwidthLimitRulesUpdateRequest = `
{
    "bandwidth_limit_rule": {
        "max_kbps": 500,
        "max_burst_kbps": 0
    }
}
`

// BandwidthLimitRulesUpdateResult represents a raw result of a Update BandwidthLimitRule call.
const BandwidthLimitRulesUpdateResult = `
{
    "bandwidth_limit_rule": {
        "max_kbps": 500,
        "id": "30a57f4a-336b-4382-8275-d708babd2241",
        "max_burst_kbps": 0
    }
}
`
