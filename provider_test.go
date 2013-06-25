package gophercloud

import (
	"testing"
)

func TestProviderRegistry(t *testing.T) {
	c := TestContext()

	_, err := c.ProviderByName("aProvider")
	if err == nil {
		t.Error("Expected error when looking for a provider by non-existant name")
		return
	}

	_ = c.RegisterProvider("aProvider", &Provider{})
	_, err = c.ProviderByName("aProvider")
	if err != nil {
		t.Error(err)
		return
	}
}
