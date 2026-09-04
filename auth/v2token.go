package auth

import "github.com/gophercloud/gophercloud/v2"

type V2TokenOpts struct {
	Username   string
	Token      string
	TenantID   string
	TenantName string
}

func (opts V2TokenOpts) ToAuthBody() (map[string]map[string]any, error) {
	type tokenCredentials struct {
		ID string `json:"id" required:"true"`
	}

	type tokenReq struct {
		// The TenantID and TenantName fields are optional for the Identity V2 API
		// Some providers allow you to specify a TenantName instead of the TenantId.
		// Some require both. Your provider's authentication policies will determine
		// how these fields influence authentication.
		TenantID   string `json:"tenantId,omitempty"`
		TenantName string `json:"tenantName,omitempty"`

		TokenCredentials tokenCredentials `json:"tokenCredentials"`
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

func (opts V2TokenOpts) CanReauth() bool {
	return true
}
