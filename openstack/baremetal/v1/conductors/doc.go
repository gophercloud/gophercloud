/*
Package conductors provides information and interaction with the conductors API
resource in the OpenStack Bare Metal service.

Example to List Conductors with Detail

	conductors.List(client, conductors.ListOpts{Detail: true}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		conductorList, err := conductors.ExtractConductors(page)
		if err != nil {
			return false, err
		}

		for _, n := range conductorList {
			// Do something
		}

		return true, nil
	})

Example to List Conductors

	listOpts := conductors.ListOpts{
		Fields:         []string{"hostname"},
	}

	conductors.List(client, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		conductorList, err := conductors.ExtractConductors(page)
		if err != nil {
			return false, err
		}

		for _, n := range conductorList {
			// Do something
		}

		return true, nil
	})

Example to Get Conductor

	showConductor, err := conductors.Get(context.TODO(), client, "compute2.localdomain").Extract()
	if err != nil {
		panic(err)
	}
*/
package conductors
