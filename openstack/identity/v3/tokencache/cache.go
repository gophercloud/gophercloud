// Package tokencache provides token caching for Identity v3 authentication.
package tokencache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
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

// Load retrieves a cached token ID. Missing, malformed, and expired entries are
// cache misses.
func Load(cache Cache, key, identityEndpoint string) (string, bool) {
	if cache == nil {
		return "", false
	}

	data, err := cache.Get(key)
	if err != nil || data == "" {
		return "", false
	}

	var cached cachedToken
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		_ = cache.Delete(key)
		return "", false
	}
	if cached.TokenID == "" || !cached.valid(identityEndpoint) {
		_ = cache.Delete(key)
		return "", false
	}
	return cached.TokenID, true
}

// Persist stores a token ID and expiration. Cache errors are ignored.
func Persist(cache Cache, key, identityEndpoint, tokenID string, expiresAt time.Time) {
	if cache == nil || tokenID == "" {
		return
	}

	data, err := json.Marshal(cachedToken{
		TokenID:   tokenID,
		ExpiresAt: expiresAt,
		Endpoint:  identityEndpoint,
	})
	if err != nil {
		return
	}
	_ = cache.Set(key, string(data))
}
