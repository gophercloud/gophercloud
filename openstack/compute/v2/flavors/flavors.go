package flavors

import (
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// ErrCannotInterpret is returned by an Extract call if the response body doesn't have the expected structure.
var ErrCannotInterpet = errors.New("Unable to interpret a response body.")

// Flavor records represent (virtual) hardware configurations for server resources in a region.
type Flavor struct {
	// The Id field contains the flavor's unique identifier.
	// For example, this identifier will be useful when specifying which hardware configuration to use for a new server instance.
	ID string `mapstructure:"id"`

	// The Disk and RA< fields provide a measure of storage space offered by the flavor, in GB and MB, respectively.
	Disk int `mapstructure:"disk"`
	RAM  int `mapstructure:"ram"`

	// The Name field provides a human-readable moniker for the flavor.
	Name string `mapstructure:"name"`

	RxTxFactor float64 `mapstructure:"rxtx_factor"`

	// Swap indicates how much space is reserved for swap.
	// If not provided, this field will be set to 0.
	Swap int `mapstructure:"swap"`

	// VCPUs indicates how many (virtual) CPUs are available for this flavor.
	VCPUs int `mapstructure:"vcpus"`
}

func defaulter(from, to reflect.Kind, v interface{}) (interface{}, error) {
	if (from == reflect.String) && (to == reflect.Int) {
		return 0, nil
	}
	return v, nil
}

// ExtractFlavors provides access to the list of flavors in a page acquired from the List operation.
func ExtractFlavors(page pagination.Page) ([]Flavor, error) {
	casted := page.(ListResult).Body
	var flavors []Flavor

	cfg := &mapstructure.DecoderConfig{
		DecodeHook: defaulter,
		Result:     &flavors,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return flavors, err
	}
	err = decoder.Decode(casted)
	if err != nil {
		return flavors, err
	}

	return flavors, nil
}

// ExtractFlavor provides access to the individual flavor returned by the Get function.
func ExtractFlavor(gr GetResults) (*Flavor, error) {
	f, ok := gr["flavor"]
	if !ok {
		return nil, ErrCannotInterpet
	}

	flav := new(Flavor)
	cfg := &mapstructure.DecoderConfig{
		DecodeHook: defaulter,
		Result:     flav,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return flav, err
	}
	err = decoder.Decode(f)
	return flav, err
}
