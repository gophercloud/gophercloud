package testing

import (
	"context"
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestNewInstanceHAV1(t *testing.T) {
	ctx := context.Background()
	want := gophercloud.EndpointOpts{Type: "instance-ha", Version: 1, Region: "RegionOne", Availability: gophercloud.AvailabilityInternal}
	provider := &gophercloud.ProviderClient{EndpointLocator: func(actualCtx context.Context, opts gophercloud.EndpointOpts) (string, error) {
		th.AssertEquals(t, ctx, actualCtx)
		th.AssertDeepEquals(t, want, opts)
		return "https://masakari.example/v1/", nil
	}}
	client, err := openstack.NewInstanceHAV1(ctx, provider, gophercloud.EndpointOpts{Region: want.Region, Availability: want.Availability})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "instance-ha", client.Type)
	th.AssertEquals(t, "https://masakari.example/v1/", client.Endpoint)
	wantErr := errors.New("endpoint unavailable")
	provider.EndpointLocator = func(context.Context, gophercloud.EndpointOpts) (string, error) { return "", wantErr }
	_, err = openstack.NewInstanceHAV1(ctx, provider, gophercloud.EndpointOpts{})
	th.AssertEquals(t, wantErr, err)
}
