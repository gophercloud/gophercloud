/*
Package noauth creates a "noauth" *gophercloud.ServiceClient for use in Cinder
environments configured with the noauth authentication middleware.

Example of Creating a noauth Service Client

	ao, err := openstack.AuthOptionsFromEnv()
	provider, err := noauth.NewClient(ao)
	client, err := noauth.NewBlockStorageV2(provider, noauth.EndpointOpts{
		CinderEndpoint: os.Getenv("CINDER_ENDPOINT"),
	})

	An example of a CinderEndpoint would be: http://example.com:8776/v2,
*/
package noauth
