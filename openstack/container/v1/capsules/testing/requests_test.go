package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/container/v1/capsules"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetCapsule_OldTime(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCapsuleGetOldTimeSuccessfully(t, fakeServer)

	createdAt, _ := time.Parse(gophercloud.RFC3339ZNoT, "2018-01-12 09:37:25+00:00")
	updatedAt, _ := time.Parse(gophercloud.RFC3339ZNoT, "2018-01-12 09:37:26+00:00")
	startedAt, _ := time.Parse(gophercloud.RFC3339ZNoT, "2018-01-12 09:37:26+00:00")

	ec := GetFakeCapsule()
	ec.CreatedAt = createdAt
	ec.UpdatedAt = updatedAt
	ec.Containers[0].CreatedAt = createdAt
	ec.Containers[0].UpdatedAt = updatedAt
	ec.Containers[0].StartedAt = startedAt

	actualCapsule, err := capsules.Get(context.TODO(), client.ServiceClient(fakeServer), ec.UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &ec, actualCapsule)
}

func TestGetCapsule_NewTime(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCapsuleGetNewTimeSuccessfully(t, fakeServer)

	ec := GetFakeCapsule()

	actualCapsule, err := capsules.Get(context.TODO(), client.ServiceClient(fakeServer), ec.UUID).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &ec, actualCapsule)
}

func TestCreateCapsule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCapsuleCreateSuccessfully(t, fakeServer)

	ec := GetFakeCapsule()

	template := new(capsules.Template)
	template.Bin = []byte(ValidJSONTemplate)
	createOpts := capsules.CreateOpts{
		TemplateOpts: template,
	}
	actualCapsule, err := capsules.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &ec, actualCapsule)
}

func TestListCapsule(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCapsuleListSuccessfully(t, fakeServer)

	createdAt, _ := time.Parse(gophercloud.RFC3339ZNoT, "2018-01-12 09:37:25+00:00")
	updatedAt, _ := time.Parse(gophercloud.RFC3339ZNoT, "2018-01-12 09:37:25+01:00")

	ec := GetFakeCapsule()
	ec.CreatedAt = createdAt
	ec.UpdatedAt = updatedAt
	ec.Containers = nil

	expected := []capsules.Capsule{ec}

	count := 0
	results := capsules.List(client.ServiceClient(fakeServer), nil)
	err := results.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := capsules.ExtractCapsules(page)
		if err != nil {
			t.Errorf("Failed to extract capsules: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListCapsuleV132(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCapsuleV132ListSuccessfully(t, fakeServer)

	createdAt, _ := time.Parse(gophercloud.RFC3339ZNoTNoZ, "2018-01-12 09:37:25")
	updatedAt, _ := time.Parse(gophercloud.RFC3339ZNoTNoZ, "2018-01-12 09:37:25")

	ec := GetFakeCapsuleV132()
	ec.CreatedAt = createdAt
	ec.UpdatedAt = updatedAt
	ec.Containers = nil

	expected := []capsules.CapsuleV132{ec}

	count := 0
	results := capsules.List(client.ServiceClient(fakeServer), nil)
	err := results.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := capsules.ExtractCapsules(page)
		if err != nil {
			t.Errorf("Failed to extract capsules: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCapsuleDeleteSuccessfully(t, fakeServer)

	res := capsules.Delete(context.TODO(), client.ServiceClient(fakeServer), "963a239d-3946-452b-be5a-055eab65a421")
	th.AssertNoErr(t, res.Err)
}
