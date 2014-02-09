package identity

import "github.com/mitchellh/mapstructure"

type ServiceCatalogDesc struct {
	serviceDescriptions []interface{}
}

type CatalogEntry struct {
	Name      string
	Type      string
	Endpoints []Endpoint
}

type Endpoint struct {
	TenantId    string
	PublicURL   string
	InternalURL string
	Region      string
	VersionId   string
	VersionInfo string
	VersionList string
}

func ServiceCatalog(ar AuthResults) (*ServiceCatalogDesc, error) {
	access := ar["access"].(map[string]interface{})
	sds := access["serviceCatalog"].([]interface{})
	sc := &ServiceCatalogDesc{
		serviceDescriptions: sds,
	}
	return sc, nil
}

func (sc *ServiceCatalogDesc) NumberOfServices() int {
	return len(sc.serviceDescriptions)
}

func (sc *ServiceCatalogDesc) CatalogEntries() ([]CatalogEntry, error) {
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
