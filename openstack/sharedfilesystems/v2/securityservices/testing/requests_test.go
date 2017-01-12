package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// Verifies that a security service can be created correctly
func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &securityservices.CreateOpts{
		Name:        "SecServ1",
		Description: "Creating my first Security Service",
		DNSIP:       "10.0.0.0/24",
		User:        "demo",
		Password:    "***",
		Type:        "kerberos",
	}

	s, err := securityservices.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "SecServ1")
	th.AssertEquals(t, s.Description, "Creating my first Security Service")
	th.AssertEquals(t, s.User, "demo")
	th.AssertEquals(t, s.DNSIP, "10.0.0.0/24")
	th.AssertEquals(t, s.Password, "supersecret")
	th.AssertEquals(t, s.Type, "kerberos")
}

// Verifies that a security service cannot be created without a type
func TestCreateFails(t *testing.T) {
	options := &securityservices.CreateOpts{
		Name:        "SecServ1",
		Description: "Creating my first Security Service",
		DNSIP:       "10.0.0.0/24",
		User:        "demo",
		Password:    "***",
	}

	_, err := securityservices.Create(client.ServiceClient(), options).Extract()
	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

// Verifies that security service deletion works
func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := securityservices.Delete(client.ServiceClient(), "securityServiceID")
	th.AssertNoErr(t, res.Err)
}
