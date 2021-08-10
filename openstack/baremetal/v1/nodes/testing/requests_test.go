package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListDetailNodes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeListDetailSuccessfully(t)

	pages := 0
	err := nodes.ListDetail(client.ServiceClient(), nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := nodes.ExtractNodes(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 3 {
			t.Fatalf("Expected 3 nodes, got %d", len(actual))
		}
		th.CheckDeepEquals(t, NodeFoo, actual[0])
		th.CheckDeepEquals(t, NodeBar, actual[1])
		th.CheckDeepEquals(t, NodeBaz, actual[2])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListNodes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeListSuccessfully(t)

	pages := 0
	err := nodes.List(client.ServiceClient(), nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := nodes.ExtractNodes(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 3 {
			t.Fatalf("Expected 3 nodes, got %d", len(actual))
		}
		th.AssertEquals(t, "foo", actual[0].Name)
		th.AssertEquals(t, "bar", actual[1].Name)
		th.AssertEquals(t, "baz", actual[2].Name)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListOpts(t *testing.T) {
	// Detail cannot take Fields
	opts := nodes.ListOpts{
		Fields: []string{"name", "uuid"},
	}

	_, err := opts.ToNodeListDetailQuery()
	th.AssertEquals(t, err.Error(), "fields is not a valid option when getting a detailed listing of nodes")

	// Regular ListOpts can
	query, err := opts.ToNodeListQuery()
	th.AssertEquals(t, query, "?fields=name&fields=uuid")
	th.AssertNoErr(t, err)
}

func TestCreateNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeCreationSuccessfully(t, SingleNodeBody)

	actual, err := nodes.Create(client.ServiceClient(), nodes.CreateOpts{
		Name:          "foo",
		Driver:        "ipmi",
		BootInterface: "pxe",
		DriverInfo: map[string]interface{}{
			"ipmi_port":      "6230",
			"ipmi_username":  "admin",
			"deploy_kernel":  "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
			"ipmi_address":   "192.168.122.1",
			"deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
			"ipmi_password":  "admin",
		},
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, NodeFoo, *actual)
}

func TestDeleteNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeDeletionSuccessfully(t)

	res := nodes.Delete(client.ServiceClient(), "asdfasdfasdf")
	th.AssertNoErr(t, res.Err)
}

func TestGetNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeGetSuccessfully(t)

	c := client.ServiceClient()
	actual, err := nodes.Get(c, "1234asdf").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, NodeFoo, *actual)
}

func TestUpdateNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeUpdateSuccessfully(t, SingleNodeBody)

	c := client.ServiceClient()
	actual, err := nodes.Update(c, "1234asdf", nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:   nodes.ReplaceOp,
			Path: "/properties",
			Value: map[string]interface{}{
				"root_gb": 25,
			},
		},
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, NodeFoo, *actual)
}

func TestUpdateRequiredOp(t *testing.T) {
	c := client.ServiceClient()
	_, err := nodes.Update(c, "1234asdf", nodes.UpdateOpts{
		nodes.UpdateOperation{
			Path:  "/driver",
			Value: "new-driver",
		},
	}).Extract()

	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}

}

func TestUpdateRequiredPath(t *testing.T) {
	c := client.ServiceClient()
	_, err := nodes.Update(c, "1234asdf", nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:    nodes.ReplaceOp,
			Value: "new-driver",
		},
	}).Extract()

	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

func TestValidateNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeValidateSuccessfully(t)

	c := client.ServiceClient()
	actual, err := nodes.Validate(c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeFooValidation, *actual)
}

func TestInjectNMI(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInjectNMISuccessfully(t)

	c := client.ServiceClient()
	err := nodes.InjectNMI(c, "1234asdf").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestSetBootDevice(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSetBootDeviceSuccessfully(t)

	c := client.ServiceClient()
	err := nodes.SetBootDevice(c, "1234asdf", nodes.BootDeviceOpts{
		BootDevice: "pxe",
		Persistent: false,
	}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetBootDevice(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetBootDeviceSuccessfully(t)

	c := client.ServiceClient()
	bootDevice, err := nodes.GetBootDevice(c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeBootDevice, *bootDevice)
}

func TestGetSupportedBootDevices(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSupportedBootDeviceSuccessfully(t)

	c := client.ServiceClient()
	bootDevices, err := nodes.GetSupportedBootDevices(c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeSupportedBootDevice, bootDevices)
}

func TestNodeChangeProvisionStateActive(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeChangeProvisionStateActive(t)

	c := client.ServiceClient()
	err := nodes.ChangeProvisionState(c, "1234asdf", nodes.ProvisionStateOpts{
		Target:      nodes.TargetActive,
		ConfigDrive: "http://127.0.0.1/images/test-node-config-drive.iso.gz",
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestNodeChangeProvisionStateActiveWithSteps(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeChangeProvisionStateActiveWithSteps(t)

	c := client.ServiceClient()
	err := nodes.ChangeProvisionState(c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetActive,
		DeploySteps: []nodes.DeployStep{
			{
				Interface: nodes.InterfaceDeploy,
				Step:      "inject_files",
				Priority:  50,
				Args: map[string]interface{}{
					"files": []interface{}{},
				},
			},
		},
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestHandleNodeChangeProvisionStateConfigDrive(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleNodeChangeProvisionStateConfigDrive(t)

	c := client.ServiceClient()

	err := nodes.ChangeProvisionState(c, "1234asdf", nodes.ProvisionStateOpts{
		Target:      nodes.TargetActive,
		ConfigDrive: ConfigDriveMap,
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestNodeChangeProvisionStateClean(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeChangeProvisionStateClean(t)

	c := client.ServiceClient()
	err := nodes.ChangeProvisionState(c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Interface: nodes.InterfaceDeploy,
				Step:      "upgrade_firmware",
				Args: map[string]interface{}{
					"force": "True",
				},
			},
		},
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestNodeChangeProvisionStateCleanWithConflict(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleNodeChangeProvisionStateCleanWithConflict(t)

	c := client.ServiceClient()
	err := nodes.ChangeProvisionState(c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Interface: nodes.InterfaceDeploy,
				Step:      "upgrade_firmware",
				Args: map[string]interface{}{
					"force": "True",
				},
			},
		},
	}).ExtractErr()

	if _, ok := err.(gophercloud.ErrDefault409); !ok {
		t.Fatal("ErrDefault409 was expected to occur")
	}
}

func TestCleanStepRequiresInterface(t *testing.T) {
	c := client.ServiceClient()
	err := nodes.ChangeProvisionState(c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Step: "upgrade_firmware",
				Args: map[string]interface{}{
					"force": "True",
				},
			},
		},
	}).ExtractErr()

	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

func TestCleanStepRequiresStep(t *testing.T) {
	c := client.ServiceClient()
	err := nodes.ChangeProvisionState(c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Interface: nodes.InterfaceDeploy,
				Args: map[string]interface{}{
					"force": "True",
				},
			},
		},
	}).ExtractErr()

	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

func TestChangePowerState(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleChangePowerStateSuccessfully(t)

	opts := nodes.PowerStateOpts{
		Target:  nodes.PowerOn,
		Timeout: 100,
	}

	c := client.ServiceClient()
	err := nodes.ChangePowerState(c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestChangePowerStateWithConflict(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleChangePowerStateWithConflict(t)

	opts := nodes.PowerStateOpts{
		Target:  nodes.PowerOn,
		Timeout: 100,
	}

	c := client.ServiceClient()
	err := nodes.ChangePowerState(c, "1234asdf", opts).ExtractErr()
	if _, ok := err.(gophercloud.ErrDefault409); !ok {
		t.Fatal("ErrDefault409 was expected to occur")
	}
}

func TestSetRAIDConfig(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSetRAIDConfig(t)

	sizeGB := 100
	isRootVolume := true

	config := nodes.RAIDConfigOpts{
		LogicalDisks: []nodes.LogicalDisk{
			{
				SizeGB:       &sizeGB,
				IsRootVolume: &isRootVolume,
				RAIDLevel:    nodes.RAID1,
			},
		},
	}

	c := client.ServiceClient()
	err := nodes.SetRAIDConfig(c, "1234asdf", config).ExtractErr()
	th.AssertNoErr(t, err)
}

// Without specifying a size, we need to send a string: "MAX"
func TestSetRAIDConfigMaxSize(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSetRAIDConfigMaxSize(t)

	isRootVolume := true

	config := nodes.RAIDConfigOpts{
		LogicalDisks: []nodes.LogicalDisk{
			{
				IsRootVolume: &isRootVolume,
				RAIDLevel:    nodes.RAID1,
			},
		},
	}

	c := client.ServiceClient()
	err := nodes.SetRAIDConfig(c, "1234asdf", config).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestToRAIDConfigMap(t *testing.T) {
	cases := []struct {
		name     string
		opts     nodes.RAIDConfigOpts
		expected map[string]interface{}
	}{
		{
			name: "LogicalDisks is empty",
			opts: nodes.RAIDConfigOpts{},
			expected: map[string]interface{}{
				"logical_disks": nil,
			},
		},
		{
			name: "LogicalDisks is nil",
			opts: nodes.RAIDConfigOpts{
				LogicalDisks: nil,
			},
			expected: map[string]interface{}{
				"logical_disks": nil,
			},
		},
		{
			name: "PhysicalDisks is []string",
			opts: nodes.RAIDConfigOpts{
				LogicalDisks: []nodes.LogicalDisk{
					{
						RAIDLevel:     "0",
						VolumeName:    "root",
						PhysicalDisks: []interface{}{"6I:1:5", "6I:1:6", "6I:1:7"},
					},
				},
			},
			expected: map[string]interface{}{
				"logical_disks": []map[string]interface{}{
					{
						"raid_level":     "0",
						"size_gb":        "MAX",
						"volume_name":    "root",
						"physical_disks": []interface{}{"6I:1:5", "6I:1:6", "6I:1:7"},
					},
				},
			},
		},
		{
			name: "PhysicalDisks is []map[string]string",
			opts: nodes.RAIDConfigOpts{
				LogicalDisks: []nodes.LogicalDisk{
					{
						RAIDLevel:  "0",
						VolumeName: "root",
						Controller: "software",
						PhysicalDisks: []interface{}{
							map[string]string{
								"size": "> 100",
							},
							map[string]string{
								"size": "> 100",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"logical_disks": []map[string]interface{}{
					{
						"raid_level":  "0",
						"size_gb":     "MAX",
						"volume_name": "root",
						"controller":  "software",
						"physical_disks": []interface{}{
							map[string]interface{}{
								"size": "> 100",
							},
							map[string]interface{}{
								"size": "> 100",
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := c.opts.ToRAIDConfigMap()
			th.CheckDeepEquals(t, c.expected, got)
		})
	}
}

func TestListBIOSSettings(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListBIOSSettingsSuccessfully(t)

	c := client.ServiceClient()
	actual, err := nodes.ListBIOSSettings(c, "1234asdf", nil).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeBIOSSettings, actual)
}

func TestListDetailBIOSSettings(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListDetailBIOSSettingsSuccessfully(t)

	opts := nodes.ListBIOSSettingsOpts{
		Detail: true,
	}

	c := client.ServiceClient()
	actual, err := nodes.ListBIOSSettings(c, "1234asdf", opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeDetailBIOSSettings, actual)
}

func TestGetBIOSSetting(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetBIOSSettingSuccessfully(t)

	c := client.ServiceClient()
	actual, err := nodes.GetBIOSSetting(c, "1234asdf", "ProcVirtualization").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeSingleBIOSSetting, *actual)
}

func TestListBIOSSettingsOpts(t *testing.T) {
	// Detail cannot take Fields
	opts := nodes.ListBIOSSettingsOpts{
		Detail: true,
		Fields: []string{"name", "value"},
	}

	_, err := opts.ToListBIOSSettingsOptsQuery()
	th.AssertEquals(t, err.Error(), "cannot have both fields and detail options for BIOS settings")
}

func TestGetVendorPassthruMethods(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetVendorPassthruMethodsSuccessfully(t)

	c := client.ServiceClient()
	actual, err := nodes.GetVendorPassthruMethods(c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeVendorPassthruMethods, *actual)
}

func TestGetAllSubscriptions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetAllSubscriptionsVendorPassthruSuccessfully(t)

	c := client.ServiceClient()
	method := nodes.CallVendorPassthruOpts{
		Method: "get_all_subscriptions",
	}
	actual, err := nodes.GetAllSubscriptions(c, "1234asdf", method).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeGetAllSubscriptions, *actual)
}

func TestGetSubscription(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSubscriptionVendorPassthruSuccessfully(t)

	c := client.ServiceClient()
	method := nodes.CallVendorPassthruOpts{
		Method: "get_subscription",
	}
	subscriptionOpt := nodes.GetSubscriptionOpts{
		Id: "62dbd1b6-f637-11eb-b551-4cd98f20754c",
	}
	actual, err := nodes.GetSubscription(c, "1234asdf", method, subscriptionOpt).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeGetSubscription, *actual)
}

func TestCreateSubscriptionAllParameters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSubscriptionVendorPassthruAllParametersSuccessfully(t)

	c := client.ServiceClient()
	method := nodes.CallVendorPassthruOpts{
		Method: "create_subscription",
	}
	createOpt := nodes.CreateSubscriptionOpts{
		Destination: "https://someurl",
		Context:     "gophercloud",
		Protocol:    "Redfish",
		EventTypes:  []string{"Alert"},
	}
	actual, err := nodes.CreateSubscription(c, "1234asdf", method, createOpt).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeCreateSubscriptionAllParameters, *actual)
}

func TestCreateSubscriptionWithRequiredParameters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSubscriptionVendorPassthruRequiredParametersSuccessfully(t)

	c := client.ServiceClient()
	method := nodes.CallVendorPassthruOpts{
		Method: "create_subscription",
	}
	createOpt := nodes.CreateSubscriptionOpts{
		Destination: "https://somedestinationurl",
	}
	actual, err := nodes.CreateSubscription(c, "1234asdf", method, createOpt).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeCreateSubscriptionRequiredParameters, *actual)
}

func TestDeleteSubscription(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSubscriptionVendorPassthruSuccessfully(t)

	c := client.ServiceClient()
	method := nodes.CallVendorPassthruOpts{
		Method: "delete_subscription",
	}
	deleteOpt := nodes.DeleteSubscriptionOpts{
		Id: "344a3e2-978a-444e-990a-cbf47c62ef88",
	}
	err := nodes.DeleteSubscription(c, "1234asdf", method, deleteOpt).ExtractErr()
	th.AssertNoErr(t, err)
}
