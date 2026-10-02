package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

// OAuth1SignatureMethod is an OAuth1 signature method.
type OAuth1SignatureMethod string

const (
	// OAuth1HMACSHA1 signs OAuth1 requests with HMAC-SHA1.
	OAuth1HMACSHA1 OAuth1SignatureMethod = "HMAC-SHA1"

	// OAuth1Plaintext signs OAuth1 requests with PLAINTEXT. It is not
	// recommended for production use.
	OAuth1Plaintext OAuth1SignatureMethod = "PLAINTEXT"
)

// V3OAuth1Opts contains options for Keystone v3 OAuth1 authentication.
type V3OAuth1Opts struct {
	// ConsumerKey is the OAuth1 consumer key.
	ConsumerKey string `q:"oauth_consumer_key" required:"true" json:"consumer_key,omitempty"`

	// ConsumerSecret is the OAuth1 consumer secret.
	ConsumerSecret string `json:"consumer_secret,omitempty"`

	// Token is the OAuth1 access token.
	Token string `q:"oauth_token" required:"true" json:"access_key,omitempty"`

	// TokenSecret is the OAuth1 access token secret.
	TokenSecret string `json:"access_secret,omitempty"`

	// SignatureMethod selects the OAuth1 request signature method.
	SignatureMethod OAuth1SignatureMethod `q:"oauth_signature_method" required:"true" json:"-"`

	// Timestamp overrides the OAuth1 request timestamp.
	Timestamp *time.Time `json:"-"`

	// Nonce overrides the OAuth1 request nonce.
	Nonce string `q:"oauth_nonce" json:"-"`

	// AllowReauth allows Gophercloud to reauthenticate automatically.
	AllowReauth bool `json:"allow_reauth,omitempty"`
}

// ToAuthBody builds the OAuth1 identity method body.
func (opts V3OAuth1Opts) ToAuthBody() (map[string]map[string]any, error) {
	for _, required := range []struct {
		argument string
		value    string
	}{
		{argument: "ConsumerKey", value: opts.ConsumerKey},
		{argument: "ConsumerSecret", value: opts.ConsumerSecret},
		{argument: "Token", value: opts.Token},
		{argument: "TokenSecret", value: opts.TokenSecret},
		{argument: "SignatureMethod", value: string(opts.SignatureMethod)},
	} {
		if required.value == "" {
			return nil, gophercloud.ErrMissingInput{Argument: required.argument}
		}
	}

	return map[string]map[string]any{AuthV3OAuth1.toAuthMethod(): {}}, nil
}

// ToAuthHeaders builds the OAuth1 Authorization header.
func (opts V3OAuth1Opts) ToAuthHeaders(options ...RequestOption) (map[string]any, error) {
	var ro requestOptions
	for _, option := range options {
		option(&ro)
	}

	if ro.tokenURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "tokenURL"}
	}

	queryURL, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	query := queryURL.Query()

	if opts.Timestamp == nil {
		query.Set("oauth_timestamp", strconv.FormatInt(time.Now().UTC().Unix(), 10))
	} else {
		query.Set("oauth_timestamp", strconv.FormatInt(opts.Timestamp.Unix(), 10))
	}
	if query.Get("oauth_nonce") == "" {
		nonce := make([]byte, 16)
		if _, err := rand.Read(nonce); err != nil {
			return nil, err
		}
		query.Set("oauth_nonce", hex.EncodeToString(nonce))
	}
	query.Set("oauth_version", "1.0")

	stringToSign := oauth1StringToSign(http.MethodPost, ro.tokenURL, query)
	signature := url.QueryEscape(oauth1Sign(opts.SignatureMethod, stringToSign, opts.ConsumerSecret, opts.TokenSecret))

	return map[string]any{
		"Authorization": oauth1AuthorizationHeader(query, signature),
	}, nil
}

// ToAuthScope returns no scope for OAuth1 authentication.
func (opts V3OAuth1Opts) ToAuthScope() (map[string]any, error) {
	return nil, nil
}

// ToAuthType returns the Keystone v3 OAuth1 authentication type.
func (opts V3OAuth1Opts) ToAuthType() AuthType {
	return AuthV3OAuth1
}

// CanReauth reports whether OAuth1 authentication can be repeated automatically.
func (opts V3OAuth1Opts) CanReauth() bool {
	return opts.AllowReauth
}

func oauth1StringToSign(method, rawURL string, query url.Values) []byte {
	parsedURL, _ := url.Parse(rawURL)
	port := parsedURL.Port()
	if parsedURL.Scheme == "http" && port == "80" || parsedURL.Scheme == "https" && port == "443" {
		parsedURL.Host = strings.TrimSuffix(parsedURL.Host, ":"+port)
	}
	parsedURL.RawQuery = ""

	return []byte(strings.Join([]string{
		method,
		url.QueryEscape(parsedURL.String()),
		url.QueryEscape(query.Encode()),
	}, "&"))
}

func oauth1Sign(method OAuth1SignatureMethod, value []byte, consumerSecret, tokenSecret string) string {
	key := []byte(url.QueryEscape(consumerSecret) + "&" + url.QueryEscape(tokenSecret))
	if method == OAuth1Plaintext {
		return string(key)
	}

	h := hmac.New(sha1.New, key)
	_, _ = h.Write(value)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func oauth1AuthorizationHeader(query url.Values, signature string) string {
	keys := make([]string, 0, len(query))
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys)+1)
	for _, key := range keys {
		for _, value := range query[key] {
			values = append(values, fmt.Sprintf("%s=%q", key, url.QueryEscape(value)))
		}
	}
	values = append(values, fmt.Sprintf("oauth_signature=%q", signature))
	return "OAuth " + strings.Join(values, ", ")
}
