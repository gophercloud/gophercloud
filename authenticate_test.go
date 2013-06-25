package gophercloud

import (
	"testing"
)

func TestAuthProvider(t *testing.T) {
	c := TestContext()

	_, err := c.Authenticate("", AuthOptions{})
	if err == nil {
		t.Error("Expected error for empty provider string")
		return
	}
	_, err = c.Authenticate("unknown-provider", AuthOptions{Username: "u", Password: "p"})
	if err == nil {
		t.Error("Expected error for unknown service provider")
		return
	}

	err = c.RegisterProvider("provider", &Provider{})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = c.Authenticate("provider", AuthOptions{Username: "u", Password: "p"})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUserNameAndPassword(t *testing.T) {
	c := TestContext()
	c.RegisterProvider("provider", &Provider{})

	auths := []AuthOptions{
		AuthOptions{},
		AuthOptions{Username: "u"},
		AuthOptions{Password: "p"},
	}
	for i, auth := range auths {
		_, err := c.Authenticate("provider", auth)
		if err == nil {
			t.Error("Expected error from missing credentials (%d)", i)
			return
		}
	}

	_, err := c.Authenticate("provider", AuthOptions{Username: "u", Password: "p"})
	if err != nil {
		t.Error(err)
		return
	}
}
