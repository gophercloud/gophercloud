package resourceproviders

import "github.com/gophercloud/gophercloud/pagination"

type ResourceProviderLinks struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

// ResourceProvider are entities which provider consumable inventory of one or more classes of resource
type ResourceProvider struct {
	Generation         int                     `json:"generation"`
	UUID               string                  `json:"uuid"`
	Links              []ResourceProviderLinks `json:"links"`
	Name               string                  `json:"name"`
	ParentProviderUuid string                  `json:"parent_provider_uuid"`
	RootProviderUuid   string                  `json:"root_provider_uuid"`
}

// ResourceProvidersPage contains a single page of all resource providers from a List call.
type ResourceProvidersPage struct {
	pagination.SinglePageBase
}

// IsEmpty determines if a ResourceProvidersPage contains any results.
func (page ResourceProvidersPage) IsEmpty() (bool, error) {
	resourceProviders, err := ExtractResourceProviders(page)
	return len(resourceProviders) == 0, err
}

// ExtractResourceProviders returns a slice of ResourceProvider from a List operation.
func ExtractResourceProviders(r pagination.Page) ([]ResourceProvider, error) {
	var s struct {
		ResourceProviders []ResourceProvider `json:"resource_providers"`
	}
	err := (r.(ResourceProvidersPage)).ExtractInto(&s)
	return s.ResourceProviders, err
}
