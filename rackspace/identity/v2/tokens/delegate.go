package tokens

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
)

// Create authenticates to Rackspace's identity service and attempts to acquire a Token. Rather
// than interact with this service directly, users should generally call
// rackspace.AuthenticatedClient().
func Create(client *gophercloud.ServiceClient, auth gophercloud.AuthOptions) os.CreateResult {
	if auth.APIKey != "" {
		// Authenticate with the provided API key.

		if auth.Username == "" {
			return createErr(os.ErrUsernameRequired)
		}

		var request struct {
			Auth struct {
				APIKeyCredentials struct {
					Username string `json:"username"`
					APIKey   string `json:"apiKey"`
				} `json:"RAX-KSKEY:apiKeyCredentials"`
				TenantID   string `json:"tenantId,omitempty"`
				TenantName string `json:"tenantName,omitempty"`
			} `json:"auth"`
		}

		request.Auth.APIKeyCredentials.Username = auth.Username
		request.Auth.APIKeyCredentials.APIKey = auth.APIKey
		request.Auth.TenantID = auth.TenantID
		request.Auth.TenantName = auth.TenantName

		var result os.CreateResult
		_, result.Err = perigee.Request("POST", os.CreateURL(client), perigee.Options{
			ReqBody: &request,
			Results: &result.Resp,
			OkCodes: []int{200, 203},
		})
		return result
	}

	return os.Create(client, auth)
}

func createErr(err error) os.CreateResult {
	return os.CreateResult{CommonResult: gophercloud.CommonResult{Err: err}}
}
