package auth

import "github.com/gophercloud/gophercloud/v2"

type V3ApplicationCredentialOpts struct {
	Username                    string `json:"username,omitempty"`
	UserID                      string `json:"user_id,omitempty"`
	ApplicationCredentialID     string `json:"application_credential_id,omitempty"`
	ApplicationCredentialName   string `json:"application_credential_name,omitempty"`
	ApplicationCredentialSecret string `json:"application_credential_secret,omitempty"`
	UserDomainID                string `json:"user_domain_id,omitempty"`
	UserDomainName              string `json:"user_domain_name,omitempty"`

	// AllowReauth grants permission for Gophercloud to cache these
	// credentials in memory and automatically re-authenticate when the
	// token expires. If false, credentials are not cached and
	// re-authentication is not possible. Defaults to false.
	AllowReauth bool `json:"allow_reauth,omitempty"`
}

func (opts V3ApplicationCredentialOpts) ToAuthBody() (map[string]map[string]any, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type userReq struct {
		ID     *string    `json:"id,omitempty"`
		Name   *string    `json:"name,omitempty"`
		Domain *domainReq `json:"domain,omitempty"`
	}

	type applicationCredentialReq struct {
		ID     string  `json:"id,omitempty"`
		Name   string  `json:"name,omitempty"`
		User   userReq `json:"user"`
		Secret string  `json:"secret"`
	}

	req := applicationCredentialReq{
		User: userReq{},
	}

	// There are three kinds of possible application_credential requests
	//
	// 1. application_credential id + secret
	// 2. application_credential name + secret + user_id
	// 3. application_credential name + secret + username + domain_id / domain_name
	if opts.ApplicationCredentialSecret == "" {
		return nil, gophercloud.ErrAppCredMissingSecret{}
	}

	// Exactly one of ApplicationCredentialID and ApplicationCredentialName must be specified
	if opts.ApplicationCredentialID == "" && opts.ApplicationCredentialName == "" {
		return nil, gophercloud.ErrAppCredNameOrAppCredID{}
	} else if opts.ApplicationCredentialID != "" && opts.ApplicationCredentialName != "" {
		return nil, gophercloud.ErrAppCredNameOrAppCredID{}
	}

	req.ID = opts.ApplicationCredentialID
	req.Name = opts.ApplicationCredentialName
	req.Secret = opts.ApplicationCredentialSecret

	// Of the three supported application credential requests, only the
	// name-based ones (2 and 3) require a user to disambiguate application
	// credentials that share a name; the ID-based one (1) uniquely
	// identifies the application credential on its own, so Username/UserID
	// are optional there.
	userRequired := opts.ApplicationCredentialID == ""

	// At most one of Username and UserID may be specified; if the request
	// requires a user, exactly one must be specified.
	if opts.Username != "" && opts.UserID != "" {
		return nil, gophercloud.ErrUsernameOrUserID{}
	} else if userRequired && opts.Username == "" && opts.UserID == "" {
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
	} else if opts.UserID != "" {
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
		AuthV3ApplicationCredential.toAuthMethod(): b,
	}

	return result, nil
}

func (opts V3ApplicationCredentialOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3ApplicationCredentialOpts) CanReauth() bool {
	return opts.AllowReauth
}

func (opts V3ApplicationCredentialOpts) ToAuthScope() (map[string]any, error) {
	return nil, nil
}
