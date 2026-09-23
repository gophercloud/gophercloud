/*
Package ec2tokens provides information and interaction with the EC2 token API
resource for the OpenStack Identity service.

For more information, see:
https://docs.openstack.org/api-ref/identity/v2-ext/

Example to Create a Token From an EC2 access and secret keys

	authOptions := auth.EC2TokenOpts{
		Access: "a7f1e798b7c2417cba4a02de97dc3cdc",
		Secret: "18f4f6761ada4e3795fa5273c30349b9",
	}

	token, err := ec2tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

Example auth client using EC2 access and secret keys

	ao := auth.AuthOptionsEC2{
		ServiceProvider: serviceProvider,
		AuthURL:         "http://localhost:5000/v3",
		Auth: auth.EC2TokenOpts{
			Access:      "a7f1e798b7c2417cba4a02de97dc3cdc",
			Secret:      "18f4f6761ada4e3795fa5273c30349b9",
			AllowReauth: true,
		},
	}

	client, err := openstack.AuthenticatedClient(context.TODO(), ao)
	if err != nil {
		panic(err)
	}
*/
package ec2tokens
