// +build acceptance

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/oauth1"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestOAuth1CRUD(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	ao, err := openstack.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	authOptions := tokens.AuthOptions{
		Username:   ao.Username,
		Password:   ao.Password,
		DomainName: ao.DomainName,
		DomainID:   ao.DomainID,
		// We need a scope to get the token roles list
		Scope: tokens.Scope{
			ProjectID:   ao.TenantID,
			ProjectName: ao.TenantName,
			DomainID:    ao.DomainID,
			DomainName:  ao.DomainName,
		},
	}
	tokenRes := tokens.Create(client, &authOptions)
	token, err := tokenRes.Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)

	user, err := tokenRes.ExtractUser()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, user)

	roles, err := tokenRes.ExtractRoles()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, roles)

	project, err := tokenRes.ExtractProject()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, project)

	// Create a consumer
	createConsumerOpts := oauth1.CreateConsumerOpts{
		Description: "My test consumer",
	}
	// NOTE: secret is available only in create response
	consumer, err := oauth1.CreateConsumer(client, createConsumerOpts).Extract()
	th.AssertNoErr(t, err)

	// Delete a consumer
	defer oauth1.DeleteConsumer(client, consumer.ID)
	tools.PrintResource(t, consumer)

	th.AssertEquals(t, consumer.Description, createConsumerOpts.Description)

	// Update a consumer
	updateConsumerOpts := oauth1.UpdateConsumerOpts{
		Description: "",
	}
	updatedConsumer, err := oauth1.UpdateConsumer(client, consumer.ID, updateConsumerOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, updatedConsumer)
	th.AssertEquals(t, updatedConsumer.ID, consumer.ID)
	th.AssertEquals(t, updatedConsumer.Description, updateConsumerOpts.Description)

	// Get a consumer
	getConsumer, err := oauth1.GetConsumer(client, consumer.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getConsumer)
	th.AssertEquals(t, getConsumer.ID, consumer.ID)
	th.AssertEquals(t, getConsumer.Description, updateConsumerOpts.Description)

	// List consumers
	consumersPages, err := oauth1.ListConsumers(client).AllPages()
	th.AssertNoErr(t, err)
	consumers, err := oauth1.ExtractConsumers(consumersPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, consumers[0].ID, updatedConsumer.ID)
	th.AssertEquals(t, consumers[0].Description, updatedConsumer.Description)

	// test HMACSHA1 and PLAINTEXT signature methods
	for _, method := range []oauth1.SignatureMethod{oauth1.HMACSHA1, oauth1.PLAINTEXT} {
		oauth1MethodTest(t, client, consumer, method, user, project, roles, ao.IdentityEndpoint)
	}
}

func oauth1MethodTest(t *testing.T, client *gophercloud.ServiceClient, consumer *oauth1.Consumer, method oauth1.SignatureMethod, user *tokens.User, project *tokens.Project, roles []tokens.Role, identityEndpoint string) {
	// Request a token
	requestTokenOpts := oauth1.RequestTokenOpts{
		OAuthConsumerKey:     consumer.ID,
		OAuthConsumerSecret:  consumer.Secret,
		OAuthSignatureMethod: method,
		RequestedProjectID:   project.ID,
	}
	requestToken, err := oauth1.RequestToken(client, requestTokenOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, requestToken)

	// Authorize token
	authorizeTokenOpts := oauth1.AuthorizeTokenOpts{
		Roles: []oauth1.Role{
			// test role by ID
			{ID: roles[0].ID},
		},
	}
	if len(roles) > 1 {
		// test role by name
		authorizeTokenOpts.Roles = append(authorizeTokenOpts.Roles, oauth1.Role{Name: roles[1].Name})
	}
	authToken, err := oauth1.AuthorizeToken(client, requestToken.OAuthToken, authorizeTokenOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, authToken)

	// Create access token
	accessTokenOpts := oauth1.CreateAccessTokenOpts{
		OAuthConsumerKey:     consumer.ID,
		OAuthConsumerSecret:  consumer.Secret,
		OAuthToken:           requestToken.OAuthToken,
		OAuthTokenSecret:     requestToken.OAuthTokenSecret,
		OAuthVerifier:        authToken.OAuthVerifier,
		OAuthSignatureMethod: method,
	}
	accessToken, err := oauth1.CreateAccessToken(client, accessTokenOpts).Extract()
	th.AssertNoErr(t, err)
	defer oauth1.RevokeAccessToken(client, user.ID, accessToken.OAuthToken)
	tools.PrintResource(t, accessToken)

	// Get access token
	getAccessToken, err := oauth1.GetAccessToken(client, user.ID, accessToken.OAuthToken).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getAccessToken)
	th.AssertEquals(t, getAccessToken.ID, accessToken.OAuthToken)
	th.AssertEquals(t, getAccessToken.ConsumerID, consumer.ID)
	th.AssertEquals(t, getAccessToken.AuthorizingUserID, user.ID)
	th.AssertEquals(t, getAccessToken.ProjectID, project.ID)

	// List access tokens
	accessTokensPages, err := oauth1.ListAccessTokens(client, user.ID).AllPages()
	th.AssertNoErr(t, err)
	accessTokens, err := oauth1.ExtractAccessTokens(accessTokensPages)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, accessTokens)
	th.AssertDeepEquals(t, accessTokens[0], *getAccessToken)

	// List access token roles
	accessTokenRolesPages, err := oauth1.ListAccessTokenRoles(client, user.ID, accessToken.OAuthToken).AllPages()
	th.AssertNoErr(t, err)
	accessTokenRoles, err := oauth1.ExtractAccessTokenRoles(accessTokenRolesPages)
	th.AssertNoErr(t, err)
	tools.PrintResource(t, accessTokenRoles)
	th.AssertEquals(t, accessTokenRoles[0].ID, roles[0].ID)
	if len(accessTokenRoles) > 1 {
		th.AssertEquals(t, accessTokenRoles[1].Name, roles[1].Name)
	}

	// Get access token role
	getAccessTokenRole, err := oauth1.GetAccessTokenRole(client, user.ID, accessToken.OAuthToken, roles[0].ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, getAccessTokenRole)
	th.AssertDeepEquals(t, *getAccessTokenRole, accessTokenRoles[0])

	// Test auth using OAuth1
	newClient, err := clients.NewIdentityV3UnauthenticatedClient()
	th.AssertNoErr(t, err)

	// Opts to auth using an oauth1 credential
	authOptions := &oauth1.AuthOptions{
		OAuthConsumerKey:     consumer.ID,
		OAuthConsumerSecret:  consumer.Secret,
		OAuthToken:           accessToken.OAuthToken,
		OAuthTokenSecret:     accessToken.OAuthTokenSecret,
		OAuthSignatureMethod: method,
	}
	err = openstack.AuthenticateV3(newClient.ProviderClient, authOptions, gophercloud.EndpointOpts{})
	th.AssertNoErr(t, err)

	// Test OAuth1 token extract
	var token struct {
		tokens.Token
		oauth1.TokenExt
	}
	tokenRes := tokens.Get(newClient, newClient.ProviderClient.TokenID)
	err = tokenRes.ExtractInto(&token)
	th.AssertNoErr(t, err)
	oauth1Roles, err := tokenRes.ExtractRoles()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, token)
	tools.PrintResource(t, oauth1Roles)
	th.AssertEquals(t, token.OAuth1.ConsumerID, consumer.ID)
	th.AssertEquals(t, token.OAuth1.AccessTokenID, accessToken.OAuthToken)
}
