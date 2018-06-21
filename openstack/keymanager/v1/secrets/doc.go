/*
Package secrets manages and retrieves secrets in the OpenStack Key Manager
Service.

Example to List Secrets

	createdQuery := &secrets.DateQuery{
		Date:   time.Date(2049, 6, 7, 1, 2, 3, 0, time.UTC),
		Filter: secrets.DateFilterLT,
	}

	listOpts := secrets.ListOpts{
		CreatedQuery: createdQuery,
	}

	allPages, err := secrets.List(client, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allSecrets, err := secrets.ExtractSecrets(allPages)
	if err != nil {
		panic(err)
	}

	for _, v := range allSecrets {
		fmt.Printf("%v\n", v)
	}

Example to Get a Secret

	secret, err := secrets.Get(client, secretID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", secret)

Example to Get a Payload

	payload, err := secrets.GetPayload(client, secretID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(payload))

Example to Create a Secrets

	createOpts := secrets.CreateOpts{
		Algorithm:         "aes",
		BitLength:          256,
		Mode:               "cbc",
		Name:               "mysecret",
		Payload:            "super-secret",
		PayloadContentType: "text/plain",
		SecretType:         secrets.OpaqueSecret,
	}

	secret, err := secrets.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Println(secret.SecretRef)

Example to Add a Payload

	updateOpts := secrets.UpdateOpts{
		ContentType: "text/plain",
		Payload:     "super-secret",
	}

	err := secrets.Update(client, secretID, updateOpts).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Delete a Secrets

	err := secrets.Delete(client, secretID).ExtractErr()
	if err != nil {
		panic(err)
	}

*/
package secrets
