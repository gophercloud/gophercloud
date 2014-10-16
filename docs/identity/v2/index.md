---
layout: page
title: Getting Started with Identity v2
---

## Tokens

A token is an arbitrary bit of text that is used to access resources. Each
token has a scope that describes which resources are accessible with it. A
token may be revoked at anytime and is valid for a finite duration.

### Generate a token

The nature of required and optional auth options will depend on your provider,
but generally the `Username` and `IdentityEndpoint` fields are always
required. Some providers will insist on a `Password` instead of an `APIKey`,
others will prefer `TenantID` over `TenantName` - so it is always worth
checking before writing your implementation in Go.

{% highlight go %}
import "github.com/rackspace/gophercloud/openstack/identity/v2/tokens"

opts := tokens.AuthOptions{
  IdentityEndpoint: "{identityURL}",
  Username:         "{username}",
  APIKey:           "{apiKey}",
}

token, err := tokens.Create(client, opts).Extract()
{% endhighlight %}

## Tenants

A tenant is a container used to group or isolate API resources. Depending on
the provider, a tenant can map to a customer, account, organization, or project.

### List tenants

{% highlight go %}
import (
  "github.com/rackspace/gophercloud/pagination"
  "github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
)

// We have the option of filtering the tenant list. If we want the full
// collection, leave it as an empty struct
opts := tenants.ListOpts{Limit: 10}

// Retrieve a pager (i.e. a paginated collection)
pager := tenants.List(client, opts)

// Define an anonymous function to be executed on each page's iteration
err := pager.EachPage(func(page pagination.Page) (bool, error) {
  tenantList, err := tenants.ExtractTenants(page)

  for _, t := range tenantList {
    // "t" will be a tenants.Tenant
  }
})
{% endhighlight %}
