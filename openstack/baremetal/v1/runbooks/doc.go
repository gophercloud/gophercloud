/*
Package runbooks contains the functionality for listing, getting, creating,
updating and deleting bare metal runbooks, and for managing runbook traits.

Runbooks require microversion 1.92 or later. The Description field and the
trait functions require microversion 1.112 or later.

API reference: https://docs.openstack.org/api-ref/baremetal/#runbooks-runbooks

Example to List Runbooks with Detail

	runbooks.List(client, runbooks.ListOpts{Detail: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		runbookList, err := runbooks.ExtractRunbooks(page)
		if err != nil {
			return false, err
		}

		for _, r := range runbookList {
			// Do something
		}

		return true, nil
	})

Example to Create a Runbook

	createOpts := runbooks.CreateOpts{
		Name: "CUSTOM_AWESOME",
		Steps: []runbooks.Step{
			{
				Interface: nodes.InterfaceBIOS,
				Step:      "apply_configuration",
				Args: map[string]any{
					"settings": []map[string]string{{"name": "LogicalProc", "value": "Enabled"}},
				},
				Order: 1,
			},
		},
	}

	runbook, err := runbooks.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Get a Runbook

	runbook, err := runbooks.Get(context.TODO(), client, "CUSTOM_AWESOME").Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Runbook

	updateOpts := runbooks.UpdateOpts{
		runbooks.UpdateOperation{
			Op:    runbooks.ReplaceOp,
			Path:  "/name",
			Value: "CUSTOM_AWESOME2",
		},
	}

	runbook, err := runbooks.Update(context.TODO(), client, "CUSTOM_AWESOME", updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Set the Traits of a Runbook

	err := runbooks.SetTraits(context.TODO(), client, "CUSTOM_AWESOME", []string{"CUSTOM_GPU"}).ExtractErr()
	if err != nil {
		panic(err)
	}

Example to Delete a Runbook

	err := runbooks.Delete(context.TODO(), client, "CUSTOM_AWESOME").ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package runbooks
