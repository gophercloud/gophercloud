/*
Package auth authenticates to OpenStack Identity (Keystone), builds authentication
requests, and interprets token responses independently of the Identity resource
packages. Applications should use this package rather than manually creating
tokens through openstack/identity.

# Password authentication

Wrap credential options in AuthOptionsV2 or AuthOptionsV3 to select the Identity
API version. Authenticate returns an AuthResult containing the token, expiry,
user, scope, and service catalog. A nil provider uses a default provider client;
pass a configured provider to control HTTP transport and TLS settings.

	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3PasswordOpts{
			Username:       "username",
			Password:       "password",
			UserDomainName: "Default",
			Scope:          &auth.Scope{ProjectID: "project-id"},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.TokenID)

Use UserID instead of Username to identify a user by ID. When using Username,
identify the user's domain with UserDomainID or UserDomainName. The user's domain
is distinct from the requested token scope.

For Identity v2, use V2PasswordOpts with TenantID or TenantName for a scoped token;
leave both empty for an unscoped token:

	opts := auth.AuthOptionsV2{
		AuthURL: "https://identity.example.com/v2.0",
		Auth: auth.V2PasswordOpts{
			Username: "username",
			Password: "password",
			TenantID: "tenant-id",
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

# Scoping and existing tokens

For Identity v3, set Scope on password, token, or multifactor options. A nil scope
requests an unscoped token. Select one scope, for example:

	projectByID := auth.Scope{ProjectID: "project-id"}
	projectByName := auth.Scope{ProjectName: "project", ProjectDomainName: "Default"}
	domain := auth.Scope{DomainID: "domain-id"}
	system := auth.Scope{System: true}
	trust := auth.Scope{TrustID: "trust-id"}

ProjectDomainID or ProjectDomainName identifies a named project's domain;
DomainID or DomainName instead requests a domain-scoped token.

Use V3TokenOpts to exchange an existing token for a new token, optionally with
a different scope:

	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3TokenOpts{
			Token: "existing-token",
			Scope: &auth.Scope{ProjectID: "project-id"},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

V2TokenOpts provides the corresponding Identity v2 operation, with TenantID or
TenantName selecting its scope.

# Application credentials

Application credentials carry their own project and access restrictions. Supply
an application credential ID and secret, or a name and user identity:

	credentials := auth.V3ApplicationCredentialOpts{
		ApplicationCredentialID:     "credential-id",
		ApplicationCredentialSecret: "credential-secret",
	}
	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth:    credentials,
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

# TOTP and multifactor authentication

V3TOTPOpts supplies a time-based one-time passcode. Combine it with other methods
using V3MultifactorOpts, placing the scope on the combined options:

	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3MultifactorOpts{
			AuthMethods: []auth.AuthOptionsBuilderV3{
				auth.V3PasswordOpts{UserID: "user-id", Password: "password"},
				auth.V3TOTPOpts{UserID: "user-id", Passcode: "123456"},
			},
			Scope: &auth.Scope{ProjectID: "project-id"},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

# Loading options and configuring a provider

AuthOptionsFromEnv selects an authenticator from the OS_* environment variables.
Use the openstack convenience functions to authenticate and configure token
storage, endpoint lookup, and optional reauthentication together:

	opts, err := auth.AuthOptionsFromEnv()
	if err != nil {
		panic(err)
	}
	provider, err := openstack.AuthenticatedClient(context.Background(), opts)
	if err != nil {
		panic(err)
	}

For clouds.yaml, parse the cloud configuration, then pass it to
AuthOptionsFromCloud. Apply the parsed TLS configuration both to authentication
version discovery and to the provider's HTTP transport:

	cloud, endpointOpts, tlsConfig, err := clouds.Parse()
	if err != nil {
		panic(err)
	}
	opts, err := auth.AuthOptionsFromCloud(cloud, auth.WithTLSConfig(tlsConfig))
	if err != nil {
		panic(err)
	}
	provider, err := config.NewProviderClient(context.Background(), opts, config.WithTLSConfig(tlsConfig))
	if err != nil {
		panic(err)
	}
	network, err := openstack.NewNetworkV2(context.Background(), provider, endpointOpts)
	if err != nil {
		panic(err)
	}

If authentication has already been performed, openstack.NewClientFromAuthResult
configures a provider from the returned result. Pass
openstack.WithReauth(opts, *result) to enable reauthentication when permitted:

	provider, err := openstack.NewClientFromAuthResult(
		*result, openstack.WithReauth(opts, *result),
	)
	if err != nil {
		panic(err)
	}

# OAuth1 authentication

V3OAuth1Opts exchanges an existing OAuth1 consumer and access token for a
Keystone token. It signs the request against the authentication endpoint.
Manage consumers and obtain OAuth1 access tokens separately using the
openstack/identity/v3/oauth1 package.

	opts := auth.AuthOptionsV3{
		AuthURL: "https://identity.example.com/v3",
		Auth: auth.V3OAuth1Opts{
			ConsumerKey:     "consumer-key",
			ConsumerSecret:  "consumer-secret",
			Token:           "access-token",
			TokenSecret:     "access-secret",
			SignatureMethod: auth.OAuth1HMACSHA1,
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

# EC2 authentication

AuthOptionsEC2 uses EC2TokenOpts to sign an AWS v2 or v4 request and exchange it
for a Keystone token. ServiceProvider is a separately authenticated provider
with permission to call Keystone's EC2 token endpoint; it remains responsible
for authenticating requests to that endpoint.

	opts := auth.AuthOptionsEC2{
		AuthURL:         "https://identity.example.com/v3",
		ServiceProvider: serviceProvider,
		Auth: auth.EC2TokenOpts{
			Access:  "access-key",
			Secret:  "secret-key",
			Host:    "ec2.example.com",
			Path:    "/",
			Verb:    "GET",
			Params: map[string]string{
				"SignatureVersion": "2",
				"SignatureMethod":  auth.EC2CredentialsHmacSha256V2,
			},
		},
	}
	result, err := opts.Authenticate(context.Background(), nil)
	if err != nil {
		panic(err)
	}

# Reauthentication

Credential options supporting AllowReauth default to false. Set it to true to
permit reuse of credentials for automatic reauthentication when configuring a
provider client. Authenticate itself only performs one authentication request.
TOTP and multifactor authentication do not permit automatic reauthentication;
a one-time passcode cannot be reused safely. Reauthenticating with an existing
token cannot extend that source token's validity.
*/
package auth
