package tokens

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// Scope allows a created token to be limited to a specific domain or project.
type Scope struct {
	ProjectID   string
	ProjectName string
	DomainID    string
	DomainName  string
}

// Create authenticates and either generates a new token, or changes the Scope of an existing token.
func Create(c *gophercloud.ServiceClient, scope *Scope) (gophercloud.AuthResults, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type projectReq struct {
		Domain *domainReq `json:"domain,omitempty"`
		Name   *string    `json:"name,omitempty"`
		ID     *string    `json:"id,omitempty"`
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

	type tokenReq struct {
		ID string `json:"id"`
	}

	type identityReq struct {
		Methods  []string     `json:"methods"`
		Password *passwordReq `json:"password,omitempty"`
		Token    *tokenReq    `json:"token,omitempty"`
	}

	type scopeReq struct {
		Domain  *domainReq  `json:"domain,omitempty"`
		Project *projectReq `json:"project,omitempty"`
	}

	type authReq struct {
		Identity identityReq `json:"identity"`
		Scope    *scopeReq   `json:"scope,omitempty"`
	}

	type request struct {
		Auth authReq `json:"auth"`
	}

	ao := c.Options

	// Populate the request structure based on the provided arguments. Create and return an error
	// if insufficient or incompatible information is present.
	var req request

	// Test first for unrecognized arguments.
	if ao.APIKey != "" {
		return nil, ErrAPIKeyProvided
	}
	if ao.TenantID != "" {
		return nil, ErrTenantIDProvided
	}
	if ao.TenantName != "" {
		return nil, ErrTenantNameProvided
	}

	if ao.Password == "" {
		if c.TokenID != "" {
			// Because we aren't using password authentication, it's an error to also provide any of the user-based authentication
			// parameters.
			if ao.Username != "" {
				return nil, ErrUsernameWithToken
			}
			if ao.UserID != "" {
				return nil, ErrUserIDWithToken
			}
			if ao.DomainID != "" {
				return nil, ErrDomainIDWithToken
			}
			if ao.DomainName != "" {
				return nil, ErrDomainNameWithToken
			}

			// Configure the request for Token authentication.
			req.Auth.Identity.Methods = []string{"token"}
			req.Auth.Identity.Token = &tokenReq{
				ID: c.TokenID,
			}
		} else {
			// If no password or token ID are available, authentication can't continue.
			return nil, ErrMissingPassword
		}
	} else {
		// Password authentication.
		req.Auth.Identity.Methods = []string{"password"}

		// At least one of Username and UserID must be specified.
		if ao.Username == "" && ao.UserID == "" {
			return nil, ErrUsernameOrUserID
		}

		if ao.Username != "" {
			// If Username is provided, UserID may not be provided.
			if ao.UserID != "" {
				return nil, ErrUsernameOrUserID
			}

			// Either DomainID or DomainName must also be specified.
			if ao.DomainID == "" && ao.DomainName == "" {
				return nil, ErrDomainIDOrDomainName
			}

			if ao.DomainID != "" {
				if ao.DomainName != "" {
					return nil, ErrDomainIDOrDomainName
				}

				// Configure the request for Username and Password authentication with a DomainID.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &ao.Username,
						Password: ao.Password,
						Domain:   &domainReq{ID: &ao.DomainID},
					},
				}
			}

			if ao.DomainName != "" {
				// Configure the request for Username and Password authentication with a DomainName.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &ao.Username,
						Password: ao.Password,
						Domain:   &domainReq{Name: &ao.DomainName},
					},
				}
			}
		}

		if ao.UserID != "" {
			// If UserID is specified, neither DomainID nor DomainName may be.
			if ao.DomainID != "" {
				return nil, ErrDomainIDWithUserID
			}
			if ao.DomainName != "" {
				return nil, ErrDomainNameWithUserID
			}

			// Configure the request for UserID and Password authentication.
			req.Auth.Identity.Password = &passwordReq{
				User: userReq{ID: &ao.UserID, Password: ao.Password},
			}
		}
	}

	// Add a "scope" element if a Scope has been provided.
	if scope != nil {
		if scope.ProjectName != "" {
			// ProjectName provided: either DomainID or DomainName must also be supplied.
			// ProjectID may not be supplied.
			if scope.DomainID == "" && scope.DomainName == "" {
				return nil, ErrScopeDomainIDOrDomainName
			}
			if scope.ProjectID != "" {
				return nil, ErrScopeProjectIDOrProjectName
			}

			if scope.DomainID != "" {
				// ProjectName + DomainID
				req.Auth.Scope = &scopeReq{
					Project: &projectReq{
						Name:   &scope.ProjectName,
						Domain: &domainReq{ID: &scope.DomainID},
					},
				}
			}

			if scope.DomainName != "" {
				// ProjectName + DomainName
				req.Auth.Scope = &scopeReq{
					Project: &projectReq{
						Name:   &scope.ProjectName,
						Domain: &domainReq{Name: &scope.DomainName},
					},
				}
			}
		} else if scope.ProjectID != "" {
			// ProjectID provided. ProjectName, DomainID, and DomainName may not be provided.
			if scope.ProjectName != "" {
				return nil, ErrScopeProjectIDOrProjectName
			}
			if scope.DomainID != "" {
				return nil, ErrScopeProjectIDAlone
			}
			if scope.DomainName != "" {
				return nil, ErrScopeProjectIDAlone
			}

			// ProjectID
			req.Auth.Scope = &scopeReq{
				Project: &projectReq{ID: &scope.ProjectID},
			}
		} else if scope.DomainID != "" {
			// DomainID provided. ProjectID, ProjectName, and DomainName may not be provided.
			if scope.DomainName != "" {
				return nil, ErrScopeDomainIDOrDomainName
			}

			// DomainID
			req.Auth.Scope = &scopeReq{
				Domain: &domainReq{ID: &scope.DomainID},
			}
		} else if scope.DomainName != "" {
			return nil, ErrScopeDomainName
		} else {
			return nil, ErrScopeEmpty
		}
	}

	var resp TokenCreateResult
	perigee.Post(getTokenURL(c), perigee.Options{
		ReqBody: &req,
		Results: &resp,
		OkCodes: []int{201},
	})
	return &resp, nil
}
