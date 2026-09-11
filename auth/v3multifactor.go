package auth

import (
	"maps"

	"github.com/gophercloud/gophercloud/v2"
)

type V3MultifactorOpts struct {
	AuthMethods []AuthOptionsBuilderV3
	Scope       *Scope
}

func (opts V3MultifactorOpts) ToAuthBody() (map[string]map[string]any, error) {
	if len(opts.AuthMethods) == 0 {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthMethods"}
	}

	result := make(map[string]map[string]any)
	for _, authMethod := range opts.AuthMethods {
		var authResult map[string]map[string]any
		var err error

		switch authMethod.(type) {
		case V3PasswordOpts, V3TOTPOpts, V3ApplicationCredentialOpts, V3TokenOpts:
		default:
			return nil, gophercloud.ErrUnsupportedAuthType{AuthType: string(authMethod.ToAuthType())}
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

func (opts V3MultifactorOpts) ToAuthType() AuthType {
	return AuthV3MultiFactor
}
