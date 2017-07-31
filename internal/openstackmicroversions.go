package internal

import (
	"fmt"
	"strconv"
	"strings"
)

// Microversion represents OpenStack Microversion (for more details refer to https://developer.openstack.org/api-guide/compute/microversions.html)
type Microversion struct {
	Major uint
	Minor uint
}

// New returns:
// - Either creates a new Microversion in case inStr string is a valid Microversion that has format Major.Minor, all parts, i.e. Major, '.' and Minor are compulsory.
// - Or an error in case inStr string is not a valid Microversion.
func New(inStr string) (*Microversion, error) {
	subVersions := strings.Split(inStr, ".")
	if len(subVersions) != 2 {
		return nil, fmt.Errorf("Invalid microversion %q; required format X.Y where both X and Y must be numbers greater or equal to zero", inStr)
	}
	var major, minor int
	var err error
	if major, err = strconv.Atoi(subVersions[0]); err != nil {
		return nil, fmt.Errorf("Invalid microversion %q; required format X.Y where both X and Y must be numbers greater or equal to zero", inStr)
	} else {
		if major < 0 {
			return nil, fmt.Errorf("Invalid microversion %q; required format X.Y where both X and Y must be numbers greater or equal to zero", inStr)
		}
	}
	if minor, err = strconv.Atoi(subVersions[1]); err != nil {
		return nil, fmt.Errorf("Invalid microversion %q; required format X.Y where both X and Y must be numbers greater or equal to zero", inStr)
	} else {
		if minor < 0 {
			return nil, fmt.Errorf("Invalid microversion %q; required format X.Y where both X and Y must be numbers greater or equal to zero", inStr)
		}
	}
	return &Microversion{Major: uint(major), Minor: uint(minor)}, nil
}

// LessThan compares itself (v1) with another Microversion v2 provided as an input parameter. Returns:
// - Either true in case v1 < v2.
// - Or false in case v1 >= v2.
func (v *Microversion) LessThan(greater *Microversion) bool {
	if v.Major > greater.Major {
		return false
	}
	if v.Major == greater.Major && v.Minor > greater.Minor {
		return false
	}
	if v.Major == greater.Major && v.Minor == greater.Minor {
		return false
	}
	return true
}

// String returns the Microversion as a string
func (v *Microversion) String() string {
	return strconv.Itoa(int(v.Major)) + "." + strconv.Itoa(int(v.Minor))
}
