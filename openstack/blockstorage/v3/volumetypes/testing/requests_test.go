package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allPages, err := volumetypes.List(client.ServiceClient()).AllPages()
	th.AssertNoErr(t, err)
	actual, err := volumetypes.ExtractVolumeTypes(allPages)
	th.AssertNoErr(t, err)

	expected := []volumetypes.VolumeType{
		{
			ID:          "6685584b-1eac-4da6-b5c3-555430cf68ff",
			Name:        "SSD",
			ExtraSpecs:  map[string]string{"volume_backend_name": "lvmdriver-1"},
			IsPublic:    true,
			Description: "",
		}, {
			ID:          "8eb69a46-df97-4e41-9586-9a40a7533803",
			Name:        "SATA",
			ExtraSpecs:  map[string]string{"volume_backend_name": "lvmdriver-1"},
			IsPublic:    true,
			Description: "",
		},
	}
	th.CheckDeepEquals(t, expected, actual)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	v, err := volumetypes.Get(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, v.Name, "vol-type-001")
	th.AssertEquals(t, v.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertEquals(t, v.ExtraSpecs["capabilities"], "gpu")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &volumetypes.CreateOpts{Name: "test_type", IsPublic: true, Description: "test_type_desc"}
	n, err := volumetypes.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Name, "test_type")
	th.AssertEquals(t, n.Description, "test_type_desc")
	th.AssertEquals(t, n.IsPublic, true)
	th.AssertEquals(t, n.ID, "6d0ff92a-0007-4780-9ece-acfe5876966a")
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := volumetypes.Delete(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateResponse(t)

	options := volumetypes.UpdateOpts{Name: "vol-type-002"}
	v, err := volumetypes.Update(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", options).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "vol-type-002", v.Name)
}
