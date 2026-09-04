package auth

import (
	"fmt"
	"maps"

	"github.com/gophercloud/gophercloud/v2"
)

type V3MultifactorOpts struct {
	AuthMethods []AuthOptionsBuilderV3
	Scope       *Scope
}

func (opts V3MultifactorOpts) ToAuthBody() (map[string]map[string]any, error) {
	result := make(map[string]map[string]any)
	var authTypes []AuthType
	for _, authMethod := range opts.AuthMethods {
		var authResult map[string]map[string]any
		var err error

		switch authMethod.(type) {
		case V3PasswordOpts:
			authTypes = append(authTypes, V3Password)
		case V3TOTPOpts:
			authTypes = append(authTypes, V3TOTP)
		case V3ApplicationCredentialOpts:
			authTypes = append(authTypes, V3ApplicationCredential)
		case V3TokenOpts:
			authTypes = append(authTypes, V3Token)
		default:
			return nil, gophercloud.ErrUnsupportedAuthType{AuthType: fmt.Sprintf("%T", authMethod)}
		}

		authResult, err = authMethod.ToAuthBody()
		if err != nil {
			return nil, err
		}

		maps.Copy(result, authResult)
	}

	return result, nil
}

func (opts V3MultifactorOpts) ToAuthHeaders() (map[string]any, error) {
	return nil, nil
}

func (opts V3MultifactorOpts) CanReauth() bool {
	return false
}

func (opts V3MultifactorOpts) ToAuthScope() (map[string]any, error) {
	result, err := opts.Scope.ToScopeMap()

	if err != nil {
		return nil, err
	}
	return result, nil
}
