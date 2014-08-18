package identity

import (
	"os"
)

type Client struct {
	Endpoint  string
	Authority AuthResults
	Options   AuthOptions
}

type ClientOpts struct {
	Type    string
	Name    string
	Region  string
	URLType string
}

func (ao AuthOptions) NewClient(opts ClientOpts) (Client, error) {
	client := Client{
		Options: ao,
	}

	ar, err := Authenticate(ao)
	if err != nil {
		return client, err
	}

	client.Authority = ar

	sc, err := GetServiceCatalog(ar)
	if err != nil {
		return client, err
	}

	ces, err := sc.CatalogEntries()
	if err != nil {
		return client, err
	}

	var eps []Endpoint

	if opts.Name != "" {
		for _, ce := range ces {
			if ce.Type == opts.Type && ce.Name == opts.Name {
				eps = ce.Endpoints
			}
		}
	} else {
		for _, ce := range ces {
			if ce.Type == opts.Type {
				eps = ce.Endpoints
			}
		}
	}

	region := os.Getenv("OS_REGION_NAME")
	if opts.Region != "" {
		region = opts.Region
	}

	var rep string
	for _, ep := range eps {
		if ep.Region == region {
			switch opts.URLType {
			case "public":
				rep = ep.PublicURL
			case "private":
				rep = ep.InternalURL
			default:
				rep = ep.PublicURL
			}
		}
	}

	client.Endpoint = rep

	return client, nil
}
