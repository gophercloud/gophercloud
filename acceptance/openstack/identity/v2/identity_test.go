// +build acceptance

package v2

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"

	"github.com/rackspace/gophercloud"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
	"github.com/rackspace/gophercloud/openstack/utils"
)

type extractor func(*identity.Token) string

func TestAuthentication(t *testing.T) {
	// Create an initialized set of authentication options based on available OS_*
	// environment variables.
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Error(err)
		return
	}

	// Attempt to authenticate with them.
	client := &gophercloud.ServiceClient{Endpoint: ao.IdentityEndpoint + "/"}
	r, err := identity.Authenticate(client, ao)
	if err != nil {
		t.Error(err)
		return
	}

	// We're authenticated; now let's grab our authentication token.
	tok, err := identity.GetToken(r)
	if err != nil {
		t.Error(err)
		return
	}

	// Authentication tokens have a variety of fields which might be of some interest.
	// Let's print a few of them out.
	table := map[string]extractor{
		"ID":      func(t *identity.Token) string { return tok.ID },
		"Expires": func(t *identity.Token) string { return tok.Expires },
	}

	for attr, fn := range table {
		fmt.Printf("Your token's %s is %s\n", attr, fn(tok))
	}

	// With each authentication, you receive a master directory of all the services
	// your account can access.  This "service catalog", as OpenStack calls it,
	// provides you the means to exploit other OpenStack services.
	sc, err := identity.GetServiceCatalog(r)
	if err != nil {
		t.Error(err)
		return
	}

	// Prepare our elastic tabstopped writer for our table.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 2, 8, 2, ' ', 0)

	// Different providers will provide different services.  Let's print them
	// in summary.
	ces, err := sc.CatalogEntries()
	fmt.Println("Service Catalog Summary:")
	fmt.Fprintln(w, "Name\tType\t")
	for _, ce := range ces {
		fmt.Fprintf(w, "%s\t%s\t\n", ce.Name, ce.Type)
	}
	w.Flush()

	// Now let's print them in greater detail.
	for _, ce := range ces {
		fmt.Printf("Endpoints for %s/%s\n", ce.Name, ce.Type)
		fmt.Fprintln(w, "Version\tRegion\tTenant\tPublic URL\tInternal URL\t")
		for _, ep := range ce.Endpoints {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", ep.VersionID, ep.Region, ep.TenantID, ep.PublicURL, ep.InternalURL)
		}
		w.Flush()
	}
}

func TestExtensions(t *testing.T) {
	// Create an initialized set of authentication options based on available OS_*
	// environment variables.
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Error(err)
		return
	}

	// Attempt to query extensions.
	client := &gophercloud.ServiceClient{Endpoint: ao.IdentityEndpoint + "/"}
	exts, err := identity.GetExtensions(client, ao)
	if err != nil {
		t.Error(err)
		return
	}

	// Print out a summary of supported extensions
	aliases, err := exts.Aliases()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("Extension Aliases:")
	for _, alias := range aliases {
		fmt.Printf("  %s\n", alias)
	}
}
