package v1

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

func CheckContainerName(s string) error {
	if len(s) < 1 {
		return ErrEmptyContainerName{}
	}
	if strings.ContainsRune(s, '/') {
		return ErrInvalidContainerName{name: s}
	}
	return nil
}

func CheckObjectName(s string) error {
	if s == "" {
		return ErrEmptyObjectName{}
	}
	// Check if objet name has a leading slash or two or more consecutive slashes
	// anywhere in the string
	re := regexp.MustCompile(`^\/|\/{2,}`)
	match := re.MatchString(s)
	if match {
		return ErrInvalidObjectName{name: s}
	}
	return nil
}

// ErrInvalidContainerName signals a container name containing an illegal
// character.
type ErrInvalidContainerName struct {
	name string
	gophercloud.BaseError
}

func (e ErrInvalidContainerName) Error() string {
	return fmt.Sprintf("invalid name %q: a container name cannot contain a slash (/) character", e.name)
}

// ErrEmptyContainerName signals an empty container name.
type ErrEmptyContainerName struct {
	gophercloud.BaseError
}

func (e ErrEmptyContainerName) Error() string {
	return "a container name must not be empty"
}

// ErrEmptyObjectName signals an empty container name.
type ErrEmptyObjectName struct {
	gophercloud.BaseError
}

// ErrInvalidObjectName signals an invalid object name.
type ErrInvalidObjectName struct {
	name string
	gophercloud.BaseError
}

func (e ErrInvalidObjectName) Error() string {
	return "an object name must not start with a leading slash or have two or more consecutive slashes"
}

func (e ErrEmptyObjectName) Error() string {
	return "an object name must not be empty"
}
