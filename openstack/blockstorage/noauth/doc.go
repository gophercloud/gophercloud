/*
Package noauth creates a "noauth" *gophercloud.ServiceClient for use in Cinder
environments configured with the noauth authentication middleware.

Example of Creating a noauth Service Client

	provider, err := noauth.NewClient(auth.NoAuthV2Opts{
		Username:   os.Getenv("OS_USERNAME"),
		TenantName: os.Getenv("OS_TENANT_NAME"),
	})
	client, err := noauth.NewBlockStorageNoAuthV2(provider, noauth.EndpointOpts{
		CinderEndpoint: os.Getenv("CINDER_ENDPOINT"),
	})

Cinder API v3 uses synthetic user and project IDs:

	provider, err := noauth.NewClient(auth.NoAuthV3Opts{
		UserID:    os.Getenv("OS_USER_ID"),
		ProjectID: os.Getenv("OS_PROJECT_ID"),
	})
	client, err := noauth.NewBlockStorageNoAuthV3(provider, noauth.EndpointOpts{
		CinderEndpoint: os.Getenv("CINDER_ENDPOINT"),
	})

An example CinderEndpoint is http://example.com:8776/v3.
*/
package noauth
