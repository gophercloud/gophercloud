package testing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListServers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerListSuccessfully(t, fakeServer)

	pages := 0
	err := servers.List(client.ServiceClient(fakeServer), servers.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 3 {
			t.Fatalf("Expected 3 servers, got %d", len(actual))
		}
		th.CheckDeepEquals(t, ServerHerp, actual[0])
		th.CheckDeepEquals(t, ServerDerp, actual[1])
		th.CheckDeepEquals(t, ServerMerp, actual[2])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllServers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerListSimpleSuccessfully(t, fakeServer)

	allPages, err := servers.ListSimple(client.ServiceClient(fakeServer), servers.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := servers.ExtractServers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ServerHerp, actual[0])
	th.CheckDeepEquals(t, ServerDerp, actual[1])
}

func TestListAllServersWithExtensions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerListSuccessfully(t, fakeServer)

	allPages, err := servers.List(client.ServiceClient(fakeServer), servers.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	var actual []servers.Server
	err = servers.ExtractServersInto(allPages, &actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 3, len(actual))
	th.AssertEquals(t, "nova", actual[0].AvailabilityZone)
	th.AssertEquals(t, "RUNNING", actual[0].PowerState.String())
	th.AssertEquals(t, "", actual[0].TaskState)
	th.AssertEquals(t, "active", actual[0].VmState)
	th.AssertEquals(t, servers.Manual, actual[0].DiskConfig)
}

func TestCreateServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerCreationSuccessfully(t, fakeServer, SingleServerBody)

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServerNoNetwork(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerNoNetworkCreationSuccessfully(t, fakeServer, SingleServerBody)

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
		Networks:  "none",
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServersCreationSuccessfully(t, fakeServer, SingleServerBody)

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
		Min:       3,
		Max:       3,
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServerWithCustomField(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerCreationWithCustomFieldSuccessfully(t, fakeServer, SingleServerBody)

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), CreateOptsWithCustomField{
		CreateOpts: servers.CreateOpts{
			Name:      "derp",
			ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
			FlavorRef: "1",
		},
		Foo: "bar",
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServerWithMetadata(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerCreationWithMetadata(t, fakeServer, SingleServerBody)

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
		Metadata: map[string]string{
			"abc": "def",
		},
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServerWithUserdataString(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerCreationWithUserdata(t, fakeServer, SingleServerBody)

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
		UserData:  []byte("userdata string"),
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServerWithUserdataEncoded(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerCreationWithUserdata(t, fakeServer, SingleServerBody)

	encoded := base64.StdEncoding.EncodeToString([]byte("userdata string"))

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
		UserData:  []byte(encoded),
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServerWithHostname(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerCreationWithHostname(t, fakeServer, SingleServerBody)

	actual, err := servers.Create(context.TODO(), client.ServiceClient(fakeServer), servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
		Hostname:  "derp.local",
	}, nil).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestCreateServerWithDiskConfig(t *testing.T) {
	opts := servers.CreateOpts{
		Name:       "createdserver",
		ImageRef:   "asdfasdfasdf",
		FlavorRef:  "performance1-1",
		DiskConfig: servers.Manual,
	}
	expected := `
		{
			"server": {
				"name": "createdserver",
				"imageRef": "asdfasdfasdf",
				"flavorRef": "performance1-1",
				"OS-DCF:diskConfig": "MANUAL"
			}
		}
	`

	actual, err := opts.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestCreateServerWithBFVBootFromNewVolume(t *testing.T) {
	opts := servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		BlockDevice: []servers.BlockDevice{
			{
				UUID:                "123456",
				SourceType:          servers.SourceImage,
				DestinationType:     servers.DestinationVolume,
				VolumeSize:          10,
				DeleteOnTermination: true,
			},
		},
	}
	expected := `
	{
		"server": {
			"name":"createdserver",
			"flavorRef":"performance1-1",
			"imageRef":"",
			"block_device_mapping_v2":[
				{
					"uuid":"123456",
					"source_type":"image",
					"destination_type":"volume",
					"boot_index": 0,
					"delete_on_termination": true,
					"volume_size": 10
				}
			]
		}
	}
	`

	actual, err := opts.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestCreateServerWithBFVBootFromExistingVolume(t *testing.T) {
	opts := servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		BlockDevice: []servers.BlockDevice{
			{
				UUID:                "123456",
				SourceType:          servers.SourceVolume,
				DestinationType:     servers.DestinationVolume,
				DeleteOnTermination: true,
			},
		},
	}
	expected := `
	{
		"server": {
			"name":"createdserver",
			"flavorRef":"performance1-1",
			"imageRef":"",
			"block_device_mapping_v2":[
				{
					"uuid":"123456",
					"source_type":"volume",
					"destination_type":"volume",
					"boot_index": 0,
					"delete_on_termination": true
				}
			]
		}
	}
	`

	actual, err := opts.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestCreateServerWithBFVBootFromImage(t *testing.T) {
	var ImageRequest = servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		ImageRef:  "asdfasdfasdf",
		BlockDevice: []servers.BlockDevice{
			{
				BootIndex:           0,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationLocal,
				SourceType:          servers.SourceImage,
				UUID:                "asdfasdfasdf",
			},
		},
	}
	const ExpectedImageRequest = `
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
					"uuid":"asdfasdfasdf"
				}
			]
		}
	}
	`

	actual, err := ImageRequest.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, ExpectedImageRequest, actual)
}

func TestCreateServerWithBFVCreateMultiEphemeralOpts(t *testing.T) {
	var MultiEphemeralRequest = servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		ImageRef:  "asdfasdfasdf",
		BlockDevice: []servers.BlockDevice{
			{
				BootIndex:           0,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationLocal,
				SourceType:          servers.SourceImage,
				UUID:                "asdfasdfasdf",
			},
			{
				BootIndex:           -1,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationLocal,
				GuestFormat:         "ext4",
				SourceType:          servers.SourceBlank,
				VolumeSize:          1,
			},
			{
				BootIndex:           -1,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationLocal,
				GuestFormat:         "ext4",
				SourceType:          servers.SourceBlank,
				VolumeSize:          1,
			},
		},
	}
	const ExpectedMultiEphemeralRequest = `
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
					"uuid":"asdfasdfasdf"
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

	actual, err := MultiEphemeralRequest.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, ExpectedMultiEphemeralRequest, actual)
}

func TestCreateServerWithBFVAttachNewVolume(t *testing.T) {
	opts := servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		ImageRef:  "asdfasdfasdf",
		BlockDevice: []servers.BlockDevice{
			{
				BootIndex:           0,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationLocal,
				SourceType:          servers.SourceImage,
				UUID:                "asdfasdfasdf",
			},
			{
				BootIndex:           1,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationVolume,
				SourceType:          servers.SourceBlank,
				VolumeSize:          1,
				DeviceType:          "disk",
				DiskBus:             "scsi",
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
					"uuid":"asdfasdfasdf"
				},
				{
					"boot_index": 1,
					"delete_on_termination": true,
					"destination_type":"volume",
					"source_type":"blank",
					"volume_size": 1,
					"device_type": "disk",
					"disk_bus": "scsi"
				}
			]
		}
	}
	`

	actual, err := opts.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestCreateServerWithBFVAttachExistingVolume(t *testing.T) {
	opts := servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		ImageRef:  "asdfasdfasdf",
		BlockDevice: []servers.BlockDevice{
			{
				BootIndex:           0,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationLocal,
				SourceType:          servers.SourceImage,
				UUID:                "asdfasdfasdf",
			},
			{
				BootIndex:           1,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationVolume,
				SourceType:          servers.SourceVolume,
				UUID:                "123456",
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
					"uuid":"asdfasdfasdf"
				},
				{
					"boot_index": 1,
					"delete_on_termination": true,
					"destination_type":"volume",
					"source_type":"volume",
					"uuid":"123456",
					"volume_size": 1
				}
			]
		}
	}
	`

	actual, err := opts.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestCreateServerWithBFVBootFromNewVolumeType(t *testing.T) {
	var NewVolumeTypeRequest = servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		BlockDevice: []servers.BlockDevice{
			{
				UUID:                "123456",
				SourceType:          servers.SourceImage,
				DestinationType:     servers.DestinationVolume,
				VolumeSize:          10,
				DeleteOnTermination: true,
				VolumeType:          "ssd",
			},
		},
	}
	const ExpectedNewVolumeTypeRequest = `
	{
		"server": {
			"name":"createdserver",
			"flavorRef":"performance1-1",
			"imageRef":"",
			"block_device_mapping_v2":[
				{
					"uuid":"123456",
					"source_type":"image",
					"destination_type":"volume",
					"boot_index": 0,
					"delete_on_termination": true,
					"volume_size": 10,
					"volume_type": "ssd"
				}
			]
		}
	}
	`

	actual, err := NewVolumeTypeRequest.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, ExpectedNewVolumeTypeRequest, actual)
}

func TestCreateServerWithBFVAttachExistingVolumeWithTag(t *testing.T) {
	var ImageAndExistingVolumeWithTagRequest = servers.CreateOpts{
		Name:      "createdserver",
		FlavorRef: "performance1-1",
		ImageRef:  "asdfasdfasdf",
		BlockDevice: []servers.BlockDevice{
			{
				BootIndex:           0,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationLocal,
				SourceType:          servers.SourceImage,
				UUID:                "asdfasdfasdf",
			},
			{
				BootIndex:           -1,
				DeleteOnTermination: true,
				DestinationType:     servers.DestinationVolume,
				SourceType:          servers.SourceVolume,
				Tag:                 "volume-tag",
				UUID:                "123456",
				VolumeSize:          1,
			},
		},
	}
	const ExpectedImageAndExistingVolumeWithTagRequest = `
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
					"uuid":"asdfasdfasdf"
				},
				{
					"boot_index": -1,
					"delete_on_termination": true,
					"destination_type":"volume",
					"source_type":"volume",
					"tag": "volume-tag",
					"uuid":"123456",
					"volume_size": 1
				}
			]
		}
	}
	`

	actual, err := ImageAndExistingVolumeWithTagRequest.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, ExpectedImageAndExistingVolumeWithTagRequest, actual)
}

func TestCreateSchedulerHints(t *testing.T) {
	opts := servers.SchedulerHintOpts{
		Group: "101aed42-22d9-4a3e-9ba1-21103b0d1aba",
		DifferentHost: []string{
			"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
			"8c19174f-4220-44f0-824a-cd1eeef10287",
		},
		SameHost: []string{
			"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
			"8c19174f-4220-44f0-824a-cd1eeef10287",
		},
		Query:      []any{"=", "$free_ram_mb", "1024"},
		TargetCell: "foobar",
		DifferentCell: []string{
			"bazbar",
			"barbaz",
		},
		BuildNearHostIP:      "192.168.1.1/24",
		AdditionalProperties: map[string]any{"reservation": "a0cf03a5-d921-4877-bb5c-86d26cf818e1"},
	}

	expected := `
		{
			"os:scheduler_hints": {
				"group": "101aed42-22d9-4a3e-9ba1-21103b0d1aba",
				"different_host": [
					"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
					"8c19174f-4220-44f0-824a-cd1eeef10287"
				],
				"same_host": [
					"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
					"8c19174f-4220-44f0-824a-cd1eeef10287"
				],
				"query": "[\"=\",\"$free_ram_mb\",\"1024\"]",
				"target_cell": "foobar",
				"different_cell": [
					"bazbar",
					"barbaz"
				],
				"build_near_host_ip": "192.168.1.1",
				"cidr": "/24",
				"reservation": "a0cf03a5-d921-4877-bb5c-86d26cf818e1"
			}
		}
	`
	actual, err := opts.ToSchedulerHintsMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestCreateComplexSchedulerHints(t *testing.T) {
	opts := servers.SchedulerHintOpts{
		Group: "101aed42-22d9-4a3e-9ba1-21103b0d1aba",
		DifferentHost: []string{
			"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
			"8c19174f-4220-44f0-824a-cd1eeef10287",
		},
		SameHost: []string{
			"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
			"8c19174f-4220-44f0-824a-cd1eeef10287",
		},
		Query:      []any{"and", []string{"=", "$free_ram_mb", "1024"}, []string{"=", "$free_disk_mb", "204800"}},
		TargetCell: "foobar",
		DifferentCell: []string{
			"bazbar",
			"barbaz",
		},
		BuildNearHostIP:      "192.168.1.1/24",
		AdditionalProperties: map[string]any{"reservation": "a0cf03a5-d921-4877-bb5c-86d26cf818e1"},
	}

	expected := `
		{
			"os:scheduler_hints": {
				"group": "101aed42-22d9-4a3e-9ba1-21103b0d1aba",
				"different_host": [
					"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
					"8c19174f-4220-44f0-824a-cd1eeef10287"
				],
				"same_host": [
					"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
					"8c19174f-4220-44f0-824a-cd1eeef10287"
				],
				"query": "[\"and\",[\"=\",\"$free_ram_mb\",\"1024\"],[\"=\",\"$free_disk_mb\",\"204800\"]]",
				"target_cell": "foobar",
				"different_cell": [
					"bazbar",
					"barbaz"
				],
				"build_near_host_ip": "192.168.1.1",
				"cidr": "/24",
				"reservation": "a0cf03a5-d921-4877-bb5c-86d26cf818e1"
			}
		}
	`
	actual, err := opts.ToSchedulerHintsMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestDeleteServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerDeletionSuccessfully(t, fakeServer)

	res := servers.Delete(context.TODO(), client.ServiceClient(fakeServer), "asdfasdfasdf")
	th.AssertNoErr(t, res.Err)
}

func TestForceDeleteServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerForceDeletionSuccessfully(t, fakeServer)

	res := servers.ForceDelete(context.TODO(), client.ServiceClient(fakeServer), "asdfasdfasdf")
	th.AssertNoErr(t, res.Err)
}

func TestGetServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerGetSuccessfully(t, fakeServer)

	client := client.ServiceClient(fakeServer)
	actual, err := servers.Get(context.TODO(), client, "1234asdf").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestGetFaultyServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerGetFaultSuccessfully(t, fakeServer)

	client := client.ServiceClient(fakeServer)
	actual, err := servers.Get(context.TODO(), client, "1234asdf").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	FaultyServer := ServerDerp
	FaultyServer.Fault = DerpFault
	FaultyServer.Locked = nil
	th.CheckDeepEquals(t, FaultyServer, *actual)
}

func TestGetServerWithExtensions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerGetSuccessfully(t, fakeServer)

	var s struct {
		servers.Server
	}

	client := client.ServiceClient(fakeServer)
	err := servers.Get(context.TODO(), client, "1234asdf").ExtractInto(&s)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "nova", s.AvailabilityZone)
	th.AssertEquals(t, "RUNNING", s.PowerState.String())
	th.AssertEquals(t, "", s.TaskState)
	th.AssertEquals(t, "active", s.VmState)
	th.AssertEquals(t, servers.Manual, s.DiskConfig)

	err = servers.Get(context.TODO(), client, "1234asdf").ExtractInto(s)
	if err == nil {
		t.Errorf("Expected error when providing non-pointer struct")
	}
}

func TestUpdateServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerUpdateSuccessfully(t, fakeServer)

	client := client.ServiceClient(fakeServer)
	actual, err := servers.Update(context.TODO(), client, "1234asdf", servers.UpdateOpts{Name: "new-name"}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestChangeServerAdminPassword(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAdminPasswordChangeSuccessfully(t, fakeServer)

	res := servers.ChangeAdminPassword(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", "new-password")
	th.AssertNoErr(t, res.Err)
}

func TestShowConsoleOutput(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleShowConsoleOutputSuccessfully(t, fakeServer, ConsoleOutputBody)

	outputOpts := &servers.ShowConsoleOutputOpts{
		Length: 50,
	}
	actual, err := servers.ShowConsoleOutput(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", outputOpts).Extract()

	th.AssertNoErr(t, err)
	th.AssertByteArrayEquals(t, []byte(ConsoleOutput), []byte(actual))
}

func TestGetPassword(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePasswordGetSuccessfully(t, fakeServer)

	res := servers.GetPassword(context.TODO(), client.ServiceClient(fakeServer), "1234asdf")
	th.AssertNoErr(t, res.Err)
}

func TestRebootServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRebootSuccessfully(t, fakeServer)

	res := servers.Reboot(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", servers.RebootOpts{
		Type: servers.SoftReboot,
	})
	th.AssertNoErr(t, res.Err)
}

func TestRebuildServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleRebuildSuccessfully(t, fakeServer, SingleServerBody)

	opts := servers.RebuildOpts{
		Name:       "new-name",
		AdminPass:  "swordfish",
		ImageRef:   "f90f6034-2570-4974-8351-6b49732ef2eb",
		AccessIPv4: "1.2.3.4",
	}

	actual, err := servers.Rebuild(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", opts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ServerDerp, *actual)
}

func TestRebuildServerWithDiskConfig(t *testing.T) {
	opts := servers.RebuildOpts{
		Name:       "rebuiltserver",
		AdminPass:  "swordfish",
		ImageRef:   "asdfasdfasdf",
		DiskConfig: servers.Auto,
	}
	expected := `
		{
			"rebuild": {
				"name": "rebuiltserver",
				"imageRef": "asdfasdfasdf",
				"adminPass": "swordfish",
				"OS-DCF:diskConfig": "AUTO"
			}
		}
	`

	actual, err := opts.ToServerRebuildMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestResizeServer(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/servers/1234asdf/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{ "resize": { "flavorRef": "2" } }`)

		w.WriteHeader(http.StatusAccepted)
	})

	res := servers.Resize(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", servers.ResizeOpts{FlavorRef: "2"})
	th.AssertNoErr(t, res.Err)
}

func TestResizeServerWithDiskConfig(t *testing.T) {
	opts := servers.ResizeOpts{
		FlavorRef:  "performance1-8",
		DiskConfig: servers.Auto,
	}
	expected := `
		{
			"resize": {
				"flavorRef": "performance1-8",
				"OS-DCF:diskConfig": "AUTO"
			}
		}
	`

	actual, err := opts.ToServerResizeMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestConfirmResize(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/servers/1234asdf/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{ "confirmResize": null }`)

		w.WriteHeader(http.StatusNoContent)
	})

	res := servers.ConfirmResize(context.TODO(), client.ServiceClient(fakeServer), "1234asdf")
	th.AssertNoErr(t, res.Err)
}

func TestRevertResize(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/servers/1234asdf/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{ "revertResize": null }`)

		w.WriteHeader(http.StatusAccepted)
	})

	res := servers.RevertResize(context.TODO(), client.ServiceClient(fakeServer), "1234asdf")
	th.AssertNoErr(t, res.Err)
}

func TestGetMetadatum(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleMetadatumGetSuccessfully(t, fakeServer)

	expected := map[string]string{"foo": "bar"}
	actual, err := servers.Metadatum(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", "foo").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestCreateMetadatum(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleMetadatumCreateSuccessfully(t, fakeServer)

	expected := map[string]string{"foo": "bar"}
	actual, err := servers.CreateMetadatum(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", servers.MetadatumOpts{"foo": "bar"}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteMetadatum(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleMetadatumDeleteSuccessfully(t, fakeServer)

	err := servers.DeleteMetadatum(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", "foo").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetMetadata(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleMetadataGetSuccessfully(t, fakeServer)

	expected := map[string]string{"foo": "bar", "this": "that"}
	actual, err := servers.Metadata(context.TODO(), client.ServiceClient(fakeServer), "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestResetMetadata(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleMetadataResetSuccessfully(t, fakeServer)

	expected := map[string]string{"foo": "bar", "this": "that"}
	actual, err := servers.ResetMetadata(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", servers.MetadataOpts{
		"foo":  "bar",
		"this": "that",
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestUpdateMetadata(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleMetadataUpdateSuccessfully(t, fakeServer)

	expected := map[string]string{"foo": "baz", "this": "those"}
	actual, err := servers.UpdateMetadata(context.TODO(), client.ServiceClient(fakeServer), "1234asdf", servers.MetadataOpts{
		"foo":  "baz",
		"this": "those",
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestListAddresses(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAddressListSuccessfully(t, fakeServer)

	expected := ListAddressesExpected
	pages := 0
	err := servers.ListAddresses(client.ServiceClient(fakeServer), "asdfasdfasdf").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := servers.ExtractAddresses(page)
		th.AssertNoErr(t, err)

		if len(actual) != 2 {
			t.Fatalf("Expected 2 networks, got %d", len(actual))
		}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, pages)
}

func TestListAddressesByNetwork(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNetworkAddressListSuccessfully(t, fakeServer)

	expected := ListNetworkAddressesExpected
	pages := 0
	err := servers.ListAddressesByNetwork(client.ServiceClient(fakeServer), "asdfasdfasdf", "public").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := servers.ExtractNetworkAddresses(page)
		th.AssertNoErr(t, err)

		if len(actual) != 2 {
			t.Fatalf("Expected 2 addresses, got %d", len(actual))
		}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, pages)
}

func TestCreateServerImageBeforeMicroversion_2_45(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	expected := HandleCreateServerImageSuccessfullyBeforeMicroversion_2_45(t, fakeServer)

	imageID, err := servers.CreateImage(context.TODO(), client.ServiceClient(fakeServer), "serverimage", servers.CreateImageOpts{Name: "test"}).ExtractImageID()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected, imageID)
}

func TestCreateServerImageSinceMicroversion_2_45(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	expected := HandleCreateServerImageSuccessfullySinceMicroversion_2_45(t, fakeServer)

	imageID, err := servers.CreateImage(context.TODO(), client.ServiceClient(fakeServer), "serverimage", servers.CreateImageOpts{Name: "test"}).ExtractImageID()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected, imageID)
}

func TestMarshalPersonality(t *testing.T) {
	name := "/etc/test"
	contents := []byte("asdfasdf")

	personality := servers.Personality{
		&servers.File{
			Path:     name,
			Contents: contents,
		},
	}

	data, err := json.Marshal(personality)
	if err != nil {
		t.Fatal(err)
	}

	var actual []map[string]string
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Fatal(err)
	}

	if len(actual) != 1 {
		t.Fatal("expected personality length 1")
	}

	if actual[0]["path"] != name {
		t.Fatal("file path incorrect")
	}

	if actual[0]["contents"] != base64.StdEncoding.EncodeToString(contents) {
		t.Fatal("file contents incorrect")
	}
}

func TestCreateServerWithTags(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerWithTagsCreationSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	c.Microversion = "2.52"

	tags := []string{"foo", "bar"}
	ServerDerpTags := ServerDerp
	ServerDerpTags.Tags = &tags
	ServerDerpTags.Locked = nil

	createOpts := servers.CreateOpts{
		Name:      "derp",
		ImageRef:  "f90f6034-2570-4974-8351-6b49732ef2eb",
		FlavorRef: "1",
		Tags:      tags,
	}
	res := servers.Create(context.TODO(), c, createOpts, nil)
	th.AssertNoErr(t, res.Err)
	actualServer, err := res.Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ServerDerpTags, *actualServer)
}

func TestCreateServerWithHypervisorHostname(t *testing.T) {
	opts := servers.CreateOpts{
		Name:               "createdserver",
		FlavorRef:          "performance1-1",
		ImageRef:           "asdfasdfasdf",
		HypervisorHostname: "test-hypervisor",
	}
	expected := `
    {
        "server": {
            "name":"createdserver",
            "flavorRef":"performance1-1",
            "imageRef":"asdfasdfasdf",
            "hypervisor_hostname":"test-hypervisor"
        }
    }
    `
	actual, err := opts.ToServerCreateMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestUpdateServerHostname(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleServerHostnameUpdateSuccessfully(t, fakeServer)

	client := client.ServiceClient(fakeServer)
	newHostname := "new-hostname"
	actual, err := servers.Update(context.TODO(), client, "1234asdf", servers.UpdateOpts{Hostname: &newHostname}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, ServerDerp, *actual)
}
