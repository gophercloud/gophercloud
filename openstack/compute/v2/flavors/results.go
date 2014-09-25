package flavors

import (
	"errors"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ErrCannotInterpret is returned by an Extract call if the response body doesn't have the expected structure.
var ErrCannotInterpet = errors.New("Unable to interpret a response body.")

// GetResult temporarily holds the reponse from a Get call.
type GetResult struct {
	gophercloud.CommonResult
}

// Extract provides access to the individual Flavor returned by the Get function.
func (gr GetResult) Extract() (*Flavor, error) {
	if gr.Err != nil {
		return nil, gr.Err
	}

	var result struct {
		Flavor Flavor `mapstructure:"flavor"`
	}

	cfg := &mapstructure.DecoderConfig{
		DecodeHook: defaulter,
		Result:     &result,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(gr.Resp)
	return &result.Flavor, err
}

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
	casted := page.(ListPage).Body
	var container struct {
		Flavors []Flavor `mapstructure:"flavors"`
	}

	cfg := &mapstructure.DecoderConfig{
		DecodeHook: defaulter,
		Result:     &container,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return container.Flavors, err
	}
	err = decoder.Decode(casted)
	if err != nil {
		return container.Flavors, err
	}

	return container.Flavors, nil
}
