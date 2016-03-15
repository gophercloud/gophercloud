package gophercloud

// ScopeOptsV3 allows a created token to be limited to a specific domain or project.
type ScopeOptsV3 struct {
	ProjectID   string `json:"scope.project.id,omitempty" not:"ProjectName,DomainID,DomainName"`
	ProjectName string `json:"scope.project.name,omitempty"`
	DomainID    string `json:"scope.project.id,omitempty" not:"ProjectName,ProjectID,DomainName"`
	DomainName  string `json:"scope.project.id,omitempty"`
}

type ScopeDomainV3 struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ScopeProjectDomainV3 struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ScopeProjectV3 struct {
	Domain *ScopeProjectDomainV3 `json:"domain,omitempty"`
	Name   string                `json:"name,omitempty"`
	ID     string                `json:"id,omitempty"`
}

type ScopeV3 struct {
	Domain  *ScopeDomainV3  `json:"domain,omitempty"`
	Project *ScopeProjectV3 `json:"project,omitempty"`
}

type DomainV3 struct {
	ID   string `json:"id,omitempty" xor:"Name"`
	Name string `json:"name,omitempty" xor:"ID"`
}

type UserV3 struct {
	ID       string    `json:"id,omitempty" xor:"Name"`
	Name     string    `json:"name,omitempty" xor:"ID"`
	Password string    `json:"password" required:"true"`
	Domain   *DomainV3 `json:"domain,omitempty"`
}

type PasswordCredentialsV3 struct {
	User *UserV3 `json:"user" required:"true"`
}

type TokenCredentialsV3 struct {
	ID string `json:"id" required:"true"`
}

type IdentityCredentialsV3 struct {
	Methods             []string               `json:"methods" required:"true"`
	PasswordCredentials *PasswordCredentialsV3 `json:"password,omitempty" xor:"TokenCredentials"`
	TokenCredentials    *TokenCredentialsV3    `json:"token,omitempty" xor:"PasswordCredentials"`
}

type AuthOptionsV3 struct {
	Identity *IdentityCredentialsV3 `json:"identity" required:"true"`
	Scope    *ScopeV3               `json:"scope,omitempty"`
}

func (opts AuthOptionsV3) ToTokenV3CreateMap(scope *ScopeOptsV3) (map[string]interface{}, error) {
	if scope != nil {
		opts.Scope = &ScopeV3{
			Domain: &ScopeDomainV3{
				ID:   scope.DomainID,
				Name: scope.DomainName,
			},
			Project: &ScopeProjectV3{
				Domain: &ScopeProjectDomainV3{
					ID:   scope.DomainID,
					Name: scope.DomainName,
				},
				ID:   scope.ProjectID,
				Name: scope.ProjectName,
			},
		}
	}
	return BuildRequestBody(opts, "auth")
}
