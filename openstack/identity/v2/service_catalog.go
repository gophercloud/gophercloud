package v2

import "github.com/mitchellh/mapstructure"

// ServiceCatalog provides a view into the service catalog from a previous, successful authentication.
// OpenStack extensions may alter the structure of the service catalog in ways unpredictable to Go at compile-time,
// so this structure serves as a convenient anchor for type-safe accessors and methods.
type ServiceCatalog struct {
	serviceDescriptions []interface{}
}

// CatalogEntry provides a type-safe interface to an Identity API V2 service
// catalog listing.  Each class of service, such as cloud DNS or block storage
// services, will have a single CatalogEntry representing it.
//
// Name will contain the provider-specified name for the service.
//
// If OpenStack defines a type for the service, this field will contain that
// type string.  Otherwise, for provider-specific services, the provider may
// assign their own type strings.
//
// Endpoints will let the caller iterate over all the different endpoints that
// may exist for the service.
//
// Note: when looking for the desired service, try, whenever possible, to key
// off the type field.  Otherwise, you'll tie the representation of the service
// to a specific provider.
type CatalogEntry struct {
	Name      string
	Type      string
	Endpoints []Endpoint
}

// Endpoint represents a single API endpoint offered by a service.
// It provides the public and internal URLs, if supported, along with a region specifier, again if provided.
// The significance of the Region field will depend upon your provider.
//
// In addition, the interface offered by the service will have version information associated with it
// through the VersionId, VersionInfo, and VersionList fields, if provided or supported.
//
// In all cases, fields which aren't supported by the provider and service combined will assume a zero-value ("").
type Endpoint struct {
	TenantID    string
	PublicURL   string
	InternalURL string
	Region      string
	VersionID   string
	VersionInfo string
	VersionList string
}

// GetServiceCatalog acquires the service catalog from a successful authentication's results.
func GetServiceCatalog(ar AuthResults) (*ServiceCatalog, error) {
	access := ar["access"].(map[string]interface{})
	sds := access["serviceCatalog"].([]interface{})
	sc := &ServiceCatalog{
		serviceDescriptions: sds,
	}
	return sc, nil
}

// NumberOfServices yields the number of services the caller may use.  Note
// that this does not necessarily equal the number of endpoints available for
// use.
func (sc *ServiceCatalog) NumberOfServices() int {
	return len(sc.serviceDescriptions)
}

// CatalogEntries returns a slice of service catalog entries.
// Each entry corresponds to a specific class of service offered by the API provider.
// See the CatalogEntry structure for more details.
func (sc *ServiceCatalog) CatalogEntries() ([]CatalogEntry, error) {
	var err error
	ces := make([]CatalogEntry, sc.NumberOfServices())
	for i, sd := range sc.serviceDescriptions {
		d := sd.(map[string]interface{})
		eps, err := parseEndpoints(d["endpoints"].([]interface{}))
		if err != nil {
			return ces, err
		}
		ces[i] = CatalogEntry{
			Name:      d["name"].(string),
			Type:      d["type"].(string),
			Endpoints: eps,
		}
	}
	return ces, err
}

func parseEndpoints(eps []interface{}) ([]Endpoint, error) {
	var err error
	result := make([]Endpoint, len(eps))
	for i, ep := range eps {
		e := Endpoint{}
		err = mapstructure.Decode(ep, &e)
		if err != nil {
			return result, err
		}
		result[i] = e
	}
	return result, err
}
