// +build acceptance

package v2

import (
	"errors"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
)

func newClient() (*gophercloud.ServiceClient, error) {
	// Obtain credentials from the environment.
	options := gophercloud.AuthOptions{
		Username: os.Getenv("RS_USERNAME"),
		APIKey:   os.Getenv("RS_APIKEY"),
	}
	region := os.Getenv("RS_REGION")

	if options.Username == "" {
		return nil, errors.New("Please provide a Rackspace username as RS_USERNAME.")
	}
	if options.APIKey == "" {
		return nil, errors.New("Please provide a Rackspace API key as RS_APIKEY.")
	}
	if region == "" {
		return nil, errors.New("Please provide a Rackspace region as RS_REGION.")
	}

	client, err := rackspace.AuthenticatedClient(options)
	if err != nil {
		return nil, err
	}

	return rackspace.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: region,
	})
}

type serverOpts struct {
	imageID  string
	flavorID string
}

func optionsFromEnv() (*serverOpts, error) {
	options := &serverOpts{
		imageID:  os.Getenv("RS_IMAGE_ID"),
		flavorID: os.Getenv("RS_FLAVOR_ID"),
	}
	if options.imageID == "" {
		return nil, errors.New("Please provide a valid Rackspace image ID as RS_IMAGE_ID")
	}
	if options.flavorID == "" {
		return nil, errors.New("Please provide a valid Rackspace flavor ID as RS_FLAVOR_ID")
	}
	return options, nil
}
