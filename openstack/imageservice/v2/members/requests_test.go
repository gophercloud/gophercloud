package members

import (
	"strings"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fakeclient "github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateMemberSuccessfully(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateImageMemberSuccessfully(t)
	im, err := Create(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, ImageMember{
		CreatedAt: "2013-09-20T19:22:19Z",
		ImageID:   "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		MemberID:  "8989447062e04a818baf9e073fd04fa7",
		Schema:    "/v2/schemas/member",
		Status:    "pending",
		UpdatedAt: "2013-09-20T19:25:31Z",
	}, *im)

}

func TestCreateMemberMemberConflict(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateImageMemberConflict(t)

	result := Create(fakeclient.ServiceClient(), "da3b75d9-memberConflict",
		"8989447062e04a818baf9e073fd04fa7")

	if result.Err == nil {
		t.Fatalf("Expected error in result defined (Err: %v)", result.Err)
	}

	message := result.Err.Error()
	if !strings.Contains(message, "is already member for image") {
		t.Fatalf("Wrong error message: %s", message)
	}

}
func TestCreateMemberInvalidVisibility(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateImageMemberInvalidVisibility(t)

	result := Create(fakeclient.ServiceClient(), "da3b75d9-invalid-visibility",
		"8989447062e04a818baf9e073fd04fa7")

	if result.Err == nil {
		t.Fatalf("Expected error in result defined (Err: %v)", result.Err)
	}

	message := result.Err.Error()
	if !strings.Contains(message, "which 'visibility' attribute is private") {
		t.Fatalf("Wrong error message: %s", message)
	}
}

func TestMemberListSuccessfully(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageMemberList(t)

	images, err := List(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea").Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, images)
	th.AssertEquals(t, 2, len(*images))
}

func TestMemberListEmpty(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageMemberEmptyList(t)

	images, err := List(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea").Extract()
	th.AssertNoErr(t, err)
	th.AssertNotNil(t, images)
	th.AssertEquals(t, 0, len(*images))
}

func TestShowMemberDetails(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageMemberDetails(t)
	md, err := Get(fakeclient.ServiceClient(),
		"da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7").Extract()

	th.AssertNoErr(t, err)
	th.AssertNotNil(t, md)

	th.AssertDeepEquals(t, ImageMember{
		CreatedAt: "2013-11-26T07:21:21Z",
		ImageID:   "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		MemberID:  "8989447062e04a818baf9e073fd04fa7",
		Schema:    "/v2/schemas/member",
		Status:    "pending",
		UpdatedAt: "2013-11-26T07:21:21Z",
	}, *md)
}

func TestDeleteMember(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	counter := HandleImageMemberDeleteSuccessfully(t)

	result := Delete(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7")
	th.AssertEquals(t, 1, counter.Counter)
	th.AssertNoErr(t, result.Err)
}

func TestDeleteMemberByNonOwner(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	counter := HandleImageMemberDeleteByNonOwner(t)

	result := Delete(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7")
	th.AssertEquals(t, 1, counter.Counter)

	if result.Err == nil {
		t.Fatalf("Expected error in result defined (Err: %v)", result.Err)
	}

	message := result.Err.Error()
	if !strings.Contains(message, "You must be the owner of the specified image") {
		t.Fatalf("Wrong error message: %s", message)
	}
}

func TestMemberUpdateSuccessfully(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	counter := HandleImageMemberUpdate(t)
	im, err := Update(fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		"8989447062e04a818baf9e073fd04fa7", "accepted").Extract()
	th.AssertEquals(t, 1, counter.Counter)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, ImageMember{
		CreatedAt: "2013-11-26T07:21:21Z",
		ImageID:   "da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		MemberID:  "8989447062e04a818baf9e073fd04fa7",
		Schema:    "/v2/schemas/member",
		Status:    "accepted",
		UpdatedAt: "2013-11-26T07:21:21Z",
	}, *im)

}
