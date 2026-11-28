/*
Package networksegmentranges provides the ability to retrieve and manage
network segment ranges through the Neutron API.

Network segment ranges define pools of segmentation IDs that can be used for
dynamic segment allocation in provider networks.

Example of Listing Network Segment Ranges

	listOpts := networksegmentranges.ListOpts{
		NetworkType: "vxlan",
	}

	allPages, err := networksegmentranges.List(networkClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allRanges, err := networksegmentranges.ExtractNetworkSegmentRanges(allPages)
	if err != nil {
		panic(err)
	}

	for _, r := range allRanges {
		fmt.Printf("%+v\n", r)
	}

Example to Get a Network Segment Range

	rangeID := "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5"
	segmentRange, err := networksegmentranges.Get(context.TODO(), networkClient, rangeID).Extract()
	if err != nil {
		panic(err)
	}

Example to Create a Network Segment Range

	createOpts := networksegmentranges.CreateOpts{
		Name:        "vxlan-range",
		NetworkType: "vxlan",
		Minimum:     1000,
		Maximum:     2000,
	}
	segmentRange, err := networksegmentranges.Create(context.TODO(), networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Update a Network Segment Range

	rangeID := "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5"
	name := "updated-range"
	minimum := 1500
	updateOpts := networksegmentranges.UpdateOpts{
		Name:    &name,
		Minimum: &minimum,
	}
	segmentRange, err := networksegmentranges.Update(context.TODO(), networkClient, rangeID, updateOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Network Segment Range

	rangeID := "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5"
	err := networksegmentranges.Delete(context.TODO(), networkClient, rangeID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package networksegmentranges
