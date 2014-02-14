package flavors

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// Flavor records represent (virtual) hardware configurations for server resources in a region.
//
// The Id field contains the flavor's unique identifier.
// For example, this identifier will be useful when specifying which hardware configuration to use for a new server instance.
//
// The Disk and Ram fields provide a measure of storage space offered by the flavor, in GB and MB, respectively.
//
// The Name field provides a human-readable moniker for the flavor.
//
// Swap indicates how much space is reserved for swap.
// If not provided, this field will be set to 0.
//
// VCpus indicates how many (virtual) CPUs are available for this flavor.
type Flavor struct {
	Disk       int
	Id         string
	Name       string
	Ram        int
	RxTxFactor float64 `mapstructure:"rxtx_factor"`
	Swap       int
	VCpus      int
}

func defaulter(from, to reflect.Kind, v interface{}) (interface{}, error) {
	if (from == reflect.String) && (to == reflect.Int) {
		return 0, nil
	}
	return v, nil
}

func GetFlavors(lr ListResults) ([]Flavor, error) {
	fa, ok := lr["flavors"]
	if !ok {
		return nil, ErrNotImplemented
	}
	fms := fa.([]interface{})

	flavors := make([]Flavor, len(fms))
	for i, fm := range fms {
		flavorObj := fm.(map[string]interface{})
		cfg := &mapstructure.DecoderConfig{
			DecodeHook: defaulter,
			Result:     &flavors[i],
		}
		decoder, err := mapstructure.NewDecoder(cfg)
		if err != nil {
			return flavors, err
		}
		err = decoder.Decode(flavorObj)
		if err != nil {
			return flavors, err
		}
	}
	return flavors, nil
}
