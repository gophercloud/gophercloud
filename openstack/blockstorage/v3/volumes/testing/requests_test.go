package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListWithExtensions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	count := 0

	err := volumes.List(client.ServiceClient(), &volumes.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := volumes.ExtractVolumes(page)
		if err != nil {
			t.Errorf("Failed to extract volumes: %v", err)
			return false, err
		}

		expected := []volumes.Volume{
			{
				ID:   "289da7f8-6440-407c-9fb4-7db01ec49164",
				Name: "vol-001",
				Attachments: []volumes.Attachment{{
					ServerID:     "83ec2e3b-4321-422b-8706-a84185f52a0a",
					AttachmentID: "05551600-a936-4d4a-ba42-79a037c1-c91a",
					AttachedAt:   time.Date(2016, 8, 6, 14, 48, 20, 0, time.UTC),
					HostName:     "foobar",
					VolumeID:     "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
					Device:       "/dev/vdc",
					ID:           "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
				}},
				AvailabilityZone:   "nova",
				Bootable:           "false",
				ConsistencyGroupID: "",
				CreatedAt:          time.Date(2015, 9, 17, 3, 35, 3, 0, time.UTC),
				Description:        "",
				Encrypted:          false,
				Host:               "host-001",
				Metadata:           map[string]string{"foo": "bar"},
				Multiattach:        false,
				TenantID:           "304dc00909ac4d0da6c62d816bcb3459",
				//ReplicationDriverData:     "",
				//ReplicationExtendedStatus: "",
				ReplicationStatus: "disabled",
				Size:              75,
				SnapshotID:        "",
				SourceVolID:       "",
				Status:            "available",
				UserID:            "ff1ce52c03ab433aaba9108c2e3ef541",
				VolumeType:        "lvmdriver-1",
			},
			{
				ID:                 "96c3bda7-c82a-4f50-be73-ca7621794835",
				Name:               "vol-002",
				Attachments:        []volumes.Attachment{},
				AvailabilityZone:   "nova",
				Bootable:           "false",
				ConsistencyGroupID: "",
				CreatedAt:          time.Date(2015, 9, 17, 3, 32, 29, 0, time.UTC),
				Description:        "",
				Encrypted:          false,
				Metadata:           map[string]string{},
				Multiattach:        false,
				TenantID:           "304dc00909ac4d0da6c62d816bcb3459",
				//ReplicationDriverData:     "",
				//ReplicationExtendedStatus: "",
				ReplicationStatus: "disabled",
				Size:              75,
				SnapshotID:        "",
				SourceVolID:       "",
				Status:            "available",
				UserID:            "ff1ce52c03ab433aaba9108c2e3ef541",
				VolumeType:        "lvmdriver-1",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListAllWithExtensions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allPages, err := volumes.List(client.ServiceClient(), &volumes.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	var actual []volumes.Volume
	err = volumes.ExtractVolumesInto(allPages, &actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(actual))
	th.AssertEquals(t, "host-001", actual[0].Host)
	th.AssertEquals(t, "", actual[1].Host)
	th.AssertEquals(t, "304dc00909ac4d0da6c62d816bcb3459", actual[0].TenantID)
}

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allPages, err := volumes.List(client.ServiceClient(), &volumes.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := volumes.ExtractVolumes(allPages)
	th.AssertNoErr(t, err)

	expected := []volumes.Volume{
		{
			ID:   "289da7f8-6440-407c-9fb4-7db01ec49164",
			Name: "vol-001",
			Attachments: []volumes.Attachment{{
				ServerID:     "83ec2e3b-4321-422b-8706-a84185f52a0a",
				AttachmentID: "05551600-a936-4d4a-ba42-79a037c1-c91a",
				AttachedAt:   time.Date(2016, 8, 6, 14, 48, 20, 0, time.UTC),
				HostName:     "foobar",
				VolumeID:     "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
				Device:       "/dev/vdc",
				ID:           "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
			}},
			AvailabilityZone:   "nova",
			Bootable:           "false",
			ConsistencyGroupID: "",
			CreatedAt:          time.Date(2015, 9, 17, 3, 35, 3, 0, time.UTC),
			Description:        "",
			Encrypted:          false,
			Host:               "host-001",
			Metadata:           map[string]string{"foo": "bar"},
			Multiattach:        false,
			TenantID:           "304dc00909ac4d0da6c62d816bcb3459",
			//ReplicationDriverData:     "",
			//ReplicationExtendedStatus: "",
			ReplicationStatus: "disabled",
			Size:              75,
			SnapshotID:        "",
			SourceVolID:       "",
			Status:            "available",
			UserID:            "ff1ce52c03ab433aaba9108c2e3ef541",
			VolumeType:        "lvmdriver-1",
		},
		{
			ID:                 "96c3bda7-c82a-4f50-be73-ca7621794835",
			Name:               "vol-002",
			Attachments:        []volumes.Attachment{},
			AvailabilityZone:   "nova",
			Bootable:           "false",
			ConsistencyGroupID: "",
			CreatedAt:          time.Date(2015, 9, 17, 3, 32, 29, 0, time.UTC),
			Description:        "",
			Encrypted:          false,
			Metadata:           map[string]string{},
			Multiattach:        false,
			TenantID:           "304dc00909ac4d0da6c62d816bcb3459",
			//ReplicationDriverData:     "",
			//ReplicationExtendedStatus: "",
			ReplicationStatus: "disabled",
			Size:              75,
			SnapshotID:        "",
			SourceVolID:       "",
			Status:            "available",
			UserID:            "ff1ce52c03ab433aaba9108c2e3ef541",
			VolumeType:        "lvmdriver-1",
		},
	}

	th.CheckDeepEquals(t, expected, actual)

}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	v, err := volumes.Get(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, v.Name, "vol-001")
	th.AssertEquals(t, v.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &volumes.CreateOpts{Size: 75, Name: "vol-001"}
	n, err := volumes.Create(context.TODO(), client.ServiceClient(), options, nil).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Size, 75)
	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestCreateSchedulerHints(t *testing.T) {
	base := volumes.SchedulerHintOpts{
		DifferentHost: []string{
			"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
			"8c19174f-4220-44f0-824a-cd1eeef10287",
		},
		SameHost: []string{
			"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
			"8c19174f-4220-44f0-824a-cd1eeef10287",
		},
		LocalToInstance:      "0ffb2c1b-d621-4fc1-9ae4-88d99c088ff6",
		AdditionalProperties: map[string]any{"mark": "a0cf03a5-d921-4877-bb5c-86d26cf818e1"},
	}
	expected := `
		{
			"OS-SCH-HNT:scheduler_hints": {
				"different_host": [
					"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
					"8c19174f-4220-44f0-824a-cd1eeef10287"
				],
				"same_host": [
					"a0cf03a5-d921-4877-bb5c-86d26cf818e1",
					"8c19174f-4220-44f0-824a-cd1eeef10287"
				],
				"local_to_instance": "0ffb2c1b-d621-4fc1-9ae4-88d99c088ff6",
				"mark": "a0cf03a5-d921-4877-bb5c-86d26cf818e1"
			}
		}
	`
	actual, err := base.ToSchedulerHintsMap()
	th.AssertNoErr(t, err)
	th.CheckJSONEquals(t, expected, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := volumes.Delete(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", volumes.DeleteOpts{})
	th.AssertNoErr(t, res.Err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateResponse(t)

	var name = "vol-002"
	options := volumes.UpdateOpts{Name: &name}
	v, err := volumes.Update(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", options).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "vol-002", v.Name)
}

func TestGetWithExtensions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	var v volumes.Volume
	err := volumes.Get(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").ExtractInto(&v)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "304dc00909ac4d0da6c62d816bcb3459", v.TenantID)
	th.AssertEquals(t, "centos", v.VolumeImageMetadata["image_name"])

	err = volumes.Get(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").ExtractInto(v)
	if err == nil {
		t.Errorf("Expected error when providing non-pointer struct")
	}
}

func TestCreateFromBackup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateVolumeFromBackupResponse(t)

	options := volumes.CreateOpts{
		Name:     "vol-001",
		BackupID: "20c792f0-bb03-434f-b653-06ef238e337e",
	}
	v, err := volumes.Create(context.TODO(), client.ServiceClient(), options, nil).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, v.Size, 30)
	th.AssertEquals(t, v.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertEquals(t, *v.BackupID, "20c792f0-bb03-434f-b653-06ef238e337e")
}

func TestAttach(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockAttachResponse(t)

	options := &volumes.AttachOpts{
		MountPoint:   "/mnt",
		Mode:         "rw",
		InstanceUUID: "50902f4f-a974-46a0-85e9-7efc5e22dfdd",
	}
	err := volumes.Attach(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestBeginDetaching(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockBeginDetachingResponse(t)

	err := volumes.BeginDetaching(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDetach(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDetachResponse(t)

	err := volumes.Detach(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", &volumes.DetachOpts{}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUploadImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	MockUploadImageResponse(t)
	options := &volumes.UploadImageOpts{
		ContainerFormat: "bare",
		DiskFormat:      "raw",
		ImageName:       "test",
		Force:           true,
	}

	actual, err := volumes.UploadImage(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).Extract()
	th.AssertNoErr(t, err)

	expected := volumes.VolumeImage{
		VolumeID:        "cd281d77-8217-4830-be95-9528227c105c",
		ContainerFormat: "bare",
		DiskFormat:      "raw",
		Description:     "",
		ImageID:         "ecb92d98-de08-45db-8235-bbafe317269c",
		ImageName:       "test",
		Size:            5,
		Status:          "uploading",
		UpdatedAt:       time.Date(2017, 7, 17, 9, 29, 22, 0, time.UTC),
		VolumeType: volumes.ImageVolumeType{
			ID:          "b7133444-62f6-4433-8da3-70ac332229b7",
			Name:        "basic.ru-2a",
			Description: "",
			IsPublic:    true,
			ExtraSpecs:  map[string]any{"volume_backend_name": "basic.ru-2a"},
			QosSpecsID:  "",
			Deleted:     false,
			DeletedAt:   time.Time{},
			CreatedAt:   time.Date(2016, 5, 4, 8, 54, 14, 0, time.UTC),
			UpdatedAt:   time.Date(2016, 5, 4, 9, 15, 33, 0, time.UTC),
		},
	}
	th.AssertDeepEquals(t, expected, actual)
}

func TestReserve(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockReserveResponse(t)

	err := volumes.Reserve(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUnreserve(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUnreserveResponse(t)

	err := volumes.Unreserve(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestInitializeConnection(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockInitializeConnectionResponse(t)

	options := &volumes.InitializeConnectionOpts{
		IP:        "127.0.0.1",
		Host:      "stack",
		Initiator: "iqn.1994-05.com.redhat:17cf566367d2",
		Multipath: gophercloud.Disabled,
		Platform:  "x86_64",
		OSType:    "linux2",
	}
	_, err := volumes.InitializeConnection(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).Extract()
	th.AssertNoErr(t, err)
}

func TestTerminateConnection(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockTerminateConnectionResponse(t)

	options := &volumes.TerminateConnectionOpts{
		IP:        "127.0.0.1",
		Host:      "stack",
		Initiator: "iqn.1994-05.com.redhat:17cf566367d2",
		Multipath: gophercloud.Enabled,
		Platform:  "x86_64",
		OSType:    "linux2",
	}
	err := volumes.TerminateConnection(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestExtendSize(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockExtendSizeResponse(t)

	options := &volumes.ExtendSizeOpts{
		NewSize: 3,
	}

	err := volumes.ExtendSize(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestForceDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockForceDeleteResponse(t)

	res := volumes.ForceDelete(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}

func TestSetImageMetadata(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockSetImageMetadataResponse(t)

	options := &volumes.ImageMetadataOpts{
		Metadata: map[string]string{
			"label": "test",
		},
	}

	err := volumes.SetImageMetadata(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestSetBootable(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockSetBootableResponse(t)

	options := volumes.BootableOpts{
		Bootable: true,
	}

	err := volumes.SetBootable(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestReImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockReImageResponse(t)

	options := volumes.ReImageOpts{
		ImageID:         "71543ced-a8af-45b6-a5c4-a46282108a90",
		ReImageReserved: false,
	}

	err := volumes.ReImage(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestChangeType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockChangeTypeResponse(t)

	options := &volumes.ChangeTypeOpts{
		NewType:         "ssd",
		MigrationPolicy: "on-demand",
	}

	err := volumes.ChangeType(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestResetStatus(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockResetStatusResponse(t)

	options := &volumes.ResetStatusOpts{
		Status:          "error",
		AttachStatus:    "detached",
		MigrationStatus: "migrating",
	}

	err := volumes.ResetStatus(context.TODO(), client.ServiceClient(), "cd281d77-8217-4830-be95-9528227c105c", options).ExtractErr()
	th.AssertNoErr(t, err)
}
