package auth

import (
	"os"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

func AuthOptionsFromEnv() (Authenticator, error) {
	// default to Keystone Identity v3
	if os.Getenv("OS_IDENTITY_API_VERSION") == "2.0" {
		return AuthOptionsFromEnvV2()
	} else {
		return AuthOptionsFromEnvV3()
	}
}

func AuthOptionsFromEnvV2() (*AuthOptionsV2, error) {
	authURL := os.Getenv("OS_AUTH_URL")
	tenantName := os.Getenv("OS_TENANT_NAME")
	tenantID := os.Getenv("OS_TENANT_ID")
	username := os.Getenv("OS_USERNAME")
	password := os.Getenv("OS_PASSWORD")
	token := os.Getenv("OS_TOKEN")

	if authURL == "" {
		return nil, gophercloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "OS_AUTH_URL",
		}
	}

	var authType AuthType // v2 doesn't have auth type; purely for consistency
	if password != "" {
		authType = AuthV2Password
	} else {
		authType = AuthV2Token
	}

	var opts AuthOptionsBuilderV2

	switch authType {
	case AuthV2Password:
		opts = V2PasswordOpts{
			Username:   username,
			Password:   password,
			TenantID:   tenantID,
			TenantName: tenantName,
		}
	case AuthV2Token:
		opts = V2TokenOpts{
			Token:      token,
			TenantID:   tenantID,
			TenantName: tenantName,
		}
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	ao := &AuthOptionsV2{
		AuthURL: authURL,
		Auth:    opts,
	}

	return ao, nil
}

func AuthOptionsFromEnvV3() (*AuthOptionsV3, error) {
	authURL := os.Getenv("OS_AUTH_URL")
	authType := AuthType(os.Getenv("OS_AUTH_TYPE"))
	authMethodsRaw := strings.Split(os.Getenv("OS_AUTH_METHODS"), ",")
	authMethods := make([]AuthType, 0)
	for _, am := range authMethodsRaw {
		authMethods = append(authMethods, AuthType(am))
	}

	if authURL == "" {
		return nil, gophercloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "OS_AUTH_URL",
		}
	}

	// If the user didn't provide an explicit auth type, try to guess
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

	scope := &Scope{
		DomainID:          os.Getenv("OS_DOMAIN_ID"),
		DomainName:        os.Getenv("OS_DOMAIN_NAME"),
		ProjectDomainID:   os.Getenv("OS_PROJECT_DOMAIN_ID"),
		ProjectDomainName: os.Getenv("OS_PROJECT_DOMAIN_NAME"),
		ProjectID:         os.Getenv("OS_PROJECT_ID"),
		ProjectName:       os.Getenv("OS_PROJECT_NAME"),
	}

	var opts AuthOptionsBuilderV3

	if authType == AuthV3MultiFactor {
		opts = V3MultifactorOpts{
			Scope: scope,
		}
		for _, authType := range authMethods {
			var authOpts = authMechanismFromType(authType, scope)
			if authOpts == nil {
				return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
			}
		}
	} else {
		opts = authMechanismFromType(authType, scope)
		if opts == nil {
			return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
		}
	}

	ao := &AuthOptionsV3{
		AuthURL: authURL,
		Auth:    opts,
	}

	return ao, nil
}

func authMechanismFromType(authType AuthType, scope *Scope) AuthOptionsBuilderV3 {
	var opts AuthOptionsBuilderV3

	switch authType {
	case AuthV3Password:
		opts = V3PasswordOpts{
			Username:       os.Getenv("OS_USERNAME"),
			UserID:         os.Getenv("OS_USERID"),
			Password:       os.Getenv("OS_PASSWORD"),
			UserDomainID:   os.Getenv("OS_USER_DOMAIN_ID"),
			UserDomainName: os.Getenv("OS_USER_DOMAIN_NAME"),
			Scope:          scope,
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
		}
	case AuthV3Token:
		opts = V3TokenOpts{
			Token: os.Getenv("OS_TOKEN"),
			Scope: scope,
		}
	case AuthV3MultiFactor: // this should never get here
	default:
		return nil
	}

	return opts
}
