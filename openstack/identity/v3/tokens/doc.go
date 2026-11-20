/*
Package tokens provides information and interaction with the token API
resource for the OpenStack Identity service.

For more information, see:
http://developer.openstack.org/api-ref-identity-v3.html#tokens-v3

Example to Create a Token From a Username and Password

	authOptions := auth.V3PasswordOpts{
		Username:       "username",
		Password:       "password",
		UserDomainName: "Default",
	}

	token, err := tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Create a Token From a Username, Password, and Domain

	authOptions := auth.V3PasswordOpts{
		Username:     "username",
		Password:     "password",
		UserDomainID: "default",
	}

	token, err := tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

	authOptions = auth.V3PasswordOpts{
		Username:       "username",
		Password:       "password",
		UserDomainName: "default",
	}

	token, err = tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Create a Token From a Token

	authOptions := auth.V3TokenOpts{
		Token: "token_id",
	}

	token, err := tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Create a Token from a Username and Password with Project ID Scope

	scope := auth.Scope{
		ProjectID: "0fe36e73809d46aeae6705c39077b1b3",
	}

	authOptions := auth.V3PasswordOpts{
		Scope:          &scope,
		Username:       "username",
		Password:       "password",
		UserDomainName: "Default",
	}

	token, err = tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Create a Token from a Username and Password with Domain ID Scope

	scope := auth.Scope{
		DomainID: "default",
	}

	authOptions := auth.V3PasswordOpts{
		Scope:          &scope,
		Username:       "username",
		Password:       "password",
		UserDomainName: "Default",
	}

	token, err = tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Create a Token from a Username and Password with Project Name Scope

	scope := auth.Scope{
		ProjectName:     "project_name",
		ProjectDomainID: "default",
	}

	authOptions := auth.V3PasswordOpts{
		Scope:          &scope,
		Username:       "username",
		Password:       "password",
		UserDomainName: "Default",
	}

	token, err = tokens.Create(context.TODO(), identityClient, authOptions).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Get a Token

	token, err := tokens.Get(context.TODO(), identityClient, "token_id", nil).ExtractToken()
	if err != nil {
		panic(err)
	}

# Example to Get a Token Created with Application Credentials Access Rules

When validating or retrieving tokens that were created using application
credentials with access rules, the OpenStack-Identity-Access-Rules header
must be sent. Without this header, Keystone will return a 404 Not Found error.

See https://docs.openstack.org/keystone/latest/user/application_credentials.html

	getOpts := tokens.GetOpts{
		AccessRulesVersion: "1.0",
	}
	token, err := tokens.Get(context.TODO(), identityClient, "token_id", getOpts).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Validate a Token

	ok, err := tokens.Validate(context.TODO(), identityClient, "token_id", nil)
	if err != nil {
		panic(err)
	}

	if ok {
		fmt.Println("Token is valid")
	}

Example to Validate a Token Created with Application Credentials Access Rules

	validateOpts := tokens.ValidateOpts{
		AccessRulesVersion: "1.0",
	}
	ok, err := tokens.Validate(context.TODO(), identityClient, "token_id", validateOpts)
	if err != nil {
		panic(err)
	}

	if ok {
		fmt.Println("Token is valid")
	}
*/
package tokens
