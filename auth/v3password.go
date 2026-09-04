package auth

import "github.com/gophercloud/gophercloud/v2"

type V3PasswordOpts struct {
	Username       string `json:"username,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	Password       string `json:"password,omitempty"`
	UserDomainID   string `json:"user_domain_id,omitempty"`
	UserDomainName string `json:"user_domain_name,omitempty"`
	Scope          *Scope `json:"-"`

	// AllowReauth grants permission for Gophercloud to cache these
	// credentials in memory and automatically re-authenticate when the
	// token expires. If false, credentials are not cached and
	// re-authentication is not possible. Defaults to false.
	AllowReauth bool `json:"allow_reauth,omitempty"`
}

func (opts V3PasswordOpts) ToAuthBody() (map[string]map[string]any, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type userReq struct {
		ID       *string    `json:"id,omitempty"`
		Name     *string    `json:"name,omitempty"`
		Password string     `json:"password"`
		Domain   *domainReq `json:"domain,omitempty"`
	}

	type passwordReq struct {
		User userReq `json:"user"`
	}

	req := passwordReq{
		User: userReq{
			Password: opts.Password,
		},
	}

	if opts.Password == "" {
		// A password must be specified.
		return nil, gophercloud.ErrMissingPassword{}
	}

	// Exactly one of Username and UserID must be specified
	if opts.Username == "" && opts.UserID == "" {
		return nil, gophercloud.ErrUsernameOrUserID{}
	} else if opts.Username != "" && opts.UserID != "" {
		return nil, gophercloud.ErrUsernameOrUserID{}
	}

	if opts.Username != "" {
		// Exactly one of UserDomainID or UserDomainName must be specified
		if opts.UserDomainID == "" && opts.UserDomainName == "" {
			return nil, gophercloud.ErrDomainIDOrDomainName{}
		} else if opts.UserDomainID != "" && opts.UserDomainName != "" {
			return nil, gophercloud.ErrDomainIDOrDomainName{}
		}

		var domain *domainReq

		if opts.UserDomainID != "" {
			domain = &domainReq{ID: &opts.UserDomainID}
		} else { // opts.UserDomainName != ""
			domain = &domainReq{Name: &opts.UserDomainName}
		}

		req.User.Name = &opts.Username
		req.User.Domain = domain
	} else { // opts.UserID != ""
		// None of UserDomainID or UserDomainName may be specified
		if opts.UserDomainID != "" {
			return nil, gophercloud.ErrDomainIDWithUserID{}
		} else if opts.UserDomainName != "" {
			return nil, gophercloud.ErrDomainNameWithUserID{}
		}
		req.User.ID = &opts.UserID
	}

	b, err := gophercloud.BuildRequestBody(req, "")
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]any{
		AuthV3Password.toAuthMethod(): b,
	}

	return result, nil
}

func (opts V3PasswordOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3PasswordOpts) CanReauth() bool {
	return opts.AllowReauth
}

func (opts V3PasswordOpts) ToAuthScope() (map[string]any, error) {
	result, err := opts.Scope.ToScopeMap()

	if err != nil {
		return nil, err
	}

	return result, nil
}
