/*
Package quotas provides the ability to retrieve DNS quotas through the Designate API.

Example to Get a Detailed Quota Set

    projectID = "23d5d3f79dfa4f73b72b8b0b0063ec55"
    quotasInfo, err := quotas.GetDetail(dnsClient, projectID).Extract()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("quotas: %#v\n", quotasInfo)
*/
package quotas
