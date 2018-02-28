/*
Package ikepolicies allows management and retrieval of IKE policies in the
OpenStack Networking Service.


Example to Create an IKE policy

	createOpts := ikepolicies.CreateOpts{
		Name:                "ikepolicy1",
		Description:         "Description of ikepolicy1",
		EncryptionAlgorithm: ikepolicies.EncryptionAlgorithm3DES,
		PFS:                 ikepolicies.PFSGroup5,
	}

	policy, err := ikepolicies.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Show the details of a specific IKE policy by ID

	policy, err := ikepolicies.Get(client, "f2b08c1e-aa81-4668-8ae1-1401bcb0576c").Extract()
	if err != nil {
		panic(err)
	}

*/
package ikepolicies
