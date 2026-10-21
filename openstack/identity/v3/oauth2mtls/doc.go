/*
Package oauth2mtls authenticates to Keystone using OAuth 2.0 mutual TLS client
authentication (RFC 8705).

The client certificate is mapped to a Keystone user, and the returned access
token is bound to that certificate. The same certificate must be presented to
protected service endpoints. The user must have a default project.

Example to authenticate a client using OAuth 2.0 mutual TLS

	client, err := openstack.NewClient("https://keystone.example.com:5000/v3")
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = tlsConfig
	client.HTTPClient.Transport = transport

	authOptions := &oauth2mtls.AuthOptions{
		OAuth2Endpoint: "https://keystone.example.com:5000/v3/OS-OAUTH2/token",
		ClientID:       "6c3145f4-313d-4910-b3a8-9dfc72da9e75",
		AllowReauth:    true,
	}

	err = openstack.AuthenticateV3(context.Background(), client, authOptions, gophercloud.EndpointOpts{})
	if err != nil {
		panic(err)
	}

See https://docs.openstack.org/keystone/latest/admin/oauth2-mtls-usage-guide.html.
*/
package oauth2mtls
