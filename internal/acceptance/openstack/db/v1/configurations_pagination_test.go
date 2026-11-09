//go:build acceptance || db || configurations

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/configurations"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestConfigurationsPagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client, err := clients.NewDBV1Client()
	th.AssertNoErr(t, err)
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)
	tools.LimitListPages(t, client, client.ServiceURL("configurations"))

	values := make(map[string]any)
	switch choices.DBDatastoreType {
	case "redis", "keydb", "valkey":
		values["maxclients"] = 200
	default:
		values["max_connections"] = 200
	}
	var expected []string
	for range 2 {
		configuration, err := configurations.Create(ctx, client, configurations.CreateOpts{
			Name:   tools.RandomString("ACPTTEST-pagination-", 16),
			Values: values,
			Datastore: &configurations.DatastoreOpts{
				Type:    choices.DBDatastoreType,
				Version: choices.DBDatastoreVersion,
			},
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, configurations.Delete(ctx, client, configuration.ID).ExtractErr())
		})
		expected = append(expected, configuration.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return configurations.List(client)
	}, func(page pagination.Page) ([]string, error) {
		items, err := configurations.ExtractConfigs(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.ID
		}
		return ids, err
	}, expected)
}
