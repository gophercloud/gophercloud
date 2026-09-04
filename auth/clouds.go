package auth

import (
	"encoding/json"
	"maps"

	"github.com/gophercloud/gophercloud/v2"
)

// CloudOption overrides a single credential in clouds.yaml.
// Use it to supply values dynamically, for example a TOTP passcode
type CloudOption func(*CloudOptions)

// CloudOptions is the resolved set of CloudOption overrides
type CloudOptions struct {
	Username                    string
	UserID                      string
	Password                    string
	Token                       string
	Passcode                    string
	DomainID                    string
	DomainName                  string
	ProjectID                   string
	ProjectName                 string
	ApplicationCredentialID     string
	ApplicationCredentialName   string
	ApplicationCredentialSecret string
	Scope                       *Scope
}

// ResolveCloudOptions applies opts in order and returns the result.
func ResolveCloudOptions(opts ...CloudOption) CloudOptions {
	var co CloudOptions
	for _, opt := range opts {
		opt(&co)
	}
	return co
}

// WithUsername overrides the username.
func WithUsername(username string) CloudOption {
	return func(co *CloudOptions) { co.Username = username }
}

// WithUserID overrides the user ID.
func WithUserID(userID string) CloudOption {
	return func(co *CloudOptions) { co.UserID = userID }
}

// WithPassword overrides the password.
func WithPassword(password string) CloudOption {
	return func(co *CloudOptions) { co.Password = password }
}

// WithToken overrides the token.
func WithToken(token string) CloudOption {
	return func(co *CloudOptions) { co.Token = token }
}

// WithPasscode supplies a one-time TOTP passcode. Static config files
// have no passcode field, so this is the only way to authenticate with
// v3totp against a config-sourced cloud.
func WithPasscode(passcode string) CloudOption {
	return func(co *CloudOptions) { co.Passcode = passcode }
}

// WithDomainID overrides the domain ID.
func WithDomainID(domainID string) CloudOption {
	return func(co *CloudOptions) { co.DomainID = domainID }
}

// WithDomainName overrides the domain name.
func WithDomainName(domainName string) CloudOption {
	return func(co *CloudOptions) { co.DomainName = domainName }
}

// WithProjectID overrides the project ID.
func WithProjectID(projectID string) CloudOption {
	return func(co *CloudOptions) { co.ProjectID = projectID }
}

// WithProjectName overrides the project name.
func WithProjectName(projectName string) CloudOption {
	return func(co *CloudOptions) { co.ProjectName = projectName }
}

// WithApplicationCredentialID overrides the application credential ID.
func WithApplicationCredentialID(id string) CloudOption {
	return func(co *CloudOptions) { co.ApplicationCredentialID = id }
}

// WithApplicationCredentialName overrides the application credential name.
func WithApplicationCredentialName(name string) CloudOption {
	return func(co *CloudOptions) { co.ApplicationCredentialName = name }
}

// WithApplicationCredentialSecret overrides the application credential
// secret.
func WithApplicationCredentialSecret(secret string) CloudOption {
	return func(co *CloudOptions) { co.ApplicationCredentialSecret = secret }
}

// WithScope bypasses scope resolution entirely; scope is used verbatim.
func WithScope(scope *Scope) CloudOption {
	return func(co *CloudOptions) { co.Scope = scope }
}

// Interface used to get data defined in the clouds package
type CloudSource interface {
	GetAuthType() AuthType
	GetIdentityAPIVersion() string
	GetAuthData() map[string]any
}

func AuthOptionsFromCloud(c CloudSource, opts ...CloudOption) (Authenticator, error) {
	switch c.GetAuthType() {
	// TODO(danchild): what do we do with AuthToken and AuthPassword (not v2 or v3) previously defined in clouds package ??
	case AuthV2Password, AuthV2Token:
		return AuthOptionsFromCloudV2(c, opts...)
	case AuthV3Password, AuthV3Totp, AuthV3Token, AuthV3ApplicationCredential:
		return AuthOptionsFromCloudV3(c, opts...)
	}

	// If auth type isn't provided, default to identity v3
	if c.GetIdentityAPIVersion() == "2.0" {
		return AuthOptionsFromCloudV2(c, opts...)
	}

	return AuthOptionsFromCloudV3(c, opts...)
}

func AuthOptionsFromCloudV2(c CloudSource, cloudOpts ...CloudOption) (*AuthOptionsV2, error) {
	// TODO(danchild): check implementation of AuthOptionsFromCloudV2 for parity
	m, authURL := mergedAuth(c, cloudOpts)
	if authURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}

	authType := c.GetAuthType()
	if authType == "" {
		if str(m, "password") != "" {
			authType = AuthV2Password
		} else if str(m, "token") != "" {
			authType = AuthV2Token
		}
	}

	var opts AuthOptionsBuilderV2

	switch authType {
	case AuthV2Password, AuthPassword:
		var o V2PasswordOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		opts = o
	case AuthV2Token, AuthToken:
		var o V2TokenOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		opts = o
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	ao := &AuthOptionsV2{
		AuthURL: authURL,
		Auth:    opts,
	}

	return ao, nil
}

func AuthOptionsFromCloudV3(c CloudSource, cloudOpts ...CloudOption) (*AuthOptionsV3, error) {
	// TODO(danchild): check implementation of AuthOptionsFromCloudV2 for parity
	m, authURL := mergedAuth(c, cloudOpts)
	if authURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}

	authType := c.GetAuthType()
	if authType == "" {
		switch {
		case str(m, "password") != "":
			authType = AuthV3Password
		case str(m, "passcode") != "":
			authType = AuthV3Totp
		case str(m, "token") != "":
			authType = AuthV3Token
		case str(m, "application_credential_id") != "" || str(m, "application_credential_name") != "":
			authType = AuthV3ApplicationCredential
		default:
			return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
		}
	}

	userDomainID, userDomainName, projectDomainID, projectDomainName := resolveDomains(m)
	m["user_domain_id"] = userDomainID
	m["user_domain_name"] = userDomainName

	co := ResolveCloudOptions(cloudOpts...)
	scope := co.Scope
	if scope == nil {
		scope = &Scope{
			ProjectDomainID:   projectDomainID,
			ProjectDomainName: projectDomainName,
			ProjectID:         str(m, "project_id"),
			ProjectName:       str(m, "project_name"),
			System:            str(m, "system_scope") != "",
			TrustID:           str(m, "trust_id"),
		}
	}

	var opts AuthOptionsBuilderV3

	switch authType {
	case AuthV3Password, AuthPassword:
		var o V3PasswordOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		opts = o
	case AuthV3Totp:
		var o V3TOTPOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		opts = o
	case AuthV3Token, AuthToken:
		var o V3TokenOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		opts = o
	case AuthV3ApplicationCredential:
		var o V3ApplicationCredentialOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		opts = o
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	ao := &AuthOptionsV3{
		AuthURL: authURL,
		Auth:    opts,
	}

	return ao, nil
}

// mergedAuth clones c.GetAuth() and overlays opts onto it, keyed the same
// way clouds.yaml's auth: section is, and returns the merged map plus the
// resolved auth_url.
func mergedAuth(c CloudSource, opts []CloudOption) (map[string]any, string) {
	src := c.GetAuthData()
	m := make(map[string]any, len(src)+1)
	maps.Copy(m, src)

	co := ResolveCloudOptions(opts...)
	setIfNotEmpty(m, "username", co.Username)
	setIfNotEmpty(m, "user_id", co.UserID)
	setIfNotEmpty(m, "password", co.Password)
	setIfNotEmpty(m, "token", co.Token)
	setIfNotEmpty(m, "passcode", co.Passcode)
	setIfNotEmpty(m, "domain_id", co.DomainID)
	setIfNotEmpty(m, "domain_name", co.DomainName)
	setIfNotEmpty(m, "project_id", co.ProjectID)
	setIfNotEmpty(m, "project_name", co.ProjectName)
	setIfNotEmpty(m, "application_credential_id", co.ApplicationCredentialID)
	setIfNotEmpty(m, "application_credential_name", co.ApplicationCredentialName)
	setIfNotEmpty(m, "application_credential_secret", co.ApplicationCredentialSecret)

	return m, str(m, "auth_url")
}

// resolveDomains applies clouds.yaml's user/project domain fallback:
// user_domain_id/name falls back to domain_id/name, then to
// default_domain if neither user_domain_id nor user_domain_name is set;
// project_domain_id/name falls back the same way.
func resolveDomains(m map[string]any) (userDomainID, userDomainName, projectDomainID, projectDomainName string) {
	domainID, domainName := str(m, "domain_id"), str(m, "domain_name")

	userDomainID = coalesce(str(m, "user_domain_id"), domainID)
	userDomainName = coalesce(str(m, "user_domain_name"), domainName)
	if userDomainID == "" && userDomainName == "" {
		userDomainID = str(m, "default_domain")
	}

	projectDomainID = coalesce(str(m, "project_domain_id"), domainID)
	projectDomainName = coalesce(str(m, "project_domain_name"), domainName)
	if projectDomainID == "" && projectDomainName == "" {
		projectDomainID = str(m, "default_domain")
	}

	return userDomainID, userDomainName, projectDomainID, projectDomainName
}

func setIfNotEmpty(m map[string]any, key, value string) {
	if value != "" {
		m[key] = value
	}
}

func str(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

func coalesce(items ...string) string {
	for _, item := range items {
		if item != "" {
			return item
		}
	}
	return ""
}

// decode marshals m to JSON and unmarshals it into target, relying on
// target's json tags to pick out the fields it cares about (extra keys
// in m that target has no matching field for are ignored).
func decode(m map[string]any, target AuthOptionsBuilder) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
