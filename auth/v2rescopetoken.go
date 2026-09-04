package auth

import "github.com/gophercloud/gophercloud/v2"

// V2RescopeTokenOpts takes an existing token and requests a new one scoped
// to a different tenant.
type V2RescopeTokenOpts struct {
	// Token is the existing token to rescope.
	Token string

	// TenantID is the ID of the tenant to scope the new token to. At least
	// one of TenantID/TenantName is required.
	TenantID string

	// TenantName is the name of the tenant to scope the new token to. At
	// least one of TenantID/TenantName is required.
	TenantName string
}

func (opts V2RescopeTokenOpts) ToAuthBody() (map[string]map[string]any, error) {
	type tokenCredentials struct {
		ID string `json:"id" required:"true"`
	}

	type tokenReq struct {
		TenantID   string `json:"tenantId,omitempty"`
		TenantName string `json:"tenantName,omitempty"`

		TokenCredentials tokenCredentials `json:"tokenCredentials"`
	}

	if opts.Token == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "Token"}
	}

	if opts.TenantID == "" && opts.TenantName == "" {
		return nil, gophercloud.ErrScopeEmpty{}
	}

	req := tokenReq{
		TenantID:   opts.TenantID,
		TenantName: opts.TenantName,
		TokenCredentials: tokenCredentials{
			ID: opts.Token,
		},
	}

	b, err := gophercloud.BuildRequestBody(req, "")
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]any{
		V2Token.toAuthMethod(): b,
	}

	return result, nil
}

// CanReauth returns false: an automatic reauth-on-401 would retry with the
// same original token and tenant request that just failed, which cannot
// succeed.
func (opts V2RescopeTokenOpts) CanReauth() bool {
	return false
}
