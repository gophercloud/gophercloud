//go:build acceptance || sharedfilesystems || securityservices

package v2

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSecurityServiceCreateDelete(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	securityService, err := CreateSecurityService(t, client)
	if err != nil {
		t.Fatalf("Unable to create security service: %v", err)
	}

	newSecurityService, err := securityservices.Get(context.TODO(), client, securityService.ID).Extract()
	th.AssertNoErr(t, err)

	if newSecurityService.Name != securityService.Name {
		t.Fatalf("Security service name was expeted to be: %s", securityService.Name)
	}

	if newSecurityService.Description != securityService.Description {
		t.Fatalf("Security service description was expeted to be: %s", securityService.Description)
	}

	tools.PrintResource(t, securityService)

	defer DeleteSecurityService(t, client, securityService)
}

func TestSecurityServiceList(t *testing.T) {
	clients.RequireLong(t)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	expected := make(map[string]int)
	for range 2 {
		service, err := securityservices.Create(ctx, client, securityservices.CreateOpts{
			Name: tools.RandomString("ACPTTEST-list-", 16),
			Type: securityservices.Kerberos,
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cleanupCancel()
			if err := securityservices.Delete(cleanupCtx, client, service.ID).ExtractErr(); err != nil {
				t.Errorf("Unable to delete security service %s: %v", service.ID, err)
			}
		})
		expected[service.ID] = 1
	}

	for _, allPages := range []bool{false, true} {
		name := "EachPage"
		if allPages {
			name = "AllPages"
		}
		t.Run(name, func(t *testing.T) {
			actual := make(map[string]int)
			collect := func(page pagination.Page) error {
				items, err := securityservices.ExtractSecurityServices(page)
				if err != nil {
					return err
				}
				for _, item := range items {
					if _, ok := expected[item.ID]; ok {
						actual[item.ID]++
					}
				}
				return nil
			}
			pager := securityservices.List(client, nil)
			if allPages {
				page, err := pager.AllPages(ctx)
				th.AssertNoErr(t, err)
				th.AssertNoErr(t, collect(page))
			} else {
				pages := 0
				err := pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
					pages++
					return true, collect(page)
				})
				th.AssertNoErr(t, err)
				th.CheckEquals(t, 1, pages)
			}
			th.CheckDeepEquals(t, expected, actual)
		})
	}
}

// The test creates 2 security services and verifies that only the one(s) with
// a particular name are being listed
func TestSecurityServiceListFiltering(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create a shared file system client: %v", err)
	}

	securityService, err := CreateSecurityService(t, client)
	if err != nil {
		t.Fatalf("Unable to create security service: %v", err)
	}
	defer DeleteSecurityService(t, client, securityService)

	securityService, err = CreateSecurityService(t, client)
	if err != nil {
		t.Fatalf("Unable to create security service: %v", err)
	}
	defer DeleteSecurityService(t, client, securityService)

	options := securityservices.ListOpts{
		Name: securityService.Name,
	}

	allPages, err := securityservices.List(client, options).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to retrieve security services: %v", err)
	}

	allSecurityServices, err := securityservices.ExtractSecurityServices(allPages)
	if err != nil {
		t.Fatalf("Unable to extract security services: %v", err)
	}

	for _, listedSecurityService := range allSecurityServices {
		if listedSecurityService.Name != securityService.Name {
			t.Fatalf("The name of the security service was expected to be %s", securityService.Name)
		}
		tools.PrintResource(t, &listedSecurityService)
	}
}

// Create a security service and update the name and description. Get the security
// service and verify that the name and description have been updated
func TestSecurityServiceUpdate(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	securityService, err := CreateSecurityService(t, client)
	if err != nil {
		t.Fatalf("Unable to create security service: %v", err)
	}

	name := "NewName"
	description := ""
	options := securityservices.UpdateOpts{
		Name:        &name,
		Description: &description,
		Type:        "ldap",
	}

	_, err = securityservices.Update(context.TODO(), client, securityService.ID, options).Extract()
	th.AssertNoErr(t, err)

	newSecurityService, err := securityservices.Get(context.TODO(), client, securityService.ID).Extract()
	th.AssertNoErr(t, err)

	if newSecurityService.Name != name {
		t.Fatalf("Security service name was expeted to be: %s", name)
	}

	if newSecurityService.Description != description {
		t.Fatalf("Security service description was expeted to be: %s", description)
	}

	if newSecurityService.Type != options.Type {
		t.Fatalf("Security service type was expected to be: %s", options.Type)
	}

	tools.PrintResource(t, securityService)

	defer DeleteSecurityService(t, client, securityService)
}
