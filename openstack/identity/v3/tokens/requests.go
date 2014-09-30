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

func subjectTokenHeaders(c *gophercloud.ServiceClient, subjectToken string) map[string]string {
	h := c.Provider.AuthenticatedHeaders()
	h["X-Subject-Token"] = subjectToken
	return h
}

// Create authenticates and either generates a new token, or changes the Scope of an existing token.
func Create(c *gophercloud.ServiceClient, options gophercloud.AuthOptions, scope *Scope) (gophercloud.AuthResults, error) {
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

	// Populate the request structure based on the provided arguments. Create and return an error
	// if insufficient or incompatible information is present.
	var req request

	// Test first for unrecognized arguments.
	if options.APIKey != "" {
		return nil, ErrAPIKeyProvided
	}
	if options.TenantID != "" {
		return nil, ErrTenantIDProvided
	}
	if options.TenantName != "" {
		return nil, ErrTenantNameProvided
	}

	if options.Password == "" {
		if c.Provider.TokenID != "" {
			// Because we aren't using password authentication, it's an error to also provide any of the user-based authentication
			// parameters.
			if options.Username != "" {
				return nil, ErrUsernameWithToken
			}
			if options.UserID != "" {
				return nil, ErrUserIDWithToken
			}
			if options.DomainID != "" {
				return nil, ErrDomainIDWithToken
			}
			if options.DomainName != "" {
				return nil, ErrDomainNameWithToken
			}

			// Configure the request for Token authentication.
			req.Auth.Identity.Methods = []string{"token"}
			req.Auth.Identity.Token = &tokenReq{
				ID: c.Provider.TokenID,
			}
		} else {
			// If no password or token ID are available, authentication can't continue.
			return nil, ErrMissingPassword
		}
	} else {
		// Password authentication.
		req.Auth.Identity.Methods = []string{"password"}

		// At least one of Username and UserID must be specified.
		if options.Username == "" && options.UserID == "" {
			return nil, ErrUsernameOrUserID
		}

		if options.Username != "" {
			// If Username is provided, UserID may not be provided.
			if options.UserID != "" {
				return nil, ErrUsernameOrUserID
			}

			// Either DomainID or DomainName must also be specified.
			if options.DomainID == "" && options.DomainName == "" {
				return nil, ErrDomainIDOrDomainName
			}

			if options.DomainID != "" {
				if options.DomainName != "" {
					return nil, ErrDomainIDOrDomainName
				}

				// Configure the request for Username and Password authentication with a DomainID.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &options.Username,
						Password: options.Password,
						Domain:   &domainReq{ID: &options.DomainID},
					},
				}
			}

			if options.DomainName != "" {
				// Configure the request for Username and Password authentication with a DomainName.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &options.Username,
						Password: options.Password,
						Domain:   &domainReq{Name: &options.DomainName},
					},
				}
			}
		}

		if options.UserID != "" {
			// If UserID is specified, neither DomainID nor DomainName may be.
			if options.DomainID != "" {
				return nil, ErrDomainIDWithUserID
			}
			if options.DomainName != "" {
				return nil, ErrDomainNameWithUserID
			}

			// Configure the request for UserID and Password authentication.
			req.Auth.Identity.Password = &passwordReq{
				User: userReq{ID: &options.UserID, Password: options.Password},
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

	var result TokenCreateResult
	response, err := perigee.Request("POST", tokenURL(c), perigee.Options{
		ReqBody: &req,
		Results: &result.response,
		OkCodes: []int{201},
	})
	if err != nil {
		return nil, err
	}

	// Extract the token ID from the response, if present.
	result.tokenID = response.HttpResponse.Header.Get("X-Subject-Token")

	return &result, nil
}

// Get validates and retrieves information about another token.
func Get(c *gophercloud.ServiceClient, token string) (*TokenCreateResult, error) {
	var result TokenCreateResult

	response, err := perigee.Request("GET", tokenURL(c), perigee.Options{
		MoreHeaders: subjectTokenHeaders(c, token),
		Results:     &result.response,
		OkCodes:     []int{200, 203},
	})

	if err != nil {
		return nil, err
	}

	// Extract the token ID from the response, if present.
	result.tokenID = response.HttpResponse.Header.Get("X-Subject-Token")

	return &result, nil
}

// Validate determines if a specified token is valid or not.
func Validate(c *gophercloud.ServiceClient, token string) (bool, error) {
	response, err := perigee.Request("HEAD", tokenURL(c), perigee.Options{
		MoreHeaders: subjectTokenHeaders(c, token),
		OkCodes:     []int{204, 404},
	})
	if err != nil {
		return false, err
	}

	return response.StatusCode == 204, nil
}

// Revoke immediately makes specified token invalid.
func Revoke(c *gophercloud.ServiceClient, token string) error {
	_, err := perigee.Request("DELETE", tokenURL(c), perigee.Options{
		MoreHeaders: subjectTokenHeaders(c, token),
		OkCodes:     []int{204},
	})
	return err
}
