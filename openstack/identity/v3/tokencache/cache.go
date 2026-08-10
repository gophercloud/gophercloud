// Package tokencache provides token caching for Identity v3 authentication.
package tokencache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

const tokenExpiryMargin = 5 * time.Minute

// Cache stores authentication results. Implementations must be safe for concurrent use.
type Cache interface {
	// Get returns ("", nil) for a cache miss.
	Get(key string) (string, error)
	// Set stores a value for key.
	Set(key, value string) error
	// Delete removes key.
	Delete(key string) error
}

// KeyOptions identifies a cached identity.
type KeyOptions struct {
	// Flow identifies the authentication flow.
	Flow string
	// Principal identifies the client or local profile.
	Principal string
	// IdentityEndpoint is the Keystone endpoint.
	IdentityEndpoint string
	// IdentityProvider is the Keystone federation provider.
	IdentityProvider string
	// Protocol is the Keystone federation protocol.
	Protocol string
	// AuthenticationEndpoint is the external identity endpoint.
	AuthenticationEndpoint string
	// Scope is the requested Keystone scope.
	Scope tokens.Scope
}

// Key returns a deterministic cache identifier.
func Key(opts KeyOptions) string {
	opts.IdentityEndpoint = strings.TrimRight(opts.IdentityEndpoint, "/")
	data, _ := json.Marshal(opts)
	digest := sha256.Sum256(data)
	return "gophercloud-auth-" + hex.EncodeToString(digest[:])
}

type cachedToken struct {
	TokenID   string    `json:"token_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Endpoint  string    `json:"endpoint"`
}

func (ct *cachedToken) valid(endpoint string) bool {
	if strings.TrimRight(ct.Endpoint, "/") != strings.TrimRight(endpoint, "/") {
		return false
	}
	return time.Now().Add(tokenExpiryMargin).Before(ct.ExpiresAt)
}

// Load retrieves and validates a cached token against Keystone.
func Load(ctx context.Context, client *gophercloud.ProviderClient, cache Cache, key, identityEndpoint string) (tokens.CreateResult, bool) {
	var zero tokens.CreateResult
	if client == nil || cache == nil {
		return zero, false
	}

	data, err := cache.Get(key)
	if err != nil || data == "" {
		return zero, false
	}

	var cached cachedToken
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		_ = cache.Delete(key)
		return zero, false
	}
	if !cached.valid(identityEndpoint) {
		_ = cache.Delete(key)
		return zero, false
	}

	endpoint := gophercloud.NormalizeURL(identityEndpoint)
	if !strings.HasSuffix(strings.TrimRight(endpoint, "/"), "/v3") {
		endpoint = strings.TrimRight(endpoint, "/") + "/v3/"
	}
	identityClient := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{
			HTTPClient:        client.HTTPClient,
			UserAgent:         client.UserAgent,
			RetryBackoffFunc:  client.RetryBackoffFunc,
			MaxBackoffRetries: client.MaxBackoffRetries,
			RetryFunc:         client.RetryFunc,
		},
		Endpoint: endpoint,
		Type:     "identity",
	}

	identityClient.SetToken(cached.TokenID)
	result := tokens.Get(ctx, identityClient, cached.TokenID, nil)
	if result.Err != nil {
		_ = cache.Delete(key)
		return zero, false
	}

	var created tokens.CreateResult
	created.Body = result.Body
	created.Header = result.Header
	created.Err = result.Err
	return created, true
}

// Persist stores a successful token result. Cache errors are ignored.
func Persist(cache Cache, key, identityEndpoint string, result tokens.CreateResult) {
	if cache == nil {
		return
	}
	tokenID, err := result.ExtractTokenID()
	if err != nil || tokenID == "" {
		return
	}
	token, err := result.ExtractToken()
	if err != nil || token == nil {
		return
	}

	data, err := json.Marshal(cachedToken{
		TokenID:   tokenID,
		ExpiresAt: token.ExpiresAt,
		Endpoint:  identityEndpoint,
	})
	if err != nil {
		return
	}
	_ = cache.Set(key, string(data))
}
