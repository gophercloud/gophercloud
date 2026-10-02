package auth

import (
	"os"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

func AuthOptionsFromEnv() (Authenticator, error) {
	authType := AuthType(os.Getenv("OS_AUTH_TYPE"))
	switch authType {
	case AuthV2Password, AuthV2Token:
		return AuthOptionsFromEnvV2()
	case AuthV3Password, AuthV3Totp, AuthV3Token, AuthV3ApplicationCredential, AuthV3MultiFactor, AuthV3OAuth1:
		return AuthOptionsFromEnvV3()
	case AuthPassword, AuthToken:
		v2, err := AuthOptionsFromEnvV2()
		if err != nil {
			return nil, err
		}
		v3, err := AuthOptionsFromEnvV3()
		if err != nil {
			return nil, err
		}
		return authenticator{v2: v2, v3: v3}, nil
	}

	// Fallback to identity v3 if v2 isn't set explicitly.
	switch os.Getenv("OS_IDENTITY_API_VERSION") {
	case "2", "2.0":
		return AuthOptionsFromEnvV2()
	}

	return AuthOptionsFromEnvV3()
}

func AuthOptionsFromEnvV2() (AuthOptionsV2, error) {
	authURL := os.Getenv("OS_AUTH_URL")
	tenantName := os.Getenv("OS_TENANT_NAME")
	tenantID := os.Getenv("OS_TENANT_ID")
	authType := AuthType(os.Getenv("OS_AUTH_TYPE"))

	// mimics keystone auth priority for the generic "password" or "token" auth type
	if authType == AuthPassword || authType == AuthToken {
		if projectName := os.Getenv("OS_PROJECT_NAME"); projectName != "" {
			tenantName = projectName
		}
		if projectID := os.Getenv("OS_PROJECT_ID"); projectID != "" {
			tenantID = projectID
		}
	}
	username := os.Getenv("OS_USERNAME")
	password := os.Getenv("OS_PASSWORD")
	token := os.Getenv("OS_TOKEN")

	if authURL == "" {
		return AuthOptionsV2{}, gophercloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "OS_AUTH_URL",
		}
	}

	// try to guess auth type if not provided
	if authType == "" {
		if password != "" {
			authType = AuthV2Password
		} else {
			authType = AuthV2Token
		}
	}

	var opts AuthOptionsBuilderV2

	switch authType {
	case AuthV2Password, AuthPassword:
		opts = V2PasswordOpts{
			Username:    username,
			Password:    password,
			TenantID:    tenantID,
			TenantName:  tenantName,
			AllowReauth: true,
		}
	case AuthV2Token, AuthToken:
		opts = V2TokenOpts{
			Token:       token,
			TenantID:    tenantID,
			TenantName:  tenantName,
			AllowReauth: true,
		}
	default:
		return AuthOptionsV2{}, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	ao := AuthOptionsV2{
		AuthURL: authURL,
		Auth:    opts,
	}

	return ao, nil
}

func AuthOptionsFromEnvV3() (AuthOptionsV3, error) {
	authURL := os.Getenv("OS_AUTH_URL")
	authType := AuthType(os.Getenv("OS_AUTH_TYPE"))
	systemScope := os.Getenv("OS_SYSTEM_SCOPE")
	authMethodsRaw := strings.Split(os.Getenv("OS_AUTH_METHODS"), ",")
	authMethods := make([]AuthType, 0)
	for _, am := range authMethodsRaw {
		authMethods = append(authMethods, AuthType(am))
	}

	if authURL == "" {
		return AuthOptionsV3{}, gophercloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "OS_AUTH_URL",
		}
	}

	// If the user didn't provide an explicit auth type, try to guess.
	if authType == "" {
		password := os.Getenv("OS_PASSWORD")
		passcode := os.Getenv("OS_PASSCODE")
		applicationCredentialID := os.Getenv("OS_APPLICATION_CREDENTIAL_ID")
		applicationCredentialName := os.Getenv("OS_APPLICATION_CREDENTIAL_NAME")
		token := os.Getenv("OS_TOKEN")

		if password != "" {
			authType = AuthV3Password
		} else if passcode != "" {
			authType = AuthV3Totp
		} else if token != "" {
			authType = AuthV3Token
		} else if applicationCredentialID != "" || applicationCredentialName != "" {
			authType = AuthV3ApplicationCredential
		}
	}

	if systemScope != "" && systemScope != "all" {
		return AuthOptionsV3{}, gophercloud.ErrInvalidSystemScope{}
	}

	scope := &Scope{
		DomainID:          os.Getenv("OS_DOMAIN_ID"),
		DomainName:        os.Getenv("OS_DOMAIN_NAME"),
		ProjectDomainID:   os.Getenv("OS_PROJECT_DOMAIN_ID"),
		ProjectDomainName: os.Getenv("OS_PROJECT_DOMAIN_NAME"),
		ProjectID:         os.Getenv("OS_PROJECT_ID"),
		ProjectName:       os.Getenv("OS_PROJECT_NAME"),
		System:            systemScope == "all",
		TrustID:           os.Getenv("OS_TRUST_ID"),
	}

	var opts AuthOptionsBuilderV3

	if authType == AuthV3MultiFactor {
		multifactorOpts := V3MultifactorOpts{
			Scope: scope,
		}
		for _, authType := range authMethods {
			authOpts := authMechanismFromType(authType, scope)
			if authOpts == nil {
				return AuthOptionsV3{}, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
			}
			multifactorOpts.AuthMethods = append(multifactorOpts.AuthMethods, authOpts)
		}
		opts = multifactorOpts
	} else {
		opts = authMechanismFromType(authType, scope)
		if opts == nil {
			return AuthOptionsV3{}, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
		}
	}

	ao := AuthOptionsV3{
		AuthURL: authURL,
		Auth:    opts,
	}

	return ao, nil
}

func authMechanismFromType(authType AuthType, scope *Scope) AuthOptionsBuilderV3 {
	var opts AuthOptionsBuilderV3

	switch authType {
	case AuthV3Password, AuthPassword:
		opts = V3PasswordOpts{
			Username:       os.Getenv("OS_USERNAME"),
			UserID:         os.Getenv("OS_USERID"),
			Password:       os.Getenv("OS_PASSWORD"),
			UserDomainID:   os.Getenv("OS_USER_DOMAIN_ID"),
			UserDomainName: os.Getenv("OS_USER_DOMAIN_NAME"),
			Scope:          scope,
			AllowReauth:    true,
		}
	case AuthV3Totp:
		opts = V3TOTPOpts{
			Username:       os.Getenv("OS_USERNAME"),
			UserID:         os.Getenv("OS_USERID"),
			Passcode:       os.Getenv("OS_PASSCODE"),
			UserDomainID:   os.Getenv("OS_USER_DOMAIN_ID"),
			UserDomainName: os.Getenv("OS_USER_DOMAIN_NAME"),
			Scope:          scope,
		}
	case AuthV3ApplicationCredential:
		opts = V3ApplicationCredentialOpts{
			Username:                    os.Getenv("OS_USERNAME"),
			UserID:                      os.Getenv("OS_USERID"),
			ApplicationCredentialID:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
			ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
			ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
			UserDomainID:                os.Getenv("OS_USER_DOMAIN_ID"),
			UserDomainName:              os.Getenv("OS_USER_DOMAIN_NAME"),
			AllowReauth:                 true,
		}
	case AuthV3Token, AuthToken:
		opts = V3TokenOpts{
			Token: os.Getenv("OS_TOKEN"),
			Scope: scope,
		}
	case AuthV3OAuth1:
		opts = V3OAuth1Opts{
			ConsumerKey:     os.Getenv("OS_CONSUMER_KEY"),
			ConsumerSecret:  os.Getenv("OS_CONSUMER_SECRET"),
			Token:           os.Getenv("OS_ACCESS_KEY"),
			TokenSecret:     os.Getenv("OS_ACCESS_SECRET"),
			SignatureMethod: OAuth1HMACSHA1,
			AllowReauth:     true,
		}
	case AuthV3MultiFactor: // this should never get here
	default:
		return nil
	}

	return opts
}
