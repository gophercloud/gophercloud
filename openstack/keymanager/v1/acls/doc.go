/*
Package acls manages acls in the OpenStack Key Manager Service.

All functions have a Secret and Container equivalent.

Example to Get a Secret's ACL

	acl, err := acls.GetSecretACL(context.TODO(), client, secretID).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", acl)

Example to Set a Secret's ACL

	users := []string{"uuid", "uuid"}
	iFalse := false
	setOpts := acls.SetOpts{
		Type:          "read",
		users:         &users,
		ProjectAccess: &iFalse,
	}

	aclRef, err := acls.SetSecretACL(context.TODO(), client, secretID, setOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", aclRef)

Example to Update a Secret's ACL

	users := []string{}
	setOpts := acls.SetOpts{
		Type:  "read",
		users: &users,
	}

	aclRef, err := acls.UpdateSecretACL(context.TODO(), client, secretID, setOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", aclRef)

Example to Delete a Secret's ACL

	err := acls.DeleteSecretACL(context.TODO(), client, secretID).ExtractErr()
	if err != nil {
		panci(err)
	}
*/
package acls
