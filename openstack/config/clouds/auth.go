package clouds

import (
	"encoding/json"
	"maps"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

// AuthOptions builds an auth.Authenticator from a Cloud (as returned by
// Parse), the same way auth.AuthOptionsFromEnv builds one from the
// process environment. Identity v2 is used when c.AuthType is an
// explicit v2 auth type, or when c.AuthType is empty/version-agnostic
// and c.IdentityAPIVersion == "2.0"; v3 is used otherwise.
func (c Cloud) AuthOptions(opts ...auth.CloudOption) (auth.Authenticator, error) {
	switch c.AuthType {
	case auth.AuthV2Password, auth.AuthV2Token:
		return c.AuthOptionsV2(opts...)
	case auth.AuthV3Password, auth.AuthV3Totp, auth.AuthV3Token, auth.AuthV3ApplicationCredential:
		return c.AuthOptionsV3(opts...)
	}

	if c.IdentityAPIVersion == "2.0" {
		return c.AuthOptionsV2(opts...)
	}
	return c.AuthOptionsV3(opts...)
}

// AuthOptionsV2 builds an *auth.AuthOptionsV2 from a Cloud, ignoring
// c.IdentityAPIVersion.
func (c Cloud) AuthOptionsV2(opts ...auth.CloudOption) (*auth.AuthOptionsV2, error) {
	m, authURL := c.mergedAuth(opts)
	if authURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}

	authType := c.AuthType
	if authType == "" {
		if str(m, "password") != "" {
			authType = auth.AuthV2Password
		} else if str(m, "token") != "" {
			authType = auth.AuthV2Token
		}
	}

	var authOpts auth.AuthOptionsBuilderV2
	switch authType {
	case auth.AuthV2Password, auth.AuthPassword:
		var o auth.V2PasswordOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		authOpts = o
	case auth.AuthV2Token, auth.AuthToken:
		var o auth.V2TokenOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		authOpts = o
	case "":
		return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	return &auth.AuthOptionsV2{AuthURL: authURL, Auth: authOpts}, nil
}

// AuthOptionsV3 builds an *auth.AuthOptionsV3 from a Cloud, ignoring
// c.IdentityAPIVersion.
func (c Cloud) AuthOptionsV3(opts ...auth.CloudOption) (*auth.AuthOptionsV3, error) {
	m, authURL := c.mergedAuth(opts)
	if authURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}

	authType := c.AuthType
	if authType == "" {
		switch {
		case str(m, "password") != "":
			authType = auth.AuthV3Password
		case str(m, "passcode") != "":
			authType = auth.AuthV3Totp
		case str(m, "token") != "":
			authType = auth.AuthV3Token
		case str(m, "application_credential_id") != "" || str(m, "application_credential_name") != "":
			authType = auth.AuthV3ApplicationCredential
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

	co := auth.ResolveCloudOptions(opts...)
	scope := co.Scope
	if scope == nil {
		scope = &auth.Scope{
			ProjectDomainID:   projectDomainID,
			ProjectDomainName: projectDomainName,
			ProjectID:         str(m, "project_id"),
			ProjectName:       str(m, "project_name"),
			System:            str(m, "system_scope") != "",
			TrustID:           str(m, "trust_id"),
		}
	}

	var authOpts auth.AuthOptionsBuilderV3
	switch authType {
	case auth.AuthV3Password, auth.AuthPassword:
		var o auth.V3PasswordOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		authOpts = o
	case auth.AuthV3Totp:
		var o auth.V3TOTPOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		authOpts = o
	case auth.AuthV3Token, auth.AuthToken:
		var o auth.V3TokenOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		o.Scope = scope
		authOpts = o
	case auth.AuthV3ApplicationCredential:
		var o auth.V3ApplicationCredentialOpts
		if err := decode(m, &o); err != nil {
			return nil, err
		}
		authOpts = o
	default:
		return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authType)}
	}

	return &auth.AuthOptionsV3{AuthURL: authURL, Auth: authOpts}, nil
}

// mergedAuth clones c.Auth and overlays opts onto it, keyed the same way
// clouds.yaml's auth: section is, and returns the merged map plus the
// resolved auth_url.
func (c Cloud) mergedAuth(opts []auth.CloudOption) (map[string]any, string) {
	m := make(map[string]any, len(c.Auth)+1)
	maps.Copy(m, c.Auth)

	co := auth.ResolveCloudOptions(opts...)
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
