/*
Package tags manages Tags on Compute V2 servers.

This extension is available since 2.26 Compute V2 API microversion.

Example to List all server Tags

	client.Microversion = "2.62"

    serverTags, err := tags.List(client, serverID).Extract()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Tags: %v\n", serverTags)

Example to Check if the specific Tag exists on a server

    client.Microversion = "2.62"

    exists, err := tags.Check(client, serverID, tag).Extract()
    if err != nil {
        log.Fatal(err)
    }

    if exists {
        log.Printf("Tag %s is set\n", tag)
    } else {
        log.Printf("Tag %s is not set\n", tag)
    }

Example to Replace all Tags on a server

    client.Microversion = "2.62"

    newTags, err := tags.Replace(client, serverID, tags. tags.ReplaceOpts{Tags: []string{"foo", "bar"}}).Extract()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("New tags: %v\n", newTags)
*/
package tags
