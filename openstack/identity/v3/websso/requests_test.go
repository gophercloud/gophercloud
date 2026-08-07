package websso_test

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/websso"
)

func TestAuthOptions(t *testing.T) {
	opts := websso.AuthOptions{
		IdentityEndpoint: "https://keystone.example.org:5000/v3",
		IdentityProvider: "my-idp",
		Protocol:         "openid",
		RedirectHost:     "localhost",
		RedirectPort:     9990,
		AllowReauth:      true,
	}

	// Test ToTokenV3ScopeMap with no scope
	scope, err := opts.ToTokenV3ScopeMap()
	if err != nil {
		t.Fatalf("ToTokenV3ScopeMap failed: %v", err)
	}
	if scope != nil {
		t.Errorf("Expected nil scope, got %v", scope)
	}

	// Test ToTokenV3ScopeMap with project scope
	opts.Scope = gophercloud.AuthScope{
		ProjectName: "my-project",
		DomainName:  "Default",
	}
	scope, err = opts.ToTokenV3ScopeMap()
	if err != nil {
		t.Fatalf("ToTokenV3ScopeMap with project failed: %v", err)
	}
	if scope == nil {
		t.Fatal("Expected non-nil scope")
	}

	// Test CanReauth
	if !opts.CanReauth() {
		t.Error("Expected CanReauth to return true")
	}
}

func TestCacheID(t *testing.T) {
	opts := websso.AuthOptions{
		IdentityEndpoint: "https://keystone.example.org:5000/v3",
		IdentityProvider: "my-idp",
	}

	cacheID := opts.GetCacheID()
	if cacheID == "" {
		t.Error("Expected non-empty cache ID")
	}
	if cacheID != "os-https-keystone-example-org-5000-v3-my-idp" {
		t.Errorf("Unexpected cache ID: %s", cacheID)
	}
}
