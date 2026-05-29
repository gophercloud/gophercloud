package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListDetailNodes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeListDetailSuccessfully(t, fakeServer)

	pages := 0
	err := nodes.ListDetail(client.ServiceClient(fakeServer), nodes.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeListSuccessfully(t, fakeServer)

	pages := 0
	err := nodes.List(client.ServiceClient(fakeServer), nodes.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	th.AssertEquals(t, "?fields=name%2Cuuid", query)
	th.AssertNoErr(t, err)
}

func TestCreateNode(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeCreationSuccessfully(t, fakeServer, SingleNodeBody)

	actual, err := nodes.Create(context.TODO(), client.ServiceClient(fakeServer), nodes.CreateOpts{
		Name:          "foo",
		Driver:        "ipmi",
		BootInterface: "pxe",
		DriverInfo: map[string]any{
			"ipmi_port":      "6230",
			"ipmi_username":  "admin",
			"deploy_kernel":  "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
			"ipmi_address":   "192.168.122.1",
			"deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
			"ipmi_password":  "admin",
		},
		FirmwareInterface: "no-firmware",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, NodeFoo, *actual)
}

func TestDeleteNode(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeDeletionSuccessfully(t, fakeServer)

	res := nodes.Delete(context.TODO(), client.ServiceClient(fakeServer), "asdfasdfasdf")
	th.AssertNoErr(t, res.Err)
}

func TestGetNode(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeGetSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.Get(context.TODO(), c, "1234asdf").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, NodeFoo, *actual)
}

func TestUpdateNode(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeUpdateSuccessfully(t, fakeServer, SingleNodeBody)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.Update(context.TODO(), c, "1234asdf", nodes.UpdateOpts{
		nodes.UpdateOperation{
			Op:   nodes.ReplaceOp,
			Path: "/properties",
			Value: map[string]any{
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	c := client.ServiceClient(fakeServer)
	_, err := nodes.Update(context.TODO(), c, "1234asdf", nodes.UpdateOpts{
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	c := client.ServiceClient(fakeServer)
	_, err := nodes.Update(context.TODO(), c, "1234asdf", nodes.UpdateOpts{
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeValidateSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.Validate(context.TODO(), c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeFooValidation, *actual)
}

func TestInjectNMI(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleInjectNMISuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.InjectNMI(context.TODO(), c, "1234asdf").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestSetBootDevice(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSetBootDeviceSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.SetBootDevice(context.TODO(), c, "1234asdf", nodes.BootDeviceOpts{
		BootDevice: "pxe",
		Persistent: false,
	}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetBootDevice(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetBootDeviceSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	bootDevice, err := nodes.GetBootDevice(context.TODO(), c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeBootDevice, *bootDevice)
}

func TestGetSupportedBootDevices(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSupportedBootDeviceSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	bootDevices, err := nodes.GetSupportedBootDevices(context.TODO(), c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeSupportedBootDevice, bootDevices)
}

func TestNodeChangeProvisionStateActive(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeChangeProvisionStateActive(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target:      nodes.TargetActive,
		ConfigDrive: "http://127.0.0.1/images/test-node-config-drive.iso.gz",
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestNodeChangeProvisionStateActiveWithSteps(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeChangeProvisionStateActiveWithSteps(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetActive,
		DeploySteps: []nodes.DeployStep{
			{
				Interface: nodes.InterfaceDeploy,
				Step:      "inject_files",
				Priority:  50,
				Args: map[string]any{
					"files": []any{},
				},
			},
		},
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestHandleNodeChangeProvisionStateConfigDrive(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeChangeProvisionStateConfigDrive(t, fakeServer)

	c := client.ServiceClient(fakeServer)

	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target:      nodes.TargetActive,
		ConfigDrive: ConfigDriveMap,
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestNodeChangeProvisionStateClean(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeChangeProvisionStateClean(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Interface: nodes.InterfaceDeploy,
				Step:      "upgrade_firmware",
				Args: map[string]any{
					"force": "True",
				},
			},
		},
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestNodeChangeProvisionStateCleanWithConflict(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeChangeProvisionStateCleanWithConflict(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Interface: nodes.InterfaceDeploy,
				Step:      "upgrade_firmware",
				Args: map[string]any{
					"force": "True",
				},
			},
		},
	}).ExtractErr()

	if !gophercloud.ResponseCodeIs(err, http.StatusConflict) {
		t.Fatalf("expected 409 response, but got %s", err.Error())
	}
}

func TestCleanStepRequiresInterface(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Step: "upgrade_firmware",
				Args: map[string]any{
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetClean,
		CleanSteps: []nodes.CleanStep{
			{
				Interface: nodes.InterfaceDeploy,
				Args: map[string]any{
					"force": "True",
				},
			},
		},
	}).ExtractErr()

	if _, ok := err.(gophercloud.ErrMissingInput); !ok {
		t.Fatal("ErrMissingInput was expected to occur")
	}
}

func TestNodeChangeProvisionStateService(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleNodeChangeProvisionStateService(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangeProvisionState(context.TODO(), c, "1234asdf", nodes.ProvisionStateOpts{
		Target: nodes.TargetService,
		ServiceSteps: []nodes.ServiceStep{
			{
				Interface: nodes.InterfaceBIOS,
				Step:      "apply_configuration",
				Args: map[string]any{
					"settings": []string{},
				},
			},
		},
	}).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestChangePowerState(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleChangePowerStateSuccessfully(t, fakeServer)

	opts := nodes.PowerStateOpts{
		Target:  nodes.PowerOn,
		Timeout: 100,
	}

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangePowerState(context.TODO(), c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestChangePowerStateWithConflict(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleChangePowerStateWithConflict(t, fakeServer)

	opts := nodes.PowerStateOpts{
		Target:  nodes.PowerOn,
		Timeout: 100,
	}

	c := client.ServiceClient(fakeServer)
	err := nodes.ChangePowerState(context.TODO(), c, "1234asdf", opts).ExtractErr()
	if !gophercloud.ResponseCodeIs(err, http.StatusConflict) {
		t.Fatalf("expected 409 response, but got %s", err.Error())
	}
}

func TestSetRAIDConfig(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSetRAIDConfig(t, fakeServer)

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

	c := client.ServiceClient(fakeServer)
	err := nodes.SetRAIDConfig(context.TODO(), c, "1234asdf", config).ExtractErr()
	th.AssertNoErr(t, err)
}

// Without specifying a size, we need to send a string: "MAX"
func TestSetRAIDConfigMaxSize(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSetRAIDConfigMaxSize(t, fakeServer)

	isRootVolume := true

	config := nodes.RAIDConfigOpts{
		LogicalDisks: []nodes.LogicalDisk{
			{
				IsRootVolume: &isRootVolume,
				RAIDLevel:    nodes.RAID1,
			},
		},
	}

	c := client.ServiceClient(fakeServer)
	err := nodes.SetRAIDConfig(context.TODO(), c, "1234asdf", config).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestToRAIDConfigMap(t *testing.T) {
	cases := []struct {
		name     string
		opts     nodes.RAIDConfigOpts
		expected map[string]any
	}{
		{
			name: "LogicalDisks is empty",
			opts: nodes.RAIDConfigOpts{},
			expected: map[string]any{
				"logical_disks": nil,
			},
		},
		{
			name: "LogicalDisks is nil",
			opts: nodes.RAIDConfigOpts{
				LogicalDisks: nil,
			},
			expected: map[string]any{
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
						PhysicalDisks: []any{"6I:1:5", "6I:1:6", "6I:1:7"},
					},
				},
			},
			expected: map[string]any{
				"logical_disks": []map[string]any{
					{
						"raid_level":     "0",
						"size_gb":        "MAX",
						"volume_name":    "root",
						"physical_disks": []any{"6I:1:5", "6I:1:6", "6I:1:7"},
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
						PhysicalDisks: []any{
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
			expected: map[string]any{
				"logical_disks": []map[string]any{
					{
						"raid_level":  "0",
						"size_gb":     "MAX",
						"volume_name": "root",
						"controller":  "software",
						"physical_disks": []any{
							map[string]any{
								"size": "> 100",
							},
							map[string]any{
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListBIOSSettingsSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.ListBIOSSettings(context.TODO(), c, "1234asdf", nil).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeBIOSSettings, actual)
}

func TestListDetailBIOSSettings(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListDetailBIOSSettingsSuccessfully(t, fakeServer)

	opts := nodes.ListBIOSSettingsOpts{
		Detail: true,
	}

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.ListBIOSSettings(context.TODO(), c, "1234asdf", opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeDetailBIOSSettings, actual)
}

func TestGetBIOSSetting(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetBIOSSettingSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.GetBIOSSetting(context.TODO(), c, "1234asdf", "ProcVirtualization").Extract()
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetVendorPassthruMethodsSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.GetVendorPassthruMethods(context.TODO(), c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeVendorPassthruMethods, *actual)
}

func TestGetAllSubscriptions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetAllSubscriptionsVendorPassthruSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	method := nodes.CallVendorPassthruOpts{
		Method: "get_all_subscriptions",
	}
	actual, err := nodes.GetAllSubscriptions(context.TODO(), c, "1234asdf", method).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeGetAllSubscriptions, *actual)
}

func TestGetSubscription(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSubscriptionVendorPassthruSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	method := nodes.CallVendorPassthruOpts{
		Method: "get_subscription",
	}
	subscriptionOpt := nodes.GetSubscriptionOpts{
		Id: "62dbd1b6-f637-11eb-b551-4cd98f20754c",
	}
	actual, err := nodes.GetSubscription(context.TODO(), c, "1234asdf", method, subscriptionOpt).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeGetSubscription, *actual)
}

func TestCreateSubscriptionAllParameters(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSubscriptionVendorPassthruAllParametersSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	method := nodes.CallVendorPassthruOpts{
		Method: "create_subscription",
	}
	createOpt := nodes.CreateSubscriptionOpts{
		Destination: "https://someurl",
		Context:     "gophercloud",
		Protocol:    "Redfish",
		EventTypes:  []string{"Alert"},
		HttpHeaders: []map[string]string{{"Content-Type": "application/json"}},
	}
	actual, err := nodes.CreateSubscription(context.TODO(), c, "1234asdf", method, createOpt).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeCreateSubscriptionAllParameters, *actual)
}

func TestCreateSubscriptionWithRequiredParameters(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSubscriptionVendorPassthruRequiredParametersSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	method := nodes.CallVendorPassthruOpts{
		Method: "create_subscription",
	}
	createOpt := nodes.CreateSubscriptionOpts{
		Destination: "https://somedestinationurl",
	}
	actual, err := nodes.CreateSubscription(context.TODO(), c, "1234asdf", method, createOpt).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeCreateSubscriptionRequiredParameters, *actual)
}

func TestDeleteSubscription(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSubscriptionVendorPassthruSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	method := nodes.CallVendorPassthruOpts{
		Method: "delete_subscription",
	}
	deleteOpt := nodes.DeleteSubscriptionOpts{
		Id: "344a3e2-978a-444e-990a-cbf47c62ef88",
	}
	err := nodes.DeleteSubscription(context.TODO(), c, "1234asdf", method, deleteOpt).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestSetMaintenance(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleSetNodeMaintenanceSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.SetMaintenance(context.TODO(), c, "1234asdf", nodes.MaintenanceOpts{
		Reason: "I'm tired",
	}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUnsetMaintenance(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUnsetNodeMaintenanceSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.UnsetMaintenance(context.TODO(), c, "1234asdf").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetInventory(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetInventorySuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.GetInventory(context.TODO(), c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeInventoryData.Inventory, actual.Inventory)

	pluginData, err := actual.PluginData.AsMap()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "x86_64", pluginData["cpu_arch"].(string))

	compatData, err := actual.PluginData.AsInspectorData()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "x86_64", compatData.CPUArch)
}

func TestListFirmware(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListFirmwareSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.ListFirmware(context.TODO(), c, "1234asdf").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, NodeFirmwareList, actual)
}

func TestVirtualMediaOpts(t *testing.T) {
	opts := nodes.DetachVirtualMediaOpts{
		DeviceTypes: []nodes.VirtualMediaDeviceType{nodes.VirtualMediaCD, nodes.VirtualMediaDisk},
	}

	// Regular ListOpts can
	query, err := opts.ToDetachVirtualMediaOptsQuery()
	th.AssertEquals(t, "?device_types=cdrom%2Cdisk", query)
	th.AssertNoErr(t, err)
}

func TestVirtualMediaAttach(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAttachVirtualMediaSuccessfully(t, fakeServer, false)

	c := client.ServiceClient(fakeServer)
	opts := nodes.AttachVirtualMediaOpts{
		ImageURL:   "https://example.com/image",
		DeviceType: nodes.VirtualMediaCD,
	}
	err := nodes.AttachVirtualMedia(context.TODO(), c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestVirtualMediaAttachWithSource(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAttachVirtualMediaSuccessfully(t, fakeServer, true)

	c := client.ServiceClient(fakeServer)
	opts := nodes.AttachVirtualMediaOpts{
		ImageURL:            "https://example.com/image",
		DeviceType:          nodes.VirtualMediaCD,
		ImageDownloadSource: nodes.ImageDownloadSourceHTTP,
	}
	err := nodes.AttachVirtualMedia(context.TODO(), c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestVirtualMediaDetach(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDetachVirtualMediaSuccessfully(t, fakeServer, false)

	c := client.ServiceClient(fakeServer)
	err := nodes.DetachVirtualMedia(context.TODO(), c, "1234asdf", nodes.DetachVirtualMediaOpts{}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestVirtualMediaDetachWithTypes(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDetachVirtualMediaSuccessfully(t, fakeServer, true)

	c := client.ServiceClient(fakeServer)
	opts := nodes.DetachVirtualMediaOpts{
		DeviceTypes: []nodes.VirtualMediaDeviceType{nodes.VirtualMediaCD},
	}
	err := nodes.DetachVirtualMedia(context.TODO(), c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestVirtualMediaGetAttached(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetVirtualMediaSuccessfully(t, fakeServer, true)

	c := client.ServiceClient(fakeServer)
	err := nodes.GetVirtualMedia(context.TODO(), c, "1234asdf").Err
	th.AssertNoErr(t, err)
}

func TestVirtualMediaGetNotAttached(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetVirtualMediaSuccessfully(t, fakeServer, false)

	c := client.ServiceClient(fakeServer)
	err := nodes.GetVirtualMedia(context.TODO(), c, "1234asdf").Err
	th.AssertNoErr(t, err)
}

func TestListVirtualInterfaces(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListVirtualInterfacesSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	actual, err := nodes.ListVirtualInterfaces(context.TODO(), c, "1234asdf").Extract()
	th.AssertNoErr(t, err)

	expected := []nodes.VIF{
		{
			ID: "1974dcfa-836f-41b2-b541-686c100900e5",
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}

func TestAttachVirtualInterface(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAttachVirtualInterfaceSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	opts := nodes.VirtualInterfaceOpts{
		ID: "1974dcfa-836f-41b2-b541-686c100900e5",
	}
	err := nodes.AttachVirtualInterface(context.TODO(), c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAttachVirtualInterfaceWithPort(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAttachVirtualInterfaceWithPortSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	opts := nodes.VirtualInterfaceOpts{
		ID:       "1974dcfa-836f-41b2-b541-686c100900e5",
		PortUUID: "b2f96298-5172-45e9-b174-8d1ba936ab47",
	}
	err := nodes.AttachVirtualInterface(context.TODO(), c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestAttachVirtualInterfaceWithPortgroup(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleAttachVirtualInterfaceWithPortgroupSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	opts := nodes.VirtualInterfaceOpts{
		ID:            "1974dcfa-836f-41b2-b541-686c100900e5",
		PortgroupUUID: "c24944b5-a52e-4c5c-9c0a-52a0235a08a2",
	}
	err := nodes.AttachVirtualInterface(context.TODO(), c, "1234asdf", opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDetachVirtualInterface(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDetachVirtualInterfaceSuccessfully(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	err := nodes.DetachVirtualInterface(context.TODO(), c, "1234asdf", "1974dcfa-836f-41b2-b541-686c100900e5").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestVirtualInterfaceOptsValidation(t *testing.T) {
	opts := nodes.VirtualInterfaceOpts{
		ID:            "1974dcfa-836f-41b2-b541-686c100900e5",
		PortUUID:      "b2f96298-5172-45e9-b174-8d1ba936ab47",
		PortgroupUUID: "c24944b5-a52e-4c5c-9c0a-52a0235a08a2",
	}

	_, err := opts.ToVirtualInterfaceMap()
	th.AssertEquals(t, err.Error(), "cannot specify both port_uuid and portgroup_uuid")
}
