/*
Package tsigkeys provides information and interaction with the TSIG key API
resource for the OpenStack DNS service.

TSIG (Transaction SIGnature) keys are used to authenticate DNS transactions
between servers, such as zone transfers and dynamic updates.

Example to List TSIG Keys

	listOpts := tsigkeys.ListOpts{
		Scope: "POOL",
	}

	allPages, err := tsigkeys.List(dnsClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allTSIGKeys, err := tsigkeys.ExtractTSIGKeys(allPages)
	if err != nil {
		panic(err)
	}

	for _, tsigkey := range allTSIGKeys {
		fmt.Printf("%+v\n", tsigkey)
	}

Example to Create a TSIG Key

	createOpts := tsigkeys.CreateOpts{
		Name:       "mytsigkey",
		Algorithm:  "hmac-sha256",
		Secret:     "example-secret-key-value==",
		Scope:      "POOL",
		ResourceID: "pool-id-here",
	}

	tsigkey, err := tsigkeys.Create(context.TODO(), dnsClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a TSIG Key

	tsigkeyID := "99d10f68-5623-4491-91a0-6daafa32b60e"
	tsigkey, err := tsigkeys.Get(context.TODO(), dnsClient, tsigkeyID).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a TSIG Key

	tsigkeyID := "99d10f68-5623-4491-91a0-6daafa32b60e"
	updateOpts := tsigkeys.UpdateOpts{
		Name:   "updatedname",
		Secret: "updated-secret-key-value==",
	}

	tsigkey, err := tsigkeys.Update(context.TODO(), dnsClient, tsigkeyID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a TSIG Key

	tsigkeyID := "99d10f68-5623-4491-91a0-6daafa32b60e"
	err := tsigkeys.Delete(context.TODO(), dnsClient, tsigkeyID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package tsigkeys
