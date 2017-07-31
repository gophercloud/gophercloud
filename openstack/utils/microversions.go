package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ErrIncompatible is an error in case Incompatible microversion occured.
var ErrIncompatible = errors.New("Incompatible microversion")

// CompatibleMicroversion checks input microversions to find out whether they are compatible or not. The input parameters are:
// - minimumMV: minimum microversion supported by the called function.
// - maximumMV: maximum microversion supported by the called function.
// - requestedMV: microversion requested by the user of the gophercloud library.
// - serverMaximumMV: maximum microversion supported by the particular server the gophercloud library is talking to.
func CompatibleMicroversion(minimumMV, maximumMV, requestedMV, serverMaximumMV string) error {
	if requestedMV == "latest" {
		if maximumMV == "" {
			return nil
		}
		maximumMajor, maximumMinor := splitMicroversion(maximumMV)
		serverMaximumMajor, serverMaximumMinor := splitMicroversion(serverMaximumMV)
		if (maximumMajor > serverMaximumMajor) || (maximumMajor == serverMaximumMajor && maximumMinor >= serverMaximumMinor) {
			return nil
		}
		return ErrIncompatible
	}
	if err := validMicroversion(requestedMV); err != nil {
		return err
	}
	minimumMajor, minimumMinor := splitMicroversion(minimumMV)
	requestedMajor, requestedMinor := splitMicroversion(requestedMV)
	if (requestedMajor > minimumMajor) || (requestedMajor == minimumMajor && requestedMinor >= minimumMinor) {
		return nil
	}
	return ErrIncompatible
}

func splitMicroversion(mv string) (major, minor int) {
	if err := validMicroversion(mv); err != nil {
		return
	}

	mvParts := strings.Split(mv, ".")
	major, _ = strconv.Atoi(mvParts[0])
	minor, _ = strconv.Atoi(mvParts[1])

	return
}

func validMicroversion(mv string) (err error) {
	if mv == "latest" {
		return
	}

	mvRe := regexp.MustCompile("^\\d+\\.\\d+$")
	if v := mvRe.MatchString(mv); v {
		return
	}

	err = ErrIncompatible
	return
}
