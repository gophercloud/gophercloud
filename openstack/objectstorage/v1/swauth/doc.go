/*
Package swauth implements Swift's built-in authentication.

Example to Authenticate with swauth

	authOpts := swauth.AuthOpts{
		User: "project:user",
		Key:  "password",
	}

	swiftClient, err := swauth.NewObjectStorageV1(context.TODO(), providerClient, authOpts)
	if err != nil {
		panic(err)
	}
*/
package swauth
