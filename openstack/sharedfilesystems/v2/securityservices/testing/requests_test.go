package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/securityservices"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Verifies that a security service can be created correctly
func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockCreateResponse(t, fakeServer)

	options := &securityservices.CreateOpts{
		Name:        "SecServ1",
		Description: "Creating my first Security Service",
		DNSIP:       "10.0.0.0/24",
		User:        "demo",
		Password:    "***",
		Type:        "kerberos",
	}

	s, err := securityservices.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "SecServ1")
	th.AssertEquals(t, s.Description, "Creating my first Security Service")
	th.AssertEquals(t, s.User, "demo")
	th.AssertEquals(t, s.DNSIP, "10.0.0.0/24")
	th.AssertEquals(t, s.Password, "supersecret")
	th.AssertEquals(t, s.Type, "kerberos")
}

// Verifies that a security service cannot be created without a type
func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	options := &securityservices.CreateOpts{
		Name:        "SecServ1",
		Description: "Creating my first Security Service",
		DNSIP:       "10.0.0.0/24",
		User:        "demo",
		Password:    "***",
	}

	_, err := securityservices.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()
	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

// Verifies that security service deletion works
func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockDeleteResponse(t, fakeServer)

	res := securityservices.Delete(context.TODO(), client.ServiceClient(fakeServer), "securityServiceID")
	th.AssertNoErr(t, res.Err)
}

// Verifies that security services can be listed correctly
func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListResponse(t, fakeServer)

	allPages, err := securityservices.List(client.ServiceClient(fakeServer), &securityservices.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := securityservices.ExtractSecurityServices(allPages)
	th.AssertNoErr(t, err)
	var nilTime time.Time
	expected := []securityservices.SecurityService{
		{
			Status:      "new",
			Domain:      "",
			ProjectID:   "16e1ab15c35a457e9c2b2aa189f544e1",
			Name:        "SecServ1",
			CreatedAt:   time.Date(2015, 9, 7, 12, 19, 10, 0, time.UTC),
			Description: "Creating my first Security Service",
			UpdatedAt:   nilTime,
			Server:      "",
			DNSIP:       "10.0.0.0/24",
			User:        "demo",
			Password:    "supersecret",
			Type:        "kerberos",
			ID:          "3c829734-0679-4c17-9637-801da48c0d5f",
		},
		{
			Status:      "new",
			Domain:      "",
			ProjectID:   "16e1ab15c35a457e9c2b2aa189f544e1",
			Name:        "SecServ2",
			CreatedAt:   time.Date(2015, 9, 7, 12, 25, 03, 0, time.UTC),
			Description: "Creating my second Security Service",
			UpdatedAt:   nilTime,
			Server:      "",
			DNSIP:       "10.0.0.0/24",
			User:        "",
			Password:    "",
			Type:        "ldap",
			ID:          "5a1d3a12-34a7-4087-8983-50e9ed03509a",
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}

// Verifies that security services list can be called with query parameters
func TestFilteredList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockFilteredListResponse(t, fakeServer)

	options := &securityservices.ListOpts{
		Type: "kerberos",
	}

	allPages, err := securityservices.List(client.ServiceClient(fakeServer), options).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := securityservices.ExtractSecurityServices(allPages)
	th.AssertNoErr(t, err)
	var nilTime time.Time
	expected := []securityservices.SecurityService{
		{
			Status:      "new",
			Domain:      "",
			ProjectID:   "16e1ab15c35a457e9c2b2aa189f544e1",
			Name:        "SecServ1",
			CreatedAt:   time.Date(2015, 9, 7, 12, 19, 10, 0, time.UTC),
			Description: "Creating my first Security Service",
			UpdatedAt:   nilTime,
			Server:      "",
			DNSIP:       "10.0.0.0/24",
			User:        "demo",
			Password:    "supersecret",
			Type:        "kerberos",
			ID:          "3c829734-0679-4c17-9637-801da48c0d5f",
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}

// Verifies that it is possible to get a security service
func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetResponse(t, fakeServer)

	var nilTime time.Time
	expected := securityservices.SecurityService{
		ID:          "3c829734-0679-4c17-9637-801da48c0d5f",
		Name:        "SecServ1",
		CreatedAt:   time.Date(2015, 9, 7, 12, 19, 10, 0, time.UTC),
		Description: "Creating my first Security Service",
		Type:        "kerberos",
		UpdatedAt:   nilTime,
		ProjectID:   "16e1ab15c35a457e9c2b2aa189f544e1",
		Status:      "new",
		Domain:      "",
		Server:      "",
		DNSIP:       "10.0.0.0/24",
		User:        "demo",
		Password:    "supersecret",
	}

	n, err := securityservices.Get(context.TODO(), client.ServiceClient(fakeServer), "3c829734-0679-4c17-9637-801da48c0d5f").Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, &expected, n)
}

// Verifies that it is possible to update a security service
func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockUpdateResponse(t, fakeServer)
	expected := securityservices.SecurityService{
		ID:          "securityServiceID",
		Name:        "SecServ2",
		CreatedAt:   time.Date(2015, 9, 7, 12, 19, 10, 0, time.UTC),
		Description: "Updating my first Security Service",
		Type:        "kerberos",
		UpdatedAt:   time.Date(2015, 9, 7, 12, 20, 10, 0, time.UTC),
		ProjectID:   "16e1ab15c35a457e9c2b2aa189f544e1",
		Status:      "new",
		Domain:      "",
		Server:      "",
		DNSIP:       "10.0.0.0/24",
		User:        "demo",
		Password:    "supersecret",
	}

	name := "SecServ2"
	options := securityservices.UpdateOpts{Name: &name}
	s, err := securityservices.Update(context.TODO(), client.ServiceClient(fakeServer), "securityServiceID", options).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &expected, s)
}
