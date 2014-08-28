package tokens

import (
	"testing"

	"github.com/rackspace/gophercloud"
)

func TestTokenURL(t *testing.T) {
	setup()
	defer teardown()

	client := gophercloud.ServiceClient{Endpoint: endpoint()}

	expected := endpoint() + "auth/tokens"
	actual := getTokenURL(&client)
	if actual != expected {
		t.Errorf("Expected URL %s, but was %s", expected, actual)
	}
}
