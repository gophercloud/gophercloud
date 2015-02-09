package instances

import (
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateInstanceSuccessfully(t)

	opts := CreateOpts{
		Name:      "json_rack_instance",
		FlavorRef: "1",
		Databases: DatabasesOpts{
			DatabaseOpts{CharSet: "utf8", Collate: "utf8_general_ci", Name: "sampledb"},
			DatabaseOpts{Name: "nextround"},
		},
		Users: UsersOpts{
			UserOpts{
				Name:     "demouser",
				Password: "demopassword",
				Databases: DatabasesOpts{
					DatabaseOpts{Name: "sampledb"},
				},
			},
		},
		Size:         2,
		RestorePoint: "1234567890",
	}

	_ = Create(fake.ServiceClient(), opts)
}
