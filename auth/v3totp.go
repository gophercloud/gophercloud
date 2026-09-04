package auth

import "github.com/gophercloud/gophercloud/v2"

type V3TOTPOpts struct {
	Username       string `json:"username,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	Passcode       string `json:"passcode,omitempty"`
	UserDomainID   string `json:"user_domain_id,omitempty"`
	UserDomainName string `json:"user_domain_name,omitempty"`
	Scope          *Scope `json:"-"`
}

func (opts V3TOTPOpts) ToAuthBody() (map[string]map[string]any, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type userReq struct {
		ID       *string    `json:"id,omitempty"`
		Name     *string    `json:"name,omitempty"`
		Passcode string     `json:"passcode"`
		Domain   *domainReq `json:"domain,omitempty"`
	}

	type totpReq struct {
		User userReq `json:"user"`
	}

	req := totpReq{
		User: userReq{
			Passcode: opts.Passcode,
		},
	}

	if opts.Passcode == "" {
		// A passcode must be specified.
		return nil, gophercloud.ErrMissingPasscode{}
	}

	// Exactly one of Username and UserID must be specified
	if opts.Username == "" && opts.UserID == "" {
		return nil, gophercloud.ErrUsernameOrUserID{}
	} else if opts.Username != "" && opts.UserID != "" {
		return nil, gophercloud.ErrUsernameOrUserID{}
	}

	if opts.Username != "" {
		// Exactly one of DomainID or DomainName must be specified
		if opts.UserDomainID == "" && opts.UserDomainName == "" {
			return nil, gophercloud.ErrDomainIDOrDomainName{}
		} else if opts.UserDomainID != "" && opts.UserDomainName != "" {
			return nil, gophercloud.ErrDomainIDOrDomainName{}
		}

		var domain *domainReq

		if opts.UserDomainID != "" {
			domain = &domainReq{ID: &opts.UserDomainID}
		} else { // opts.DomainName != ""
			domain = &domainReq{Name: &opts.UserDomainName}
		}

		req.User.Name = &opts.Username
		req.User.Domain = domain
	} else { // opts.UserID != ""
		// None of DomainID or DomainName may be specified
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
		AuthV3Totp.toAuthMethod(): b,
	}

	return result, nil
}

func (opts V3TOTPOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3TOTPOpts) CanReauth() bool {
	return false
}

func (opts V3TOTPOpts) ToAuthScope() (map[string]any, error) {
	result, err := opts.Scope.ToScopeMap()

	if err != nil {
		return nil, err
	}

	return result, nil
}
