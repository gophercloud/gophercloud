package bmvolume

/*
 Package bmvolume contains the functionality to Listing, Searching, Creating, Updating,
 and Deleting of bare metal volume connectors or targets resources

 API reference: https://docs.openstack.org/api-ref/baremetal/?expanded=#volume-volume

Example to List volume:
	volume, err := bmvolume.List(client.ServiceClient(),bmvolume.ListOpts{}).Extract()
	if err != nil{
		panic(er)
	}else{
		// Do something
	}


Exampe to List volume connectors:
   connectorListOpts := bmvolume.ListConnectorsOpts{}
   connectorListOpts.Node = "6d85703a-565d-469a-96ce-30b6de53079d"
   bmvolume.ListConnectors(client.ServiceClient(), connectorListOpts).EachPage(func(page pagination.Page) (bool, error){
	    connectorList, err := bmvolume.ExtractConnectors(page)
		if err != nil{
			return false, err
		}
		for _, n := range connectorList{
			// Do something
		}
		return true,nil
}

Example to Create a volume connector:
	connectorCreateOpts := bmvolume.CreateConnectorOpts{}
	connectorCreateOpts.NodeUUID = "6d85703a-565d-469a-96ce-30b6de53079d"
	connectorCreateOpts.ConnectorType = "iqn"
	connectorCreateOpts.ConnectorId = "iqn.2017-07.org.openstack:01:d9a51732c3f"
	connector, err := bmvolume.CreateConnector(client.ServiceClient(), connectorCreateOpts).Extract()
	if err != nil{
		panic(err)
	}

Example to Get a volume connector:
	connector, err := bmvolume.GetConnector(client.ServiceClient(), "9bf93e01-d728-47a3-ad4b-5e66a835037c").Extract()
	if err != nil{
		panic(err)
	}

Example to Update a volume connector:
    connector, err := bmvolume.UpdateConnector(client.ServiceClient(), "9bf93e01-d728-47a3-ad4b-5e66a835037c", bmvolume.UpdateOpts{
		bmvolume.UpdateOperation{
			Op:    bmvolume.ReplaceOp,
			Path:  "/connector_id",
			Value: "iqn.2017-07.org.openstack:01:66666666666",
		},
	}).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a volume connector:
	err := bmvolume.DeleteConnector(client.ServiceClient(), "9bf93e01-d728-47a3-ad4b-5e66a835037c").ExtractErr()
	if err != nil {
		panic(err)
	}


Exampe to List volume targets:
   targetListOpts := bmvolume.ListTargetsOpts{}
   targetListOpts.Node = "6d85703a-565d-469a-96ce-30b6de53079d"
   bmvolume.ListTargets(client.ServiceClient(), targetListOpts).EachPage(func(page pagination.Page) (bool, error){
	    targetList, err := bmvolume.ExtractTargets(page)
		if err != nil{
			return false, err
		}
		for _, n := range targetList{
			// Do something
		}
		return true,nil
}

Example to Create a volume target:
	targetCreateOpts := bmvolume.CreateTargetOpts{}
	targetCreateOpts.BootIndex = "0"
	targetCreateOpts.NodeUUID = "6d85703a-565d-469a-96ce-30b6de53079d"
	targetCreateOpts.VolumeType = "iscsi"
	targetCreateOpts.VolumeId = "04452bed-5367-4202-8bf5-de4335ac56d2"
	target, err := bmvolume.CreateTarget(client.ServiceClient(), targetCreateOpts).Extract()
	if err != nil{
		panic(err)
	}

Example to Get a volume target:
	target, err := bmvolume.GetTarget(client.ServiceClient(), "bd4d008c-7d31-463d-abf9-6c23d9d55f7f").Extract()
	if err != nil{
		panic(err)
	}

Example to Update a volume target:
    target, err := bmvolume.UpdateTarget(client.ServiceClient(), "bd4d008c-7d31-463d-abf9-6c23d9d55f7f", bmvolume.UpdateOpts{
		bmvolume.UpdateOperation{
			Op:    bmvolume.ReplaceOp,
			Path:  "/volume_id",
			Value: "06666bed-5367-4202-8bf5-de4335ac56d2",
		},
	}).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a volume target:
	err := bmvolume.DeleteTarget(client.ServiceClient(), "bd4d008c-7d31-463d-abf9-6c23d9d55f7f").ExtractErr()
	if err != nil {
		panic(err)
	}


*/
