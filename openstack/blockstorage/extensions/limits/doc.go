/*
Package limits shows rate and limit information for a project you authorized for.

Example to Retrieve Limits

	limits, err := limits.Get(context.TODO(), blockStorageClient).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", limits)
*/
package limits
