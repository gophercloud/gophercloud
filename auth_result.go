package gophercloud

/*
AuthResult is the result from the request that was used to obtain a provider
client's Keystone token. It is returned from ProviderClient.GetAuthResult().

The following types satisfy this interface:

    github.com/gophercloud/gophercloud/openstack/identity/v2/tokens.CreateResult
    github.com/gophercloud/gophercloud/openstack/identity/v3/tokens.CreateResult

Unfortunately, those types do not share any common methods (or rather, none
with identical type signatures), so this interface contains a bogus method
implemented by both of them to ensure that other types cannot be used with this
interface. Usage of this type might look like this:

  import (
    "github.com/gophercloud/gophercloud"
    tokens2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
    tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
  )

  func GetAuthenticatedUserID(providerClient *gophercloud.ProviderClient) (string, error) {
    r := providerClient.AuthResult()
    if r != nil {
      //ProviderClient did not use openstack.Authenticate(), e.g. because token
      //was set manually with ProviderClient.SetToken()
      return "", errors.New("no AuthResult available")
    }
    switch r := r.(type) {
    case *tokens2.CreateResult:
      u, err := r.ExtractUser()
      if err != nil {
        return "", err
      }
      return u.ID, nil
    case *tokens3.CreateResult:
      u, err := r.ExtractUser()
      if err != nil {
        return "", err
      }
      return u.ID, nil
    default:
      panic(fmt.Sprintf("got unexpected AuthResult type %t", r))
    }
  }
*/
type AuthResult interface {
	IsAnAuthResult()
}
