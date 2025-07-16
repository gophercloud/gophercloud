package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/oauth1"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateConsumer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateConsumer(t, fakeServer)

	consumer, err := oauth1.CreateConsumer(context.TODO(), client.ServiceClient(fakeServer), oauth1.CreateConsumerOpts{
		Description: "My consumer",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, Consumer, *consumer)
}

func TestUpdateConsumer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateConsumer(t, fakeServer)

	consumer, err := oauth1.UpdateConsumer(context.TODO(), client.ServiceClient(fakeServer), "7fea2d", oauth1.UpdateConsumerOpts{
		Description: "My new consumer",
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, UpdatedConsumer, *consumer)
}

func TestDeleteConsumer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteConsumer(t, fakeServer)

	err := oauth1.DeleteConsumer(context.TODO(), client.ServiceClient(fakeServer), "7fea2d").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetConsumer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetConsumer(t, fakeServer)

	consumer, err := oauth1.GetConsumer(context.TODO(), client.ServiceClient(fakeServer), "7fea2d").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, FirstConsumer, *consumer)
}

func TestListConsumers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListConsumers(t, fakeServer)

	count := 0
	err := oauth1.ListConsumers(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListConsumers(t, fakeServer)

	allPages, err := oauth1.ListConsumers(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := oauth1.ExtractConsumers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedConsumersSlice, actual)
}

func TestRequestToken(t *testing.T) {
	fakeServer := th.SetupPersistentPortHTTP(t, 33199)
	defer fakeServer.Teardown()
	HandleRequestToken(t, fakeServer)

	ts := time.Unix(0, 0)
	token, err := oauth1.RequestToken(context.TODO(), client.ServiceClient(fakeServer), oauth1.RequestTokenOpts{
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAuthorizeToken(t, fakeServer)

	token, err := oauth1.AuthorizeToken(context.TODO(), client.ServiceClient(fakeServer), "29971f", oauth1.AuthorizeTokenOpts{
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
	fakeServer := th.SetupPersistentPortHTTP(t, 33199)
	defer fakeServer.Teardown()
	HandleCreateAccessToken(t, fakeServer)

	ts := time.Unix(1586804894, 0)
	token, err := oauth1.CreateAccessToken(context.TODO(), client.ServiceClient(fakeServer), oauth1.CreateAccessTokenOpts{
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetAccessToken(t, fakeServer)

	token, err := oauth1.GetAccessToken(context.TODO(), client.ServiceClient(fakeServer), "ce9e07", "6be26a").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, UserAccessToken, *token)
}

func TestRevokeAccessToken(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRevokeAccessToken(t, fakeServer)

	err := oauth1.RevokeAccessToken(context.TODO(), client.ServiceClient(fakeServer), "ce9e07", "6be26a").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListAccessTokens(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAccessTokens(t, fakeServer)

	count := 0
	err := oauth1.ListAccessTokens(client.ServiceClient(fakeServer), "ce9e07").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAccessTokens(t, fakeServer)

	allPages, err := oauth1.ListAccessTokens(client.ServiceClient(fakeServer), "ce9e07").AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := oauth1.ExtractAccessTokens(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedUserAccessTokensSlice, actual)
}

func TestListAccessTokenRoles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAccessTokenRoles(t, fakeServer)

	count := 0
	err := oauth1.ListAccessTokenRoles(client.ServiceClient(fakeServer), "ce9e07", "6be26a").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListAccessTokenRoles(t, fakeServer)

	allPages, err := oauth1.ListAccessTokenRoles(client.ServiceClient(fakeServer), "ce9e07", "6be26a").AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := oauth1.ExtractAccessTokenRoles(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedUserAccessTokenRolesSlice, actual)
}

func TestGetAccessTokenRole(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetAccessTokenRole(t, fakeServer)

	role, err := oauth1.GetAccessTokenRole(context.TODO(), client.ServiceClient(fakeServer), "ce9e07", "6be26a", "5ad150").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, UserAccessTokenRole, *role)
}

func TestAuthenticate(t *testing.T) {
	fakeServer := th.SetupPersistentPortHTTP(t, 33199)
	defer fakeServer.Teardown()
	HandleAuthenticate(t, fakeServer)

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

	actual, err := oauth1.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}
