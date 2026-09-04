package auth

import (
	"encoding/json"
	"maps"

	"github.com/gophercloud/gophercloud/v2"
)

// CloudOption overrides a single credential field before mechanism
// selection and scope resolution run. Use it to supply values that
// don't belong in a static clouds.yaml file (a TOTP passcode, via
// WithPasscode) or to override a value at runtime.
type CloudOption func(*CloudOptions)

// CloudOptions is the resolved set of CloudOption overrides. It's
// exported so packages that bridge an external config format (like
// openstack/config/clouds) into auth can read the overridden values
// without depending on auth's unexported internals.
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

// CloudSource is the auth-relevant subset of a parsed clouds.yaml cloud
// entry. It exists so this package can build an Authenticator from
// clouds.yaml-sourced data the same way AuthOptionsFromEnv builds one from
// the process environment, without importing the package that defines the
// full cloud entry (openstack/config/clouds) — avoiding an import cycle,
// since that package already imports auth for AuthType. Any type exposing
// these three accessors satisfies this interface; openstack/config/clouds.Cloud
// does so via GetAuthType/GetIdentityAPIVersion/GetAuth.
type CloudSource interface {
	// GetAuthType is the explicit auth_type of the cloud entry, or "" if
	// unset (mechanism is then inferred from GetAuth's populated fields).
	GetAuthType() AuthType

	// GetIdentityAPIVersion is the identity_api_version of the cloud
	// entry (e.g. "2.0"), or "" if unset.
	GetIdentityAPIVersion() string

	// GetAuth is the auth: section of the cloud entry verbatim, keyed the
	// same way clouds.yaml is (e.g. "username", "user_domain_id").
	GetAuth() map[string]any
}

// AuthOptionsFromCloud builds an Authenticator from a CloudSource (such as
// a parsed openstack/config/clouds.Cloud), the same way AuthOptionsFromEnv
// builds one from the process environment. Identity v2 is used when
// c.GetAuthType() is an explicit v2 auth type, or when it's
// empty/version-agnostic and c.GetIdentityAPIVersion() == "2.0"; v3 is
// used otherwise.
func AuthOptionsFromCloud(c CloudSource, opts ...CloudOption) (Authenticator, error) {
	switch c.GetAuthType() {
	case AuthV2Password, AuthV2Token:
		return AuthOptionsFromCloudV2(c, opts...)
	case AuthV3Password, AuthV3Totp, AuthV3Token, AuthV3ApplicationCredential:
		return AuthOptionsFromCloudV3(c, opts...)
	}

	if c.GetIdentityAPIVersion() == "2.0" {
		return AuthOptionsFromCloudV2(c, opts...)
	}
	return AuthOptionsFromCloudV3(c, opts...)
}

// AuthOptionsFromCloudV2 builds an *AuthOptionsV2 from a CloudSource,
// ignoring c.GetIdentityAPIVersion().
func AuthOptionsFromCloudV2(c CloudSource, opts ...CloudOption) (*AuthOptionsV2, error) {
	m, authURL := mergedAuth(c, opts)
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

	var authOpts AuthOptionsBuilderV2
	switch authType {
	case AuthV2Password, AuthPassword:
		var o V2PasswordOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		authOpts = o
	case AuthV2Token, AuthToken:
		var o V2TokenOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		authOpts = o
	case "":
		return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	return &AuthOptionsV2{AuthURL: authURL, Auth: authOpts}, nil
}

// AuthOptionsFromCloudV3 builds an *AuthOptionsV3 from a CloudSource,
// ignoring c.GetIdentityAPIVersion().
func AuthOptionsFromCloudV3(c CloudSource, opts ...CloudOption) (*AuthOptionsV3, error) {
	m, authURL := mergedAuth(c, opts)
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

	// user_domain_id/name and project_domain_id/name each fall back to
	// the generic domain_id/name, then to default_domain. Write the
	// resolved user-domain values back into m so the json.Unmarshal below
	// picks them up on whichever mechanism has UserDomainID/UserDomainName
	// fields; the project-domain values feed Scope below instead, since
	// Scope (not the mechanism struct) is where project domain lives.
	userDomainID, userDomainName, projectDomainID, projectDomainName := resolveDomains(m)
	m["user_domain_id"] = userDomainID
	m["user_domain_name"] = userDomainName

	co := ResolveCloudOptions(opts...)
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

	var authOpts AuthOptionsBuilderV3
	switch authType {
	case AuthV3Password, AuthPassword:
		var o V3PasswordOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		authOpts = o
	case AuthV3Totp:
		var o V3TOTPOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		authOpts = o
	case AuthV3Token, AuthToken:
		var o V3TokenOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		authOpts = o
	case AuthV3ApplicationCredential:
		var o V3ApplicationCredentialOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		authOpts = o
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	return &AuthOptionsV3{AuthURL: authURL, Auth: authOpts}, nil
}

// mergedAuth clones c.GetAuth() and overlays opts onto it, keyed the same
// way clouds.yaml's auth: section is, and returns the merged map plus the
// resolved auth_url.
func mergedAuth(c CloudSource, opts []CloudOption) (map[string]any, string) {
	src := c.GetAuth()
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
func decode(m map[string]any, target any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
