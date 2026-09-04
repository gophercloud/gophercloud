package auth

import "github.com/gophercloud/gophercloud/v2"

// V3RescopeTokenOpts takes an existing token and requests a new one scoped
// differently. Unlike V3TokenOpts (where Scope is optional, to support
// plain token pass-through), Scope is required here: the whole point of
// this mechanism is to change the scope.
type V3RescopeTokenOpts struct {
	Token string
	Scope *Scope // required
}

func (opts V3RescopeTokenOpts) ToAuthBody() (map[string]map[string]any, error) {
	type tokenReq struct {
		Token string `json:"id"`
	}

	if opts.Scope == nil {
		return nil, gophercloud.ErrScopeEmpty{}
	}

	req := tokenReq{
		Token: opts.Token,
	}

	b, err := gophercloud.BuildRequestBody(req, "")
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]any{
		V3Token.toAuthMethod(): b,
	}

	return result, nil
}

func (opts V3RescopeTokenOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

// CanReauth returns false: an automatic reauth-on-401 would retry with the
// same original token and scope request that just failed, which cannot
// succeed.
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
