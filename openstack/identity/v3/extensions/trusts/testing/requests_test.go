package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func TestCreateUserIDPasswordTrustID(t *testing.T) {
	HandleCreateTokenWithTrustID(t, tokens.AuthOptions{UserID: "me", Password: "squirrel!"}, trusts.ScopeExt{TrustID: "de0945a"}, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": { "id": "me", "password": "squirrel!" }
					}
				},
        "scope": {
            "OS-TRUST:trust": {
                "id": "de0945a"
            }
        }
			}
		}
	`)
}
