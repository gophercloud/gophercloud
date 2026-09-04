package auth

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
