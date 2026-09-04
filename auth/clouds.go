package auth

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
)

// v3totpAuthType is the clouds.AuthType string value for TOTP auth.
// clouds.AuthType exports no constant for it (clouds.yaml has no
// first-class TOTP concept), but a clouds.yaml entry can still legally set
// auth_type: v3totp, so the dispatch below must still recognize the raw
// string.
const v3totpAuthType clouds.AuthType = "v3totp"

// CloudOption overrides a single credential field read from a
// clouds.Cloud's AuthInfo before mechanism selection and scope resolution
// run. Use it to supply values that don't belong in a static clouds.yaml
// file (a TOTP passcode, via WithPasscode) or to override a value at
// runtime.
type CloudOption func(*cloudOptions)

type cloudOptions struct {
	username                    string
	userID                      string
	password                    string
	token                       string
	passcode                    string
	domainID                    string
	domainName                  string
	projectID                   string
	projectName                 string
	applicationCredentialID     string
	applicationCredentialName   string
	applicationCredentialSecret string
	scope                       *Scope
}

// WithUsername overrides the username read from clouds.Cloud.AuthInfo.
func WithUsername(username string) CloudOption {
	return func(co *cloudOptions) { co.username = username }
}

// WithUserID overrides the user ID read from clouds.Cloud.AuthInfo.
func WithUserID(userID string) CloudOption {
	return func(co *cloudOptions) { co.userID = userID }
}

// WithPassword overrides the password read from clouds.Cloud.AuthInfo.
func WithPassword(password string) CloudOption {
	return func(co *cloudOptions) { co.password = password }
}

// WithToken overrides the token read from clouds.Cloud.AuthInfo.
func WithToken(token string) CloudOption {
	return func(co *cloudOptions) { co.token = token }
}

// WithPasscode supplies a one-time TOTP passcode. clouds.AuthInfo has no
// passcode field, so this is the only way to authenticate with v3totp
// against a clouds.yaml-sourced cloud.
func WithPasscode(passcode string) CloudOption {
	return func(co *cloudOptions) { co.passcode = passcode }
}

// WithDomainID overrides the domain ID read from clouds.Cloud.AuthInfo.
func WithDomainID(domainID string) CloudOption {
	return func(co *cloudOptions) { co.domainID = domainID }
}

// WithDomainName overrides the domain name read from clouds.Cloud.AuthInfo.
func WithDomainName(domainName string) CloudOption {
	return func(co *cloudOptions) { co.domainName = domainName }
}

// WithProjectID overrides the project ID read from clouds.Cloud.AuthInfo.
func WithProjectID(projectID string) CloudOption {
	return func(co *cloudOptions) { co.projectID = projectID }
}

// WithProjectName overrides the project name read from
// clouds.Cloud.AuthInfo.
func WithProjectName(projectName string) CloudOption {
	return func(co *cloudOptions) { co.projectName = projectName }
}

// WithApplicationCredentialID overrides the application credential ID read
// from clouds.Cloud.AuthInfo.
func WithApplicationCredentialID(id string) CloudOption {
	return func(co *cloudOptions) { co.applicationCredentialID = id }
}

// WithApplicationCredentialName overrides the application credential name
// read from clouds.Cloud.AuthInfo.
func WithApplicationCredentialName(name string) CloudOption {
	return func(co *cloudOptions) { co.applicationCredentialName = name }
}

// WithApplicationCredentialSecret overrides the application credential
// secret read from clouds.Cloud.AuthInfo.
func WithApplicationCredentialSecret(secret string) CloudOption {
	return func(co *cloudOptions) { co.applicationCredentialSecret = secret }
}

// WithScope bypasses scope resolution entirely; scope is used verbatim.
func WithScope(scope *Scope) CloudOption {
	return func(co *cloudOptions) { co.scope = scope }
}

// resolvedAuthInfo is clouds.Cloud.AuthInfo (nil-normalized) merged with
// CloudOption overrides, with the user-domain/project-domain split and
// scope already computed.
type resolvedAuthInfo struct {
	username                    string
	userID                      string
	password                    string
	token                       string
	passcode                    string
	projectID                   string
	projectName                 string
	applicationCredentialID     string
	applicationCredentialName   string
	applicationCredentialSecret string
	userDomainID                string
	userDomainName              string
	projectDomainID             string
	projectDomainName           string
	trustID                     string
	systemScope                 bool
	scope                       *Scope
}

func resolve(cloud clouds.Cloud, opts []CloudOption) resolvedAuthInfo {
	var ai clouds.AuthInfo
	if cloud.AuthInfo != nil {
		ai = *cloud.AuthInfo
	}

	var co cloudOptions
	for _, opt := range opts {
		opt(&co)
	}

	domainID := coalesce(co.domainID, ai.DomainID)
	domainName := coalesce(co.domainName, ai.DomainName)

	userDomainID := coalesce(ai.UserDomainID, domainID)
	userDomainName := coalesce(ai.UserDomainName, domainName)
	if userDomainID == "" && userDomainName == "" {
		userDomainID = ai.DefaultDomain
	}

	projectDomainID := coalesce(ai.ProjectDomainID, domainID)
	projectDomainName := coalesce(ai.ProjectDomainName, domainName)
	if projectDomainID == "" && projectDomainName == "" {
		projectDomainID = ai.DefaultDomain
	}

	return resolvedAuthInfo{
		username:                    coalesce(co.username, ai.Username),
		userID:                      coalesce(co.userID, ai.UserID),
		password:                    coalesce(co.password, ai.Password),
		token:                       coalesce(co.token, ai.Token),
		passcode:                    co.passcode,
		projectID:                   coalesce(co.projectID, ai.ProjectID),
		projectName:                 coalesce(co.projectName, ai.ProjectName),
		applicationCredentialID:     coalesce(co.applicationCredentialID, ai.ApplicationCredentialID),
		applicationCredentialName:   coalesce(co.applicationCredentialName, ai.ApplicationCredentialName),
		applicationCredentialSecret: coalesce(co.applicationCredentialSecret, ai.ApplicationCredentialSecret),
		userDomainID:                userDomainID,
		userDomainName:              userDomainName,
		projectDomainID:             projectDomainID,
		projectDomainName:           projectDomainName,
		trustID:                     ai.TrustID,
		systemScope:                 ai.SystemScope != "",
		scope:                       co.scope,
	}
}

// authScope returns the CloudOption-supplied scope verbatim if one was
// given via WithScope, otherwise the scope computed from the resolved
// cloud fields.
func (r resolvedAuthInfo) authScope() *Scope {
	if r.scope != nil {
		return r.scope
	}
	return &Scope{
		ProjectDomainID:   r.projectDomainID,
		ProjectDomainName: r.projectDomainName,
		ProjectID:         r.projectID,
		ProjectName:       r.projectName,
		System:            r.systemScope,
		TrustID:           r.trustID,
	}
}

func coalesce(items ...string) string {
	for _, item := range items {
		if item != "" {
			return item
		}
	}
	return ""
}

// AuthOptionsFromCloud builds an Authenticator from a clouds.Cloud (as
// parsed by openstack/config/clouds.Parse), the same way AuthOptionsFromEnv
// builds one from the process environment. Identity v2 is used when
// cloud.AuthType is an explicit v2 auth type, or when cloud.AuthType is
// empty/version-agnostic and cloud.IdentityAPIVersion == "2.0"; v3 is used
// otherwise.
func AuthOptionsFromCloud(cloud clouds.Cloud, opts ...CloudOption) (Authenticator, error) {
	switch cloud.AuthType {
	case clouds.AuthV2Password, clouds.AuthV2Token:
		return AuthOptionsFromCloudV2(cloud, opts...)
	case clouds.AuthV3Password, v3totpAuthType, clouds.AuthV3Token, clouds.AuthV3ApplicationCredential:
		return AuthOptionsFromCloudV3(cloud, opts...)
	}

	if cloud.IdentityAPIVersion == "2.0" {
		return AuthOptionsFromCloudV2(cloud, opts...)
	}
	return AuthOptionsFromCloudV3(cloud, opts...)
}

// AuthOptionsFromCloudV2 builds an *AuthOptionsV2 from a clouds.Cloud,
// ignoring cloud.IdentityAPIVersion.
func AuthOptionsFromCloudV2(cloud clouds.Cloud, opts ...CloudOption) (*AuthOptionsV2, error) {
	if cloud.AuthInfo == nil || cloud.AuthInfo.AuthURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}

	r := resolve(cloud, opts)

	authType := cloud.AuthType
	if authType == "" {
		if r.password != "" {
			authType = clouds.AuthV2Password
		} else if r.token != "" {
			authType = clouds.AuthV2Token
		}
	}

	var authOpts AuthOptionsBuilderV2
	switch authType {
	case clouds.AuthV2Password, clouds.AuthPassword:
		authOpts = V2PasswordOpts{
			Username:   r.username,
			Password:   r.password,
			TenantID:   r.projectID,
			TenantName: r.projectName,
		}
	case clouds.AuthV2Token, clouds.AuthToken:
		authOpts = V2TokenOpts{
			Token:      r.token,
			TenantID:   r.projectID,
			TenantName: r.projectName,
		}
	case "":
		return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	return &AuthOptionsV2{
		AuthURL: cloud.AuthInfo.AuthURL,
		Auth:    authOpts,
	}, nil
}

// AuthOptionsFromCloudV3 builds an *AuthOptionsV3 from a clouds.Cloud,
// ignoring cloud.IdentityAPIVersion.
func AuthOptionsFromCloudV3(cloud clouds.Cloud, opts ...CloudOption) (*AuthOptionsV3, error) {
	if cloud.AuthInfo == nil || cloud.AuthInfo.AuthURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}

	r := resolve(cloud, opts)

	authType := cloud.AuthType
	if authType == "" {
		switch {
		case r.password != "":
			authType = clouds.AuthV3Password
		case r.passcode != "":
			authType = v3totpAuthType
		case r.token != "":
			authType = clouds.AuthV3Token
		case r.applicationCredentialID != "" || r.applicationCredentialName != "":
			authType = clouds.AuthV3ApplicationCredential
		default:
			return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
		}
	}

	var authOpts AuthOptionsBuilderV3
	switch authType {
	case clouds.AuthV3Password, clouds.AuthPassword:
		authOpts = V3PasswordOpts{
			Username:       r.username,
			UserID:         r.userID,
			Password:       r.password,
			UserDomainID:   r.userDomainID,
			UserDomainName: r.userDomainName,
			Scope:          r.authScope(),
		}
	case v3totpAuthType:
		authOpts = V3TOTPOpts{
			Username:       r.username,
			UserID:         r.userID,
			Passcode:       r.passcode,
			UserDomainID:   r.userDomainID,
			UserDomainName: r.userDomainName,
			Scope:          r.authScope(),
		}
	case clouds.AuthV3Token, clouds.AuthToken:
		authOpts = V3TokenOpts{
			Token: r.token,
			Scope: r.authScope(),
		}
	case clouds.AuthV3ApplicationCredential:
		authOpts = V3ApplicationCredentialOpts{
			Username:                    r.username,
			UserID:                      r.userID,
			ApplicationCredentialID:     r.applicationCredentialID,
			ApplicationCredentialName:   r.applicationCredentialName,
			ApplicationCredentialSecret: r.applicationCredentialSecret,
			UserDomainID:                r.userDomainID,
			UserDomainName:              r.userDomainName,
		}
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	return &AuthOptionsV3{
		AuthURL: cloud.AuthInfo.AuthURL,
		Auth:    authOpts,
	}, nil
}
