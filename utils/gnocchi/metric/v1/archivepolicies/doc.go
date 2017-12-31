/*
Package archivepolicies provides the ability to retrieve archive policies
through the Gnocchi API.

Example of Listing archive policies

	allPages, err := archivepolicies.List(gnocchiClient).AllPages()
	if err != nil {
		panic(err)
	}

	allArchivePolicies, err := archivepolicies.ExtractArchivePolicies(allPages)
	if err != nil {
		panic(err)
	}

	for _, archivePolicy := range allArchivePolicies {
		fmt.Printf("%+v\n", archivePolicy)
	}
*/
package archivepolicies
