package auth

import "github.com/gophercloud/gophercloud/v2"

type V3TokenOpts struct {
	Token string `json:"token,omitempty"`
	Scope *Scope `json:"-"`
}

func (opts V3TokenOpts) ToAuthBody() (map[string]map[string]any, error) {
	if opts.Token == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "Token"}
	}

	type tokenReq struct {
		Token string `json:"id"`
	}

	req := tokenReq{
		Token: opts.Token,
	}

	b, err := gophercloud.BuildRequestBody(req, "")
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]any{
		AuthV3Token.toAuthMethod(): b,
	}

	return result, nil
}

func (opts V3TokenOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3TokenOpts) CanReauth() bool {
	return false
}

func (opts V3TokenOpts) ToAuthScope() (map[string]any, error) {
	result, err := opts.Scope.ToScopeMap()

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (opts V3TokenOpts) ToAuthType() AuthType {
	return AuthV3Token
}
