/*
Package trunk_details provides the ability to extend a ports result with
additional information about any trunk and subports associated with the port.

Example:

	type portExt struct {
	  ports.Port
	  trunk_details.TrunkDetailsExt
	}
	var portExt portExt

	err := ports.Get(context.TODO(), networkClient, "2ba3a709-e40e-462c-a541-85e99de589bf").ExtractInto(&portExt)
	if err != nil {
	  panic(err)
	}

	fmt.Printf("%+v\n", portExt)
*/
package trunk_details
