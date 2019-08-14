/*
Package tags manages Tags on Compute V2 servers.

This extension is available since 2.26 Compute V2 API microversion.

Example to List all server Tags

	client.Microversion = "2.62"

    serverTags, err = tags.List(client, server.ID).Extract()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Tags: %v\n", serverTags)
*/
package tags
