package testing

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)
	pages := 0
	err := volumetypes.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		pages++
		actual, err := volumetypes.ExtractVolumeTypes(page)
		if err != nil {
			return false, err
		}
		expected := []volumetypes.VolumeType{
			{
				ID:           "6685584b-1eac-4da6-b5c3-555430cf68ff",
				Name:         "SSD",
				ExtraSpecs:   map[string]string{"volume_backend_name": "lvmdriver-1"},
				IsPublic:     true,
				Description:  "",
				QosSpecID:    "",
				PublicAccess: true,
			}, {
				ID:           "8eb69a46-df97-4e41-9586-9a40a7533803",
				Name:         "SATA",
				ExtraSpecs:   map[string]string{"volume_backend_name": "lvmdriver-1"},
				IsPublic:     true,
				Description:  "",
				QosSpecID:    "",
				PublicAccess: true,
			},
		}
		th.CheckDeepEquals(t, expected, actual)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, pages, 1)
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
	th.AssertEquals(t, v.QosSpecID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertEquals(t, v.PublicAccess, true)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	var isPublic = true

	options := &volumetypes.CreateOpts{
		Name:        "test_type",
		IsPublic:    &isPublic,
		Description: "test_type_desc",
		ExtraSpecs:  map[string]string{"capabilities": "gpu"},
	}

	n, err := volumetypes.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Name, "test_type")
	th.AssertEquals(t, n.Description, "test_type_desc")
	th.AssertEquals(t, n.IsPublic, true)
	th.AssertEquals(t, n.PublicAccess, true)
	th.AssertEquals(t, n.ID, "6d0ff92a-0007-4780-9ece-acfe5876966a")
	th.AssertEquals(t, n.ExtraSpecs["capabilities"], "gpu")
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

	var isPublic = true
	var name = "vol-type-002"
	options := volumetypes.UpdateOpts{
		Name:     &name,
		IsPublic: &isPublic,
	}

	v, err := volumetypes.Update(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", options).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "vol-type-002", v.Name)
	th.CheckEquals(t, true, v.IsPublic)
}

func TestVolumeTypeExtraSpecsList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleExtraSpecsListSuccessfully(t)

	expected := ExtraSpecs
	actual, err := volumetypes.ListExtraSpecs(client.ServiceClient(), "1").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}

func TestVolumeTypeExtraSpecGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleExtraSpecGetSuccessfully(t)

	expected := ExtraSpec
	actual, err := volumetypes.GetExtraSpec(client.ServiceClient(), "1", "capabilities").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}

func TestVolumeTypeExtraSpecsCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleExtraSpecsCreateSuccessfully(t)

	createOpts := volumetypes.ExtraSpecsOpts{
		"capabilities":        "gpu",
		"volume_backend_name": "ssd",
	}
	expected := ExtraSpecs
	actual, err := volumetypes.CreateExtraSpecs(client.ServiceClient(), "1", createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}

func TestVolumeTypeExtraSpecUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleExtraSpecUpdateSuccessfully(t)

	updateOpts := volumetypes.ExtraSpecsOpts{
		"capabilities": "gpu-2",
	}
	expected := UpdatedExtraSpec
	actual, err := volumetypes.UpdateExtraSpec(client.ServiceClient(), "1", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}

func TestVolumeTypeExtraSpecDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleExtraSpecDeleteSuccessfully(t)

	res := volumetypes.DeleteExtraSpec(client.ServiceClient(), "1", "capabilities")
	th.AssertNoErr(t, res.Err)
}

func TestVolumeTypeListAccesses(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/os-volume-type-access", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
			  "volume_type_access": [
			    {
            	  "project_id": "6f70656e737461636b20342065766572",
            	  "volume_type_id": "a5082c24-2a27-43a4-b48e-fcec1240e36b"
			    }
			  ]
			}
		`)
	})

	expected := []volumetypes.VolumeTypeAccess{
		{
			VolumeTypeID: "a5082c24-2a27-43a4-b48e-fcec1240e36b",
			ProjectID:    "6f70656e737461636b20342065766572",
		},
	}

	allPages, err := volumetypes.ListAccesses(client.ServiceClient(), "a5082c24-2a27-43a4-b48e-fcec1240e36b").AllPages()
	th.AssertNoErr(t, err)

	actual, err := volumetypes.ExtractAccesses(allPages)
	th.AssertNoErr(t, err)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestVolumeTypeAddAccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "accept", "application/json")
		th.TestJSONRequest(t, r, `
			{
			  "addProjectAccess": {
			    "project": "6f70656e737461636b20342065766572"
			  }
			}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})

	addAccessOpts := volumetypes.AddAccessOpts{
		Project: "6f70656e737461636b20342065766572",
	}

	err := volumetypes.AddAccess(client.ServiceClient(), "a5082c24-2a27-43a4-b48e-fcec1240e36b", addAccessOpts).ExtractErr()
	th.AssertNoErr(t, err)

}

func TestVolumeTypeRemoveAccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "accept", "application/json")
		th.TestJSONRequest(t, r, `
			{
			  "removeProjectAccess": {
			    "project": "6f70656e737461636b20342065766572"
			  }
			}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})

	removeAccessOpts := volumetypes.RemoveAccessOpts{
		Project: "6f70656e737461636b20342065766572",
	}

	err := volumetypes.RemoveAccess(client.ServiceClient(), "a5082c24-2a27-43a4-b48e-fcec1240e36b", removeAccessOpts).ExtractErr()
	th.AssertNoErr(t, err)

}

func TestCreateEncryption(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockEncryptionCreateResponse(t)

	options := &volumetypes.CreateEncryptionOpts{
		KeySize:         256,
		Provider:        "luks",
		ControlLocation: "front-end",
		Cipher:          "aes-xts-plain64",
	}
	id := "a5082c24-2a27-43a4-b48e-fcec1240e36b"
	n, err := volumetypes.CreateEncryption(client.ServiceClient(), id, options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "a5082c24-2a27-43a4-b48e-fcec1240e36b", n.VolumeTypeID)
	th.AssertEquals(t, "front-end", n.ControlLocation)
	th.AssertEquals(t, "81e069c6-7394-4856-8df7-3b237ca61f74", n.EncryptionID)
	th.AssertEquals(t, 256, n.KeySize)
	th.AssertEquals(t, "luks", n.Provider)
	th.AssertEquals(t, "aes-xts-plain64", n.Cipher)
}

func TestDeleteEncryption(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteEncryptionResponse(t)

	res := volumetypes.DeleteEncryption(client.ServiceClient(), "a5082c24-2a27-43a4-b48e-fcec1240e36b", "81e069c6-7394-4856-8df7-3b237ca61f74")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateEncryption(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockEncryptionUpdateResponse(t)

	options := &volumetypes.UpdateEncryptionOpts{
		KeySize:         256,
		Provider:        "luks",
		ControlLocation: "front-end",
		Cipher:          "aes-xts-plain64",
	}
	id := "a5082c24-2a27-43a4-b48e-fcec1240e36b"
	encryptionID := "81e069c6-7394-4856-8df7-3b237ca61f74"
	n, err := volumetypes.UpdateEncryption(client.ServiceClient(), id, encryptionID, options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "front-end", n.ControlLocation)
	th.AssertEquals(t, 256, n.KeySize)
	th.AssertEquals(t, "luks", n.Provider)
	th.AssertEquals(t, "aes-xts-plain64", n.Cipher)
}

func TestGetEncryption(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockEncryptionGetResponse(t)
	id := "a5082c24-2a27-43a4-b48e-fcec1240e36b"
	n, err := volumetypes.GetEncryption(client.ServiceClient(), id).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "a5082c24-2a27-43a4-b48e-fcec1240e36b", n.VolumeTypeID)
	th.AssertEquals(t, "front-end", n.ControlLocation)
	th.AssertEquals(t, false, n.Deleted)
	th.AssertEquals(t, "2016-12-28T02:32:25.000000", n.CreatedAt)
	th.AssertEquals(t, "", n.UpdatedAt)
	th.AssertEquals(t, "81e069c6-7394-4856-8df7-3b237ca61f74", n.EncryptionID)
	th.AssertEquals(t, 256, n.KeySize)
	th.AssertEquals(t, "luks", n.Provider)
	th.AssertEquals(t, "", n.DeletedAt)
	th.AssertEquals(t, "aes-xts-plain64", n.Cipher)
}

func TestGetEncryptionSpec(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockEncryptionGetSpecResponse(t)
	id := "a5082c24-2a27-43a4-b48e-fcec1240e36b"
	n, err := volumetypes.GetEncryptionSpec(client.ServiceClient(), id, "cipher").Extract()
	th.AssertNoErr(t, err)

	key := "cipher"
	testVar, exists := n[key]
	if exists {
		th.AssertEquals(t, "aes-xts-plain64", testVar)
	} else {
		t.Fatalf("Key %s does not exist in map.", key)
	}
}
