/*
Package trusts enables management of OpenStack Identity Trusts.

Example to Create a Trust

	expiresAt := time.Date(2019, 12, 1, 14, 0, 0, 999999999, time.UTC)
	createOpts := trusts.CreateOpts{
	    ExpiresAt:         &expiresAt,
	    Impersonation:     true,
	    AllowRedelegation: true,
	    ProjectID:         "9b71012f5a4a4aef9193f1995fe159b2",
	    Roles: []trusts.Role{
	        {
	            Name: "member",
	        },
	    },
	    TrusteeUserID: "ecb37e88cc86431c99d0332208cb6fbf",
	    TrustorUserID: "959ed913a32c4ec88c041c98e61cbbc3",
	}

	trust, err := trusts.Create(context.TODO(), identityClient, createOpts).Extract()
	if err != nil {
	    panic(err)
	}

	fmt.Printf("Trust: %+v\n", trust)

Example to Delete a Trust

	trustID := "3422b7c113894f5d90665e1a79655e23"
	err := trusts.Delete(context.TODO(), identityClient, trustID).ExtractErr()
	if err != nil {
	    panic(err)
	}

Example to Get a Trust

	trustID := "3422b7c113894f5d90665e1a79655e23"
	err := trusts.Get(context.TODO(), identityClient, trustID).ExtractErr()
	if err != nil {
	    panic(err)
	}

Example to List a Trust

	listOpts := trusts.ListOpts{
		TrustorUserId: "3422b7c113894f5d90665e1a79655e23",
	}

	allPages, err := trusts.List(identityClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allTrusts, err := trusts.ExtractTrusts(allPages)
	if err != nil {
		panic(err)
	}

	for _, trust := range allTrusts {
		fmt.Printf("%+v\n", region)
	}
*/
package trusts
