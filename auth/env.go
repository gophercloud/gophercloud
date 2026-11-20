package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
)

func AuthOptionsFromEnv() (Authenticator, error) {
	authType := AuthType(os.Getenv("OS_AUTH_TYPE"))
	switch authType {
	case AuthV2Password, AuthV2Token:
		return AuthOptionsFromEnvV2()
	case AuthV3Password, AuthV3Totp, AuthV3Token, AuthV3ApplicationCredential, AuthV3MultiFactor, AuthV3OAuth1:
		return AuthOptionsFromEnvV3()
	case AuthPassword, AuthToken:
		authURL := os.Getenv("OS_AUTH_URL")
		caCertPath := os.Getenv("OS_CACERT")

		if authURL == "" {
			return nil, gophercloud.ErrMissingEnvironmentVariable{
				EnvironmentVariable: "OS_AUTH_URL",
			}
		}

		base, err := utils.BaseEndpoint(authURL)
		if err != nil {
			return nil, err
		}

		client := &gophercloud.ProviderClient{
			IdentityBase:     gophercloud.NormalizeURL(base),
			IdentityEndpoint: gophercloud.NormalizeURL(authURL),
		}
		if caCertPath != "" {
			if after, ok := strings.CutPrefix(caCertPath, "~/"); ok {
				home, err := os.UserHomeDir()
				if err != nil {
					return nil, fmt.Errorf("failed to resolve user home directory: %w", err)
				}
				caCertPath = path.Join(home, after)
			}

			caCert, err := os.ReadFile(caCertPath)
			if err != nil {
				return nil, fmt.Errorf("failed to open the CA cert file: %w", err)
			}
			caCertPool := x509.NewCertPool()
			if ok := caCertPool.AppendCertsFromPEM(bytes.TrimSpace(caCert)); !ok {
				return nil, fmt.Errorf("failed to parse the CA Cert from %q", caCertPath)
			}
			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.TLSClientConfig = &tls.Config{RootCAs: caCertPool}
			client.HTTPClient.Transport = transport
		}

		// TODO(danchild): drop utils.ChooseVersion and replace with utils.GetServiceVersions
		version, _, err := utils.ChooseVersion(context.Background(), client, []*utils.Version{
			{ID: "v2.0", Priority: 20, Suffix: "/v2.0/"},
			{ID: "v3", Priority: 30, Suffix: "/v3/"},
		})
		if err != nil {
			return nil, err
		}

		if version.ID == "v2.0" {
			return AuthOptionsFromEnvV2()
		}

		return AuthOptionsFromEnvV3()
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
	// Normalize user domain for all scopes
	userDomainID := os.Getenv("OS_USER_DOMAIN_ID")
	userDomainName := os.Getenv("OS_USER_DOMAIN_NAME")
	if userDomainID == "" {
		userDomainID = scope.DomainID
	}
	if userDomainName == "" {
		userDomainName = scope.DomainName
	}

	// Normalize project domain when using project scope
	if scope.ProjectID != "" || scope.ProjectName != "" {
		if scope.ProjectDomainID == "" {
			scope.ProjectDomainID = scope.DomainID
		}
		if scope.ProjectDomainName == "" {
			scope.ProjectDomainName = scope.DomainName
		}
		scope.DomainID = ""
		scope.DomainName = ""
	}

	var opts AuthOptionsBuilderV3

	switch authType {
	case AuthV3Password, AuthPassword:
		opts = V3PasswordOpts{
			Username:       os.Getenv("OS_USERNAME"),
			UserID:         os.Getenv("OS_USERID"),
			Password:       os.Getenv("OS_PASSWORD"),
			UserDomainID:   userDomainID,
			UserDomainName: userDomainName,
			Scope:          scope,
			AllowReauth:    true,
		}
	case AuthV3Totp:
		opts = V3TOTPOpts{
			Username:       os.Getenv("OS_USERNAME"),
			UserID:         os.Getenv("OS_USERID"),
			Passcode:       os.Getenv("OS_PASSCODE"),
			UserDomainID:   userDomainID,
			UserDomainName: userDomainName,
			Scope:          scope,
		}
	case AuthV3ApplicationCredential:
		opts = V3ApplicationCredentialOpts{
			Username:                    os.Getenv("OS_USERNAME"),
			UserID:                      os.Getenv("OS_USERID"),
			ApplicationCredentialID:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
			ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
			ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
			UserDomainID:                userDomainID,
			UserDomainName:              userDomainName,
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
