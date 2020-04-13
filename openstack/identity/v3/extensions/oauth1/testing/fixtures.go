package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/oauth1"
	tokens "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens/testing"
	"github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const CreateConsumerRequest = `
{
    "consumer": {
        "description": "My consumer"
    }
}
`

const CreateConsumerResponse = `
{
    "consumer": {
        "secret": "secretsecret",
        "description": "My consumer",
        "id": "7fea2d",
        "links": {
            "self": "http://example.com/identity/v3/OS-OAUTH1/consumers/7fea2d"
        }
    }
}
`

const UpdateConsumerRequest = `
{
    "consumer": {
        "description": "My new consumer"
    }
}
`

const UpdateConsumerResponse = `
{
    "consumer": {
        "description": "My new consumer",
        "id": "7fea2d",
        "links": {
            "self": "http://example.com/identity/v3/OS-OAUTH1/consumers/7fea2d"
        }
    }
}
`

// GetConsumerOutput provides a Get result.
const GetConsumerResponse = `
{
    "consumer": {
        "id": "7fea2d",
        "description": "My consumer",
        "links": {
            "self": "http://example.com/identity/v3/OS-OAUTH1/consumers/7fea2d"
        }
    }
}
`

// ListConsumersResponse provides a single page of Consumers results.
const ListConsumersResponse = `
{
    "consumers": [
        {
            "description": "My consumer",
            "id": "7fea2d",
            "links": {
                "self": "http://example.com/identity/v3/OS-OAUTH1/consumers/7fea2d"
            }
        },
        {
            "id": "0c2a74",
            "links": {
                "self": "http://example.com/identity/v3/OS-OAUTH1/consumers/0c2a74"
            }
        }
    ],
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/OS-OAUTH1/consumers"
    }
}
`

const AuthorizeTokenRequest = `
{
    "roles": [
        {
            "id": "a3b29b"
        },
        {
            "id": "49993e"
        }
    ]
}
`

const AuthorizeTokenResponse = `
{
    "token": {
        "oauth_verifier": "8171"
    }
}
`

const GetUserAccessTokenResponse = `
{
    "access_token": {
        "consumer_id": "7fea2d",
        "id": "6be26a",
        "expires_at": "2013-09-11T06:07:51.501805Z",
        "links": {
            "roles": "http://example.com/identity/v3/users/ce9e07/OS-OAUTH1/access_tokens/6be26a/roles",
            "self": "http://example.com/identity/v3/users/ce9e07/OS-OAUTH1/access_tokens/6be26a"
        },
        "project_id": "b9fca3",
        "authorizing_user_id": "ce9e07"
    }
}
`

const ListUserAccessTokensResponse = `
{
    "access_tokens": [
        {
            "consumer_id": "7fea2d",
            "id": "6be26a",
            "expires_at": "2013-09-11T06:07:51.501805Z",
            "links": {
                "roles": "http://example.com/identity/v3/users/ce9e07/OS-OAUTH1/access_tokens/6be26a/roles",
                "self": "http://example.com/identity/v3/users/ce9e07/OS-OAUTH1/access_tokens/6be26a"
            },
            "project_id": "b9fca3",
            "authorizing_user_id": "ce9e07"
        }
    ],
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/users/ce9e07/OS-OAUTH1/access_tokens"
    }
}
`

const ListUserAccessTokenRolesResponse = `
{
    "roles": [
        {
            "id": "5ad150",
            "domain_id": "7cf37b",
            "links": {
                "self": "http://example.com/identity/v3/roles/5ad150"
            },
            "name": "admin"
        },
        {
            "id": "a62eb6",
            "domain_id": "7cf37b",
            "links": {
                "self": "http://example.com/identity/v3/roles/a62eb6"
            },
            "name": "member"
        }
    ],
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/users/ce9e07/OS-OAUTH1/access_tokens/6be26a/roles"
    }
}
`

const ListUserAccessTokenRoleResponse = `
{
    "role": {
        "id": "5ad150",
        "domain_id": "7cf37b",
        "links": {
            "self": "http://example.com/identity/v3/roles/5ad150"
        },
        "name": "admin"
    }
}
`

var tokenExpiresAt = time.Date(2013, time.September, 11, 06, 07, 51, 501805000, time.UTC)
var UserAccessToken = oauth1.AccessToken{
	ID:                "6be26a",
	ConsumerID:        "7fea2d",
	ProjectID:         "b9fca3",
	AuthorizingUserID: "ce9e07",
	ExpiresAt:         &tokenExpiresAt,
}

var UserAccessTokenRole = oauth1.AccessTokenRole{
	ID:       "5ad150",
	DomainID: "7cf37b",
	Name:     "admin",
}

var UserAccessTokenRoleSecond = oauth1.AccessTokenRole{
	ID:       "a62eb6",
	DomainID: "7cf37b",
	Name:     "member",
}

var ExpectedUserAccessTokensSlice = []oauth1.AccessToken{UserAccessToken}

var ExpectedUserAccessTokenRolesSlice = []oauth1.AccessTokenRole{UserAccessTokenRole, UserAccessTokenRoleSecond}

// HandleCreateConsumer creates an HTTP handler at `/OS-OAUTH1/consumers` on the
// test handler mux that tests consumer creation.
func HandleCreateConsumer(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-OAUTH1/consumers", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, CreateConsumerRequest)

		w.WriteHeader(http.StatusCreated)
		_, err := fmt.Fprintf(w, CreateConsumerResponse)
		testhelper.AssertNoErr(t, err)
	})
}

// HandleUpdateConsumer creates an HTTP handler at `/OS-OAUTH1/consumers/7fea2d` on the
// test handler mux that tests consumer update.
func HandleUpdateConsumer(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-OAUTH1/consumers/7fea2d", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "PATCH")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, UpdateConsumerRequest)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, UpdateConsumerResponse)
		testhelper.AssertNoErr(t, err)
	})
}

// HandleDeleteConsumer creates an HTTP handler at `/OS-OAUTH1/consumers/7fea2d` on the
// test handler mux that tests consumer deletion.
func HandleDeleteConsumer(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-OAUTH1/consumers/7fea2d", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "DELETE")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleGetConsumer creates an HTTP handler at `/OS-OAUTH1/consumers/7fea2d` on the
// test handler mux that responds with a single consumer.
func HandleGetConsumer(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-OAUTH1/consumers/7fea2d", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetConsumerResponse)
	})
}

var Consumer = oauth1.Consumer{
	ID:          "7fea2d",
	Description: "My consumer",
	Secret:      "secretsecret",
}

var UpdatedConsumer = oauth1.Consumer{
	ID:          "7fea2d",
	Description: "My new consumer",
}

var FirstConsumer = oauth1.Consumer{
	ID:          "7fea2d",
	Description: "My consumer",
}

var SecondConsumer = oauth1.Consumer{
	ID: "0c2a74",
}

// ExpectedConsumersSlice is the slice of consumers expected to be returned from ListOutput.
var ExpectedConsumersSlice = []oauth1.Consumer{FirstConsumer, SecondConsumer}

// HandleListConsumers creates an HTTP handler at `/OS-OAUTH1/consumers` on the
// test handler mux that responds with a list of two consumers.
func HandleListConsumers(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-OAUTH1/consumers", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListConsumersResponse)
	})
}

var Token = oauth1.Token{
	OAuthToken:       "29971f",
	OAuthTokenSecret: "238eb8",
	OAUthExpiresAt:   &tokenExpiresAt,
}

// HandleRequestToken creates an HTTP handler at `/OS-OAUTH1/request_token` on the
// test handler mux that responds with a OAuth1 unauthorized token.
func HandleRequestToken(t *testing.T) {
	testhelper.Mux.HandleFunc("/unit_test/OS-OAUTH1/request_token", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		testhelper.TestHeader(t, r, "Authorization", `OAuth oauth_callback="oob", oauth_consumer_key="7fea2d", oauth_nonce="71416001758914252991586795052", oauth_signature_method="HMAC-SHA1", oauth_timestamp="0", oauth_version="1.0", oauth_signature="PGbUCw8qdK0bAi4VmHqWqGhbma8%3D"`)
		testhelper.TestHeader(t, r, "Requested-Project-Id", "1df927e8a466498f98788ed73d3c8ab4")
		testhelper.TestBody(t, r, "")

		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `oauth_token=29971f&oauth_token_secret=238eb8&oauth_expires_at=2013-09-11T06:07:51.501805Z`)
	})
}

// HandleAuthorizeToken creates an HTTP handler at `/OS-OAUTH1/authorize/29971f` on the
// test handler mux that tests unauthorized token authorization.
func HandleAuthorizeToken(t *testing.T) {
	testhelper.Mux.HandleFunc("/OS-OAUTH1/authorize/29971f", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "PUT")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, AuthorizeTokenRequest)

		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, AuthorizeTokenResponse)
		testhelper.AssertNoErr(t, err)
	})
}

var AccessToken = oauth1.Token{
	OAuthToken:       "accd36",
	OAuthTokenSecret: "aa47da",
	OAUthExpiresAt:   &tokenExpiresAt,
}

// HandleCreateAccessToken creates an HTTP handler at `/OS-OAUTH1/access_token` on the
// test handler mux that responds with a OAuth1 access token.
func HandleCreateAccessToken(t *testing.T) {
	testhelper.Mux.HandleFunc("/unit_test/OS-OAUTH1/access_token", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		testhelper.TestHeader(t, r, "Authorization", `OAuth oauth_consumer_key="7fea2d", oauth_nonce="66148873158553341551586804894", oauth_signature_method="HMAC-SHA1", oauth_timestamp="1586804894", oauth_token="29971f", oauth_verifier="8171", oauth_version="1.0", oauth_signature="rvmK2yOhNL9NjjZMlwTn3dPK1u0%3D"`)
		testhelper.TestBody(t, r, "")

		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `oauth_token=accd36&oauth_token_secret=aa47da&oauth_expires_at=2013-09-11T06:07:51.501805Z`)
	})
}

// HandleGetAccessToken creates an HTTP handler at `/users/ce9e07/OS-OAUTH1/access_tokens/6be26a` on the
// test handler mux that responds with a single access token.
func HandleGetAccessToken(t *testing.T) {
	testhelper.Mux.HandleFunc("/users/ce9e07/OS-OAUTH1/access_tokens/6be26a", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetUserAccessTokenResponse)
	})
}

// HandleRevokeAccessToken creates an HTTP handler at `/users/ce9e07/OS-OAUTH1/access_tokens/6be26a` on the
// test handler mux that tests access token deletion.
func HandleRevokeAccessToken(t *testing.T) {
	testhelper.Mux.HandleFunc("/users/ce9e07/OS-OAUTH1/access_tokens/6be26a", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "DELETE")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleListAccessTokens creates an HTTP handler at `/users/ce9e07/OS-OAUTH1/access_tokens` on the
// test handler mux that responds with a slice of access tokens.
func HandleListAccessTokens(t *testing.T) {
	testhelper.Mux.HandleFunc("/users/ce9e07/OS-OAUTH1/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListUserAccessTokensResponse)
	})
}

// HandleListAccessTokenRoles creates an HTTP handler at `/users/ce9e07/OS-OAUTH1/access_tokens/6be26a/roles` on the
// test handler mux that responds with a slice of access token roles.
func HandleListAccessTokenRoles(t *testing.T) {
	testhelper.Mux.HandleFunc("/users/ce9e07/OS-OAUTH1/access_tokens/6be26a/roles", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListUserAccessTokenRolesResponse)
	})
}

// HandleGetAccessTokenRole creates an HTTP handler at `/users/ce9e07/OS-OAUTH1/access_tokens/6be26a/roles/5ad150` on the
// test handler mux that responds with an access token role.
func HandleGetAccessTokenRole(t *testing.T) {
	testhelper.Mux.HandleFunc("/users/ce9e07/OS-OAUTH1/access_tokens/6be26a/roles/5ad150", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListUserAccessTokenRoleResponse)
	})
}

// HandleAuthenticate creates an HTTP handler at `/auth/tokens` on the
// test handler mux that responds with an OpenStack token.
func HandleAuthenticate(t *testing.T) {
	testhelper.Mux.HandleFunc("/unit_test/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestHeader(t, r, "Authorization", `OAuth oauth_consumer_key="7fea2d", oauth_nonce="66148873158553341551586804894", oauth_signature_method="HMAC-SHA1", oauth_timestamp="0", oauth_token="accd36", oauth_version="1.0", oauth_signature="vxka5cHQee01I3AiDZFE3S53ZA4%3D"`)
		testhelper.TestJSONRequest(t, r, `{"auth": {"identity": {"oauth1": {}, "methods": ["oauth1"]}}}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, tokens.TokenOutput)
	})
}
