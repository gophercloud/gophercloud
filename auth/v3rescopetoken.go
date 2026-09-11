package auth

import "github.com/gophercloud/gophercloud/v2"

// V3RescopeTokenOpts takes an existing token and requests a new one with a different scope
type V3RescopeTokenOpts struct {
	Token string
	Scope *Scope // required
}

func (opts V3RescopeTokenOpts) ToAuthBody() (map[string]map[string]any, error) {

	if opts.Token == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "Token"}
	}

	if opts.Scope == nil {
		return nil, gophercloud.ErrScopeEmpty{}
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

func (opts V3RescopeTokenOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3RescopeTokenOpts) CanReauth() bool {
	return false
}

func (opts V3RescopeTokenOpts) ToAuthScope() (map[string]any, error) {
	result, err := opts.Scope.ToScopeMap()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (opts V3RescopeTokenOpts) ToAuthType() AuthType {
	return AuthV3Token
}
