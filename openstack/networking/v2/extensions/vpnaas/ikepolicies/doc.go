/*
Package ikepolicies allows management and retrieval of IKE policies in the
OpenStack Networking Service.


Example to Create an IKE policy

	createOpts := ikepolicies.CreateOpts{
		Name:        "ikepolicy1",
		Description: "Description of ikepolicy1",
		EncryptionAlgorithm: ikepolicies.EncryptionAlgorithm3DES,
		PFS:                 ikepolicies.PFSGroup5,
	}

	policy, err := ikepolicies.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

*/
package ikepolicies
