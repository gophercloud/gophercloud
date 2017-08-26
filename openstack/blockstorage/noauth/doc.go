/*
Package noauth creates a "noauth" *gophercloud.ServiceClient.

Example of Creating a noauth Service Client

	ao, err := openstack.AuthOptionsFromEnv()
	provider, err := noauth.UnAuthenticatedClient(ao)
	client, err := noauth.NewBlockStorageV2(provider, noauth.EndpointOpts{
		CinderEndpoint: os.Getenv("CINDER_ENDPOINT"),
	})

e.g.
	CinderEndpoint: "http://cinder:8776/v2",
*/
package noauth
