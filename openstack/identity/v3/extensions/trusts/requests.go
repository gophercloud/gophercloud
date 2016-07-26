package trusts

import "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"

type ScopeExt struct {
	tokens.AuthScopeBuilder
	TrustID string `json:"id"`
}

func (scope ScopeExt) ToTokenV3ScopeMap() (map[string]interface{}, error) {
	b := make(map[string]interface{})
	var err error
	if scope.AuthScopeBuilder != nil {
		b, err = scope.AuthScopeBuilder.ToTokenV3ScopeMap()
		if err != nil {
			return nil, err
		}
	}

	if scope.TrustID != "" {
		b["OS-TRUST:trust"] = map[string]interface{}{
			"id": scope.TrustID,
		}
	}

	return b, nil
}
