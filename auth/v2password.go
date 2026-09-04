package auth

import "github.com/gophercloud/gophercloud/v2"

type V2PasswordOpts struct {
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	TenantID   string `json:"project_id,omitempty"`
	TenantName string `json:"project_name,omitempty"`

	// AllowReauth grants permission for Gophercloud to cache these
	// credentials in memory and automatically re-authenticate when the
	// token expires. If false, credentials are not cached and
	// re-authentication is not possible. Defaults to false.
	AllowReauth bool `json:"allow_reauth,omitempty"`
}

func (opts V2PasswordOpts) ToAuthBody() (map[string]map[string]any, error) {
	type passwordCredentials struct {
		Username string `json:"username" required:"true"`
		Password string `json:"password" required:"true"`
	}

	type passwordReq struct {
		// The TenantID and TenantName fields are optional for the Identity V2 API
		// Some providers allow you to specify a TenantName instead of the TenantId.
		// Some require both. Your provider's authentication policies will determine
		// how these fields influence authentication.
		TenantID   string `json:"tenantId,omitempty"`
		TenantName string `json:"tenantName,omitempty"`

		PasswordCredentials passwordCredentials `json:"passwordCredentials"`
	}

	req := passwordReq{
		TenantID:   opts.TenantID,
		TenantName: opts.TenantName,
		PasswordCredentials: passwordCredentials{
			Username: opts.Username,
			Password: opts.Password,
		},
	}

	b, err := gophercloud.BuildRequestBody(req, "")
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]any{
		AuthV2Password.toAuthMethod(): b,
	}

	return result, nil
}

func (opts V2PasswordOpts) CanReauth() bool {
	return opts.AllowReauth
}
