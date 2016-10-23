package tools

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	mrand "math/rand"
	"reflect"
	"time"
)

// ErrTimeout is returned if WaitFor takes longer than 300 second to happen.
var ErrTimeout = errors.New("Timed out")

// WaitFor polls a predicate function once per second to wait for a certain state to arrive.
func WaitFor(predicate func() (bool, error)) error {
	for i := 0; i < 300; i++ {
		time.Sleep(1 * time.Second)

		satisfied, err := predicate()
		if err != nil {
			return err
		}
		if satisfied {
			return nil
		}
	}
	return ErrTimeout
}

// MakeNewPassword generates a new string that's guaranteed to be different than the given one.
func MakeNewPassword(oldPass string) string {
	randomPassword := RandomString("", 16)
	for randomPassword == oldPass {
		randomPassword = RandomString("", 16)
	}
	return randomPassword
}

// RandomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func RandomString(prefix string, n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return prefix + string(bytes)
}

// RandomInt will return a random integer between a specified range.
func RandomInt(min, max int) int {
	mrand.Seed(time.Now().Unix())
	return mrand.Intn(max-min) + min
}

// Elide returns the first bit of its input string with a suffix of "..." if it's longer than
// a comfortable 40 characters.
func Elide(value string) string {
	if len(value) > 40 {
		return value[0:37] + "..."
	}
	return value
}

// DumpResource returns a resource as a readable structure
func DumpResource(resource interface{}) string {
	var sf []reflect.StructField
	v := reflect.ValueOf(resource).Elem()
	t := v.Type()

	// Loop through all fields of the resource and remove all tags
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		f.Tag = ""
		sf = append(sf, f)
	}

	// Create a new struct with the omitted tags
	typ := reflect.StructOf(sf)
	ns := reflect.New(typ).Elem()

	// Loop again, this time copying all values from the old resource to the new
	for i := 0; i < v.NumField(); i++ {
		copyField(v.Field(i), ns.Field(i))
	}

	// Convert the struct into JSON and return
	s := ns.Addr().Interface()
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

// copyField copies the value of the old field to the new field.
// If the value is a struct, copyField is called recursively.
func copyField(oldField reflect.Value, newField reflect.Value) {
	switch oldField.Kind() {
	case reflect.Struct:
		for i := 0; i < oldField.NumField(); i++ {
			copyField(oldField.Field(i), newField.Field(i))
		}
	default:
		newField.Set(oldField)
	}
}
