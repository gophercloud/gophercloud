package identity

import (
	"github.com/mitchellh/mapstructure"
)

type TenantDesc struct {
	Id   string
	Name string
}

type TokenDesc struct {
	Id_      string `mapstructure:"Id"`
	Expires_ string `mapstructure:"Expires"`
	Tenant   TenantDesc
}

func Token(m AuthResults) (*TokenDesc, error) {
	accessMap := m["access"].(map[string]interface{})
	tokenMap := accessMap["token"].(map[string]interface{})
	td := &TokenDesc{}
	err := mapstructure.Decode(tokenMap, td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (td *TokenDesc) Id() string {
	return td.Id_
}

func (td *TokenDesc) Expires() string {
	return td.Expires_
}

func (td *TokenDesc) TenantId() string {
	return td.Tenant.Id
}

func (td *TokenDesc) TenantName() string {
	return td.Tenant.Name
}
