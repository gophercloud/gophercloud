package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestChooseVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	setupIdentityVersionHandler(fakeServer)

	v2 := &utils.Version{ID: "v2.0", Priority: 2, Suffix: "blarg"}
	v3 := &utils.Version{ID: "v3.0", Priority: 3, Suffix: "hargl"}

	c := &gophercloud.ProviderClient{
		IdentityBase:     fakeServer.Endpoint(),
		IdentityEndpoint: "",
	}
	v, endpoint, err := utils.ChooseVersion(context.TODO(), c, []*utils.Version{v2, v3}) //nolint

	if err != nil {
		t.Fatalf("Unexpected error from ChooseVersion: %v", err)
	}

	if v != v3 {
		t.Errorf("Expected %#v to win, but %#v did instead", v3, v)
	}

	expected := fakeServer.Endpoint() + "v3.0/"
	if endpoint != expected {
		t.Errorf("Expected endpoint [%s], but was [%s] instead", expected, endpoint)
	}
}

func TestChooseVersionOpinionatedLink(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	setupIdentityVersionHandler(fakeServer)

	v2 := &utils.Version{ID: "v2.0", Priority: 2, Suffix: "nope"}
	v3 := &utils.Version{ID: "v3.0", Priority: 3, Suffix: "northis"}

	c := &gophercloud.ProviderClient{
		IdentityBase:     fakeServer.Endpoint(),
		IdentityEndpoint: fakeServer.Endpoint() + "v2.0/",
	}
	v, endpoint, err := utils.ChooseVersion(context.TODO(), c, []*utils.Version{v2, v3}) //nolint
	if err != nil {
		t.Fatalf("Unexpected error from ChooseVersion: %v", err)
	}

	if v != v2 {
		t.Errorf("Expected %#v to win, but %#v did instead", v2, v)
	}

	expected := fakeServer.Endpoint() + "v2.0/"
	if endpoint != expected {
		t.Errorf("Expected endpoint [%s], but was [%s] instead", expected, endpoint)
	}
}

func TestChooseVersionFromSuffix(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	v2 := &utils.Version{ID: "v2.0", Priority: 2, Suffix: "/v2.0/"}
	v3 := &utils.Version{ID: "v3.0", Priority: 3, Suffix: "/v3.0/"}

	c := &gophercloud.ProviderClient{
		IdentityBase:     fakeServer.Endpoint(),
		IdentityEndpoint: fakeServer.Endpoint() + "v2.0/",
	}
	v, endpoint, err := utils.ChooseVersion(context.TODO(), c, []*utils.Version{v2, v3}) //nolint
	if err != nil {
		t.Fatalf("Unexpected error from ChooseVersion: %v", err)
	}

	if v != v2 {
		t.Errorf("Expected %#v to win, but %#v did instead", v2, v)
	}

	expected := fakeServer.Endpoint() + "v2.0/"
	if endpoint != expected {
		t.Errorf("Expected endpoint [%s], but was [%s] instead", expected, endpoint)
	}
}
