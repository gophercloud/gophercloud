/*
Package shareattach provides the ability to attach and detach Manila shares
from servers. This requires the client to be set to microversion 2.97 or later.
The instance must be in SHUTOFF status to attach or detach a share.

Example to Attach a Share

	serverID := "7ac8686c-de71-4acb-9600-ec18b1a1ed6d"
	shareID := "87463836-f0e2-4029-abf6-20c8892a3103"

	computeClient.Microversion = "2.97"

	createOpts := shareattach.CreateOpts{
		ShareID: shareID,
		Tag:     "my-share",
	}

	result, err := shareattach.Create(context.TODO(), computeClient, serverID, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Detach a Share

	serverID := "7ac8686c-de71-4acb-9600-ec18b1a1ed6d"
	shareID := "ed081613-1c9b-4231-aa5e-ebfd4d87f983"

	computeClient.Microversion = "2.97"

	err := shareattach.Delete(context.TODO(), computeClient, serverID, shareID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package shareattach
