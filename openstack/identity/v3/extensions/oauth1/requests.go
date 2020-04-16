package oauth1

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/pagination"
)

// Type SignatureMethod is a OAuth1 SignatureMethod type.
type SignatureMethod string

// Supported attributes for OAuth1 SignatureMethod.
const (
	// HMACSHA1 is a recommended OAuth1 signature method.
	HMACSHA1 SignatureMethod = "HMAC-SHA1"
	// PLAINTEXT signature method is not recommended to be used in
	// production environment.
	PLAINTEXT SignatureMethod = "PLAINTEXT"
)

// AuthOptions represents options for authenticating a user using OAuth1 tokens.
type AuthOptions struct {
	// OAuthConsumerKey is the OAuth1 Consumer Key.
	OAuthConsumerKey string `q:"oauth_consumer_key" required:"true"`
	// OAuthConsumerSecret is the OAuth1 Consumer Secret. Used to generate
	// an OAuth1 request signature.
	OAuthConsumerSecret string `required:"true"`
	// OAuthToken is the OAuth1 Request Token.
	OAuthToken string `q:"oauth_token" required:"true"`
	// OAuthTokenSecret is the OAuth1 Request Token Secret. Used to generate
	// an OAuth1 request signature.
	OAuthTokenSecret string `required:"true"`
	// OAuthSignatureMethod is the OAuth1 signature method the Consumer used
	// to sign the request. Supported values are "HMAC-SHA1" or "PLAINTEXT".
	// "PLAINTEXT" is not recommended for production usage.
	OAuthSignatureMethod SignatureMethod `q:"oauth_signature_method" required:"true"`
	// OAuthTimestamp is an OAuth1 request timestamp. If nil, current Unix
	// timestamp will be used.
	OAuthTimestamp *time.Time
	// OAuthNonce is an OAuth1 request nonce. Nonce must be a random string,
	// uniquely generated for each request. Will be generated automatically
	// when it is not set.
	OAuthNonce string `q:"oauth_nonce"`
	// AllowReauth allows Gophercloud to re-authenticate automatically
	// if/when your token expires.
	AllowReauth bool
}

// ToTokenV3CreateHeaders builds a create request headers from the AuthOptions.
func (opts AuthOptions) ToTokenV3CreateHeaders(method string, url string) (map[string]string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Authorization": genAuthHeader(
			opts.OAuthSignatureMethod,
			method,
			url,
			tokenParams(q.Query(), opts.OAuthTimestamp, ""),
			[]string{opts.OAuthConsumerSecret, opts.OAuthTokenSecret},
		),
		"X-Auth-Token": "",
	}, nil
}

// ToTokenV3CreateMap builds a create request body.
func (opts AuthOptions) ToTokenV3CreateMap(map[string]interface{}) (map[string]interface{}, error) {
	var req struct {
		Auth struct {
			Identity struct {
				Methods []string `json:"methods"`
				OAuth1  struct{} `json:"oauth1"`
			} `json:"identity"`
		} `json:"auth"`
	}
	req.Auth.Identity.Methods = []string{"oauth1"}
	return gophercloud.BuildRequestBody(req, "")
}

// ToTokenV3ScopeMap builds a scope from AuthOpts.
func (opts AuthOptions) ToTokenV3ScopeMap() (map[string]interface{}, error) {
	return nil, nil
}

func (opts AuthOptions) CanReauth() bool {
	return opts.AllowReauth
}

// Create authenticates and either generates a new OpenStack token from an
// OAuth1 token.
func Create(client *gophercloud.ServiceClient, opts tokens.AuthOptionsBuilder) (r tokens.CreateResult) {
	b, err := opts.ToTokenV3CreateMap(nil)
	if err != nil {
		r.Err = err
		return
	}

	method := "POST"
	url := authURL(client)
	h, err := opts.ToTokenV3CreateHeaders(method, url)
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(url, b, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	if resp != nil {
		r.Header = resp.Header
	}
	r.Err = err
	return
}

// CreateConsumerOptsBuilder allows extensions to add additional parameters to
// the CreateConsumer request.
type CreateConsumerOptsBuilder interface {
	ToCreateConsumerMap() (map[string]interface{}, error)
}

// CreateConsumerOpts provides options used to create a new Consumer.
type CreateConsumerOpts struct {
	// Description is the consumer description.
	Description string `json:"description"`
}

// ToCreateConsumerMap formats a CreateConsumerOpts into a create request.
func (opts CreateConsumerOpts) ToCreateConsumerMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "consumer")
}

// Create creates a new Consumer.
func CreateConsumer(client *gophercloud.ServiceClient, opts CreateConsumerOptsBuilder) (r CreateConsumerResult) {
	b, err := opts.ToCreateConsumerMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(consumersURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// Delete deletes a Consumer.
func DeleteConsumer(client *gophercloud.ServiceClient, id string) (r DeleteConsumerResult) {
	_, r.Err = client.Delete(consumerURL(client, id), nil)
	return
}

// List enumerates Consumers.
func ListConsumers(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, consumersURL(client), func(r pagination.PageResult) pagination.Page {
		return ConsumersPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetConsumer retrieves details on a single Consumer by ID.
func GetConsumer(client *gophercloud.ServiceClient, id string) (r GetConsumerResult) {
	_, r.Err = client.Get(consumerURL(client, id), &r.Body, nil)
	return
}

// UpdateConsumerOpts provides options used to update a consumer.
type UpdateConsumerOpts struct {
	// Description is the consumer description.
	Description string `json:"description"`
}

// ToUpdateConsumerMap formats an UpdateConsumerOpts into a consumer update
// request.
func (opts UpdateConsumerOpts) ToUpdateConsumerMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "consumer")
}

// UpdateConsumer updates an existing Consumer.
func UpdateConsumer(client *gophercloud.ServiceClient, id string, opts UpdateConsumerOpts) (r UpdateConsumerResult) {
	b, err := opts.ToUpdateConsumerMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(consumerURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// RequestTokenOptsBuilder allows extensions to add additional parameters to the
// RequestToken request.
type RequestTokenOptsBuilder interface {
	ToRequestTokenHeaders(string, string) (map[string]string, error)
}

// RequestTokenOpts provides options used to get a consumer unauthorized
// request token.
type RequestTokenOpts struct {
	// OAuthConsumerKey is the OAuth1 Consumer Key.
	OAuthConsumerKey string `q:"oauth_consumer_key" required:"true"`
	// OAuthConsumerSecret is the OAuth1 Consumer Secret. Used to generate
	// an OAuth1 request signature.
	OAuthConsumerSecret string `required:"true"`
	// OAuthSignatureMethod is the OAuth1 signature method the Consumer used
	// to sign the request. Supported values are "HMAC-SHA1" or "PLAINTEXT".
	// "PLAINTEXT" is not recommended for production usage.
	OAuthSignatureMethod SignatureMethod `q:"oauth_signature_method" required:"true"`
	// OAuthTimestamp is an OAuth1 request timestamp. If nil, current Unix
	// timestamp will be used.
	OAuthTimestamp *time.Time
	// OAuthNonce is an OAuth1 request nonce. Nonce must be a random string,
	// uniquely generated for each request. Will be generated automatically
	// when it is not set.
	OAuthNonce string `q:"oauth_nonce"`
	// RequestedProjectID is a Project ID a consumer user requested an
	// access to.
	RequestedProjectID string `h:"Requested-Project-Id"`
}

// ToRequestTokenHeaders formats a RequestTokenOpts into a map of request
// headers.
func (opts RequestTokenOpts) ToRequestTokenHeaders(method, url string) (map[string]string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil, err
	}

	h["Authorization"] = genAuthHeader(
		opts.OAuthSignatureMethod,
		method,
		url,
		tokenParams(q.Query(), opts.OAuthTimestamp, "oob"),
		[]string{opts.OAuthConsumerSecret},
	)

	return h, nil
}

// RequestToken requests an unauthorized OAuth1 Token.
func RequestToken(client *gophercloud.ServiceClient, opts RequestTokenOptsBuilder) (r TokenResult) {
	method := "POST"
	url := requestTokenURL(client)
	h, err := opts.ToRequestTokenHeaders(method, url)
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Request(method, url, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	if resp != nil {
		r.Header = resp.Header
		r.Body = resp.Body
	}
	r.Err = err
	return
}

// AuthorizeTokenOptsBuilder allows extensions to add additional parameters to
// the AuthorizeToken request.
type AuthorizeTokenOptsBuilder interface {
	ToAuthorizeTokenMap() (map[string]interface{}, error)
}

// AuthorizeTokenOpts provides options used to authorize a request token.
type AuthorizeTokenOpts struct {
	Roles []Role `json:"roles"`
}

// Role is a struct representing a role object in a AuthorizeTokenOpts struct.
type Role struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ToAuthorizeTokenMap formats an AuthorizeTokenOpts into an authorize token
// request.
func (opts AuthorizeTokenOpts) ToAuthorizeTokenMap() (map[string]interface{}, error) {
	for _, r := range opts.Roles {
		if r == (Role{}) {
			return nil, fmt.Errorf("role must not be empty")
		}
	}
	return gophercloud.BuildRequestBody(opts, "")
}

// AuthorizeToken authorizes an unauthorized consumer token.
func AuthorizeToken(client *gophercloud.ServiceClient, id string, opts AuthorizeTokenOptsBuilder) (r AuthorizeTokenResult) {
	b, err := opts.ToAuthorizeTokenMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(authorizeTokenURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// CreateAccessTokenOptsBuilder allows extensions to add additional parameters
// to the CreateAccessToken request.
type CreateAccessTokenOptsBuilder interface {
	ToCreateAccessTokenHeaders(string, string) (map[string]string, error)
}

// CreateAccessTokenOpts provides options used to create an OAuth1 token.
type CreateAccessTokenOpts struct {
	// OAuthConsumerKey is the OAuth1 Consumer Key.
	OAuthConsumerKey string `q:"oauth_consumer_key" required:"true"`
	// OAuthConsumerSecret is the OAuth1 Consumer Secret. Used to generate
	// an OAuth1 request signature.
	OAuthConsumerSecret string `required:"true"`
	// OAuthToken is the OAuth1 Request Token.
	OAuthToken string `q:"oauth_token" required:"true"`
	// OAuthTokenSecret is the OAuth1 Request Token Secret. Used to generate
	// an OAuth1 request signature.
	OAuthTokenSecret string `required:"true"`
	// OAuthVerifier is the OAuth1 verification code.
	OAuthVerifier string `q:"oauth_verifier" required:"true"`
	// OAuthSignatureMethod is the OAuth1 signature method the Consumer used
	// to sign the request. Supported values are "HMAC-SHA1" or "PLAINTEXT".
	// "PLAINTEXT" is not recommended for production usage.
	OAuthSignatureMethod SignatureMethod `q:"oauth_signature_method" required:"true"`
	// OAuthTimestamp is an OAuth1 request timestamp. If nil, current Unix
	// timestamp will be used.
	OAuthTimestamp *time.Time
	// OAuthNonce is an OAuth1 request nonce. Nonce must be a random string,
	// uniquely generated for each request. Will be generated automatically
	// when it is not set.
	OAuthNonce string `q:"oauth_nonce"`
}

// ToCreateAccessTokenHeaders formats a CreateAccessTokenOpts into a map of
// request headers.
func (opts CreateAccessTokenOpts) ToCreateAccessTokenHeaders(method, url string) (map[string]string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Authorization": genAuthHeader(
			opts.OAuthSignatureMethod,
			method,
			url,
			tokenParams(q.Query(), opts.OAuthTimestamp, ""),
			[]string{opts.OAuthConsumerSecret, opts.OAuthTokenSecret},
		),
	}, nil
}

// CreateAccessToken creates a new OAuth1 Access Token
func CreateAccessToken(client *gophercloud.ServiceClient, opts CreateAccessTokenOptsBuilder) (r TokenResult) {
	method := "POST"
	url := createAccessTokenURL(client)
	h, err := opts.ToCreateAccessTokenHeaders(method, url)
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Request(method, url, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{201},
	})
	if resp != nil {
		r.Header = resp.Header
		r.Body = resp.Body
	}
	r.Err = err
	return
}

// GetAccessToken retrieves details on a single OAuth1 access token by an ID.
func GetAccessToken(client *gophercloud.ServiceClient, userID string, id string) (r GetAccessTokenResult) {
	_, r.Err = client.Get(userAccessTokenURL(client, userID, id), &r.Body, nil)
	return
}

// RevokeAccessToken revokes an OAuth1 access token.
func RevokeAccessToken(client *gophercloud.ServiceClient, userID string, id string) (r RevokeAccessTokenResult) {
	_, r.Err = client.Delete(userAccessTokenURL(client, userID, id), nil)
	return
}

// ListAccessTokens enumerates authorized access tokens.
func ListAccessTokens(client *gophercloud.ServiceClient, userID string) pagination.Pager {
	url := userAccessTokensURL(client, userID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AccessTokensPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAccessTokenRoles enumerates authorized access token roles.
func ListAccessTokenRoles(client *gophercloud.ServiceClient, userID string, id string) pagination.Pager {
	url := userAccessTokenRolesURL(client, userID, id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AccessTokenRolesPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// GetAccessTokenRole retrieves details on a single OAuth1 access token role by
// an ID.
func GetAccessTokenRole(client *gophercloud.ServiceClient, userID string, id string, roleID string) (r GetAccessTokenRoleResult) {
	_, r.Err = client.Get(userAccessTokenRoleURL(client, userID, id, roleID), &r.Body, nil)
	return
}

// The following are small helper functions used to help build the signature.

// tokenParams builds a URLEncoded parameters string.
func tokenParams(query url.Values, timestamp *time.Time, callback string) string {
	if timestamp != nil {
		// use provided timestamp
		query.Set("oauth_timestamp", strconv.FormatInt(timestamp.Unix(), 10))
	} else {
		// use current timestamp
		query.Set("oauth_timestamp", strconv.FormatInt(time.Now().UTC().Unix(), 10))
	}

	if query.Get("oauth_nonce") == "" {
		// when nonce is not set, generate a random one
		query.Set("oauth_nonce", strconv.FormatInt(rand.Int63(), 10)+query.Get("oauth_timestamp"))
	}

	if callback != "" {
		query.Set("oauth_callback", callback)
	}
	query.Set("oauth_version", "1.0")

	// sorted by key
	return query.Encode()
}

// stringToSign builds a string to be signed.
func stringToSign(method string, u string, params string) []byte {
	parsedURL, _ := url.Parse(u)
	p := parsedURL.Port()
	s := parsedURL.Scheme
	// Default scheme port must be stripped
	if s == "http" && p == "80" || s == "https" && p == "443" {
		parsedURL.Host = strings.TrimSuffix(parsedURL.Host, ":"+p)
	}
	// Ensure that URL doesn't contain queries
	parsedURL.RawQuery = ""
	return []byte(strings.Join([]string{
		method,
		url.QueryEscape(parsedURL.String()),
		url.QueryEscape(params),
	}, "&"))
}

// signString signs a string using an OAuth1 signature method.
func signString(signatureMethod SignatureMethod, strToSign []byte, secrets []string) string {
	var key []byte
	for i, k := range secrets {
		key = append(key, []byte(url.QueryEscape(k))...)
		if i == 0 {
			key = append(key, '&')
		}
	}
	if signatureMethod == PLAINTEXT {
		return string(key)
	}
	h := hmac.New(sha1.New, key)
	h.Write(strToSign)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// genAuthHeader generates an OAuth1 Authorization header with a signature
// calculated using an OAuth1 signature method.
func genAuthHeader(signatureMethod SignatureMethod, method, u, b string, secrets []string) string {
	var authHeader []string
	params, _ := url.ParseQuery(b)

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		for _, v := range params[k] {
			authHeader = append(authHeader, fmt.Sprintf("%s=%q", k, url.QueryEscape(v)))
		}
	}

	strToSign := stringToSign(method, u, b)
	signature := url.QueryEscape(signString(signatureMethod, strToSign, secrets))
	authHeader = append(authHeader, fmt.Sprintf("oauth_signature=%q", signature))

	return "OAuth " + strings.Join(authHeader, ", ")
}
