package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/oauth1"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateConsumer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateConsumer(t)

	consumer, err := oauth1.CreateConsumer(client.ServiceClient(), oauth1.CreateConsumerOpts{
		Description: "My consumer",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, Consumer, *consumer)
}

func TestUpdateConsumer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateConsumer(t)

	consumer, err := oauth1.UpdateConsumer(client.ServiceClient(), "7fea2d", oauth1.UpdateConsumerOpts{
		Description: "My new consumer",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, UpdatedConsumer, *consumer)
}

func TestDeleteConsumer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteConsumer(t)

	err := oauth1.DeleteConsumer(client.ServiceClient(), "7fea2d").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetConsumer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetConsumer(t)

	consumer, err := oauth1.GetConsumer(client.ServiceClient(), "7fea2d").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, FirstConsumer, *consumer)
}

func TestListConsumers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListConsumers(t)

	count := 0
	err := oauth1.ListConsumers(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := oauth1.ExtractConsumers(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedConsumersSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListConsumersAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListConsumers(t)

	allPages, err := oauth1.ListConsumers(client.ServiceClient()).AllPages()
	th.AssertNoErr(t, err)
	actual, err := oauth1.ExtractConsumers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedConsumersSlice, actual)
}

func TestRequestToken(t *testing.T) {
	th.SetupPersistentPortHTTP(t, 33199)
	defer th.TeardownHTTP()
	HandleRequestToken(t)

	ts := time.Unix(0, 0)
	token, err := oauth1.RequestToken(client.ServiceClient(), oauth1.RequestTokenOpts{
		OAuthConsumerKey:     Consumer.ID,
		OAuthConsumerSecret:  Consumer.Secret,
		OAuthSignatureMethod: oauth1.HMACSHA1,
		OAuthTimestamp:       &ts,
		OAuthNonce:           "71416001758914252991586795052",
		RequestedProjectID:   "1df927e8a466498f98788ed73d3c8ab4",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, Token, *token)
}

func TestAuthorizeToken(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAuthorizeToken(t)

	token, err := oauth1.AuthorizeToken(client.ServiceClient(), "29971f", oauth1.AuthorizeTokenOpts{
		Roles: []oauth1.Role{
			{
				ID: "a3b29b",
			},
			{
				ID: "49993e",
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "8171", token.OAuthVerifier)
}

func TestCreateAccessToken(t *testing.T) {
	th.SetupPersistentPortHTTP(t, 33199)
	defer th.TeardownHTTP()
	HandleCreateAccessToken(t)

	ts := time.Unix(1586804894, 0)
	token, err := oauth1.CreateAccessToken(client.ServiceClient(), oauth1.CreateAccessTokenOpts{
		OAuthConsumerKey:     Consumer.ID,
		OAuthConsumerSecret:  Consumer.Secret,
		OAuthToken:           Token.OAuthToken,
		OAuthTokenSecret:     Token.OAuthTokenSecret,
		OAuthVerifier:        "8171",
		OAuthSignatureMethod: oauth1.HMACSHA1,
		OAuthTimestamp:       &ts,
		OAuthNonce:           "66148873158553341551586804894",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, AccessToken, *token)
}

func TestGetAccessToken(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetAccessToken(t)

	token, err := oauth1.GetAccessToken(client.ServiceClient(), "ce9e07", "6be26a").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, UserAccessToken, *token)
}

func TestRevokeAccessToken(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleRevokeAccessToken(t)

	err := oauth1.RevokeAccessToken(client.ServiceClient(), "ce9e07", "6be26a").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListAccessTokens(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListAccessTokens(t)

	count := 0
	err := oauth1.ListAccessTokens(client.ServiceClient(), "ce9e07").EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := oauth1.ExtractAccessTokens(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedUserAccessTokensSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAccessTokensAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListAccessTokens(t)

	allPages, err := oauth1.ListAccessTokens(client.ServiceClient(), "ce9e07").AllPages()
	th.AssertNoErr(t, err)
	actual, err := oauth1.ExtractAccessTokens(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedUserAccessTokensSlice, actual)
}

func TestListAccessTokenRoles(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListAccessTokenRoles(t)

	count := 0
	err := oauth1.ListAccessTokenRoles(client.ServiceClient(), "ce9e07", "6be26a").EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := oauth1.ExtractAccessTokenRoles(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedUserAccessTokenRolesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListAccessTokenRolesAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListAccessTokenRoles(t)

	allPages, err := oauth1.ListAccessTokenRoles(client.ServiceClient(), "ce9e07", "6be26a").AllPages()
	th.AssertNoErr(t, err)
	actual, err := oauth1.ExtractAccessTokenRoles(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedUserAccessTokenRolesSlice, actual)
}

func TestGetAccessTokenRole(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetAccessTokenRole(t)

	role, err := oauth1.GetAccessTokenRole(client.ServiceClient(), "ce9e07", "6be26a", "5ad150").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, UserAccessTokenRole, *role)
}

func TestAuthenticate(t *testing.T) {
	th.SetupPersistentPortHTTP(t, 33199)
	defer th.TeardownHTTP()
	HandleAuthenticate(t)

	expected := &tokens.Token{
		ExpiresAt: time.Date(2017, 6, 3, 2, 19, 49, 0, time.UTC),
	}

	ts := time.Unix(0, 0)
	options := &oauth1.AuthOptions{
		OAuthConsumerKey:     Consumer.ID,
		OAuthConsumerSecret:  Consumer.Secret,
		OAuthToken:           AccessToken.OAuthToken,
		OAuthTokenSecret:     AccessToken.OAuthTokenSecret,
		OAuthSignatureMethod: oauth1.HMACSHA1,
		OAuthTimestamp:       &ts,
		OAuthNonce:           "66148873158553341551586804894",
	}

	actual, err := oauth1.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}
