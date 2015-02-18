package backups

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/db/v1/datastores"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/gophercloud/testhelper/fixture"
)

var (
	backupID = "{backupID}"
	_rootURL = "/backups"
	resURL   = _rootURL + "/" + backupID
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, _rootURL, "POST", createReq, createResp, 202)

	opts := CreateOpts{
		Name:        "snapshot",
		Description: "My Backup",
		InstanceID:  "d4603f69-ec7e-4e9b-803f-600b9205576f",
	}

	instance, err := Create(fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	expected := &Backup{
		Created:     "2014-02-13T21:47:16",
		Description: "My Backup",
		ID:          "61f12fef-edb1-4561-8122-e7c00ef26a82",
		InstanceID:  "d4603f69-ec7e-4e9b-803f-600b9205576f",
		LocationRef: "",
		Name:        "snapshot",
		ParentID:    "",
		Size:        100,
		Status:      "NEW",
		Updated:     "2014-02-13T21:47:16",
		Datastore: datastores.DatastorePartial{
			Version:   "5.1",
			Type:      "MySQL",
			VersionID: "20000000-0000-0000-0000-000000000002",
		},
	}

	th.AssertDeepEquals(t, expected, instance)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, _rootURL, "GET", "", listResp, 200)

	pages := 0

	err := List(fake.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		pages++
		actual, err := ExtractBackups(page)
		th.AssertNoErr(t, err)

		expected := []Backup{
			Backup{
				Created:     "2014-06-18T21:23:35",
				Description: "Backup from Restored Instance",
				ID:          "87972694-4be2-40f5-83f8-501656e0032a",
				InstanceID:  "29af2cd9-0674-48ab-b87a-b160f00208e6",
				LocationRef: "http://localhost/path/to/backup",
				Name:        "restored_backup",
				ParentID:    "",
				Size:        0.141026,
				Status:      "COMPLETED",
				Updated:     "2014-06-18T21:24:39",
				Datastore: datastores.DatastorePartial{
					Version:   "5.1",
					Type:      "MySQL",
					VersionID: "20000000-0000-0000-0000-000000000002",
				},
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, resURL, "GET", "", getResp, 200)

	instance, err := Get(fake.ServiceClient(), backupID).Extract()
	th.AssertNoErr(t, err)

	expected := &Backup{
		Created:     "2014-02-13T21:47:16",
		Description: "My Backup",
		ID:          "61f12fef-edb1-4561-8122-e7c00ef26a82",
		InstanceID:  "d4603f69-ec7e-4e9b-803f-600b9205576f",
		LocationRef: "",
		Name:        "snapshot",
		ParentID:    "",
		Size:        100,
		Status:      "NEW",
		Updated:     "2014-02-13T21:47:16",
		Datastore: datastores.DatastorePartial{
			Version:   "5.1",
			Type:      "MySQL",
			VersionID: "20000000-0000-0000-0000-000000000002",
		},
	}

	th.AssertDeepEquals(t, expected, instance)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, resURL, "DELETE", "", "", 202)

	err := Delete(fake.ServiceClient(), backupID).ExtractErr()
	th.AssertNoErr(t, err)
}
