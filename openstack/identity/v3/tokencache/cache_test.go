package tokencache_test

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokencache"
)

func TestKeySeparatesFields(t *testing.T) {
	first := tokencache.Key(tokencache.KeyOptions{
		Flow:      "flow\nprincipal=principal",
		Principal: "",
	})
	second := tokencache.Key(tokencache.KeyOptions{
		Flow:      "flow",
		Principal: "principal\nprincipal=",
	})

	if first == second {
		t.Fatalf("different cache identities produced the same key %q", first)
	}
}
