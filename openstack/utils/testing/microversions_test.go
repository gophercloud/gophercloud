package testing

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/utils"
)

func TestCompatibleMicroversion(t *testing.T) {
	functionUnderTest := "utils.CompatibleMicroversion"
	tests := []struct {
		minMicroversion       string
		maxMicroversion       string
		requestedMicroversion string
		serverMaxMicroversion string
		want                  error
	}{
		{"2.7", "", "latest", "2.22", nil},
		{"2.7", "", "2.5", "2.22", utils.ErrIncompatible},
		{"2.7", "2.14", "latest", "2.22", utils.ErrIncompatible},
		{"2.7", "2.14", "latest", "2.11", nil},
		{"2.7", "", "2.11.5", "2.22", utils.ErrIncompatible},
		{"2.7", "", "2.11", "2.22", nil},
	}
	for i, tt := range tests {
		fmt.Printf("%v", i)
		if got := utils.CompatibleMicroversion(tt.minMicroversion, tt.maxMicroversion, tt.requestedMicroversion, tt.serverMaxMicroversion); got != tt.want {
			if got == nil {
				t.Errorf("%v(%q, %q, %q, %q) = (%v) WANT (%v)", functionUnderTest, tt.minMicroversion, tt.maxMicroversion, tt.requestedMicroversion, tt.serverMaxMicroversion, got, tt.want)
			} else {
				t.Errorf("%v(%q, %q, %q, %q) = (%v) WANT (%v)", functionUnderTest, tt.minMicroversion, tt.maxMicroversion, tt.requestedMicroversion, tt.serverMaxMicroversion, got, tt.want)
			}
		}
	}
}
