package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreateOpts(t *testing.T) {
	base := servers.CreateOpts{
		Name:      "createdserver",
		ImageRef:  "asdfasdfasdf",
		FlavorRef: "performance1-1",
	}

	ext := bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: base,
		BlockDevice: []bootfromvolume.BlockDevice{
			{
				UUID:                "123456",
				SourceType:          bootfromvolume.Image,
				DestinationType:     "volume",
				VolumeSize:          10,
				DeleteOnTermination: false,
			},
		},
	}

	expected := `
    {
      "server": {
        "name": "createdserver",
        "imageRef": "asdfasdfasdf",
        "flavorRef": "performance1-1",
        "block_device_mapping_v2":[
          {
            "uuid":"123456",
            "source_type":"image",
            "destination_type":"volume",
            "boot_index": 0,
            "delete_on_termination": false,
            "volume_size": 10
          }
        ]
      }
    }
  `
	actual, err := ext.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestCreateMultiEphemeralOpts(t *testing.T) {
	base := servers.CreateOpts{
		Name:      "createdserver",
		ImageRef:  "asdfasdfasdf",
		FlavorRef: "performance1-1",
	}

	ext := bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: base,
		BlockDevice: []bootfromvolume.BlockDevice{
			{
				BootIndex:           0,
				DeleteOnTermination: true,
				DestinationType:     "local",
				SourceType:          bootfromvolume.Image,
				UUID:                "123456",
			},
			{
				BootIndex:           -1,
				DeleteOnTermination: true,
				DestinationType:     "local",
				GuestFormat:         "ext4",
				SourceType:          bootfromvolume.Blank,
				VolumeSize:          1,
			},
			{
				BootIndex:           -1,
				DeleteOnTermination: true,
				DestinationType:     "local",
				GuestFormat:         "ext4",
				SourceType:          bootfromvolume.Blank,
				VolumeSize:          1,
			},
		},
	}

	expected := `
    {
      "server": {
        "name": "createdserver",
        "imageRef": "asdfasdfasdf",
        "flavorRef": "performance1-1",
        "block_device_mapping_v2":[
          {
            "boot_index": 0,
            "delete_on_termination": true,
            "destination_type":"local",
            "source_type":"image",
            "uuid":"123456",
            "volume_size": 0
          },
          {
            "boot_index": -1,
            "delete_on_termination": true,
            "destination_type":"local",
            "guest_format":"ext4",
            "source_type":"blank",
            "volume_size": 1
          },
          {
            "boot_index": -1,
            "delete_on_termination": true,
            "destination_type":"local",
            "guest_format":"ext4",
            "source_type":"blank",
            "volume_size": 1
          }
        ]
      }
    }
  `
	actual, err := ext.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}
