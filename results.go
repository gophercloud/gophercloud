package gophercloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

/*
Result is an internal type to be used by individual resource packages, but its
methods will be available on a wide variety of user-facing embedding types.

It acts as a base struct that other Result types, returned from request
functions, can embed for convenience. All Results capture basic information
from the HTTP transaction that was performed, including the response body,
HTTP headers, and any errors that happened.

Generally, each Result type will have an Extract method that can be used to
further interpret the result's payload in a specific context. Extensions or
providers can then provide additional extraction functions to pull out
provider- or extension-specific information as well.
*/
type Result struct {
	// Body is the payload of the HTTP response from the server. In most cases,
	// this will be the deserialized JSON structure.
	Body interface{}

	// Header contains the HTTP header structure from the original response.
	Header http.Header

	// Err is an error that occurred during the operation. It's deferred until
	// extraction to make it easier to chain the Extract call.
	Err error
}

// ExtractInto allows users to provide an object into which `Extract` will extract
// the `Result.Body`. This would be useful for OpenStack providers that have
// different fields in the response object than OpenStack proper.
func (r Result) ExtractInto(to interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	if reader, ok := r.Body.(io.Reader); ok {
		if readCloser, ok := reader.(io.Closer); ok {
			defer readCloser.Close()
		}
		return json.NewDecoder(reader).Decode(to)
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, to)
}

// ExtractIntoStructPtr will unmarshal the Result (r) into the provided
// interface{} (to).
//
// NOTE: For internal use only
//
// `to` must be a pointer to an underlying struct type
//
// If provided, `label` will be filtered out of the response
// body prior to `r` being unmarshalled into `to`.
func (r Result) ExtractIntoStructPtr(to interface{}, label string) error {
	if r.Err != nil {
		return r.Err
	}

	t := reflect.TypeOf(to)
	if k := t.Kind(); k != reflect.Ptr {
		return fmt.Errorf("Expected pointer, got %v", k)
	}

	tov := reflect.ValueOf(to)

	var rawitem map[string]interface{}
	switch bodyt := r.Body.(type) {
	case map[string]interface{}:
		if label == "" {
			rawitem = bodyt
		} else {
			rawmap := bodyt[label]
			switch rawmapt := rawmap.(type) {
			case map[string]interface{}:
				rawitem = rawmapt
			case nil:
				return nil
			default:
				return fmt.Errorf("unsupported type: %T", rawmapt)
			}
		}
	default:
		return fmt.Errorf("unsupported type for r.Body: %T", r.Body)
	}

	b, err := json.Marshal(rawitem)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, to)
	if err != nil {
		return err
	}

	switch t.Elem().Kind() {
	case reflect.Struct:
		for i := 0; i < t.Elem().NumField(); i++ {
			if f := t.Elem().Field(i); f.Anonymous {
				var jumer json.Unmarshaler
				s := reflect.PtrTo(f.Type)
				p0 := reflect.New(f.Type)
				itemb, err := json.Marshal(rawitem)
				if err != nil {
					return err
				}
				p1 := reflect.ValueOf(itemb)
				if s.Implements(reflect.TypeOf(&jumer).Elem()) {
					method, _ := s.MethodByName("UnmarshalJSON")
					params := []reflect.Value{p0, p1}
					rvs := method.Func.Call(params)
					if err, ok := rvs[0].Interface().(error); ok && err != nil {
						return err
					}
					tov.Elem().FieldByName(f.Name).Set(reflect.Indirect(params[0]))
				} else {
					err := json.Unmarshal(itemb, p0.Interface())
					if err != nil {
						return err
					}
					tov.Elem().FieldByName(f.Name).Set(reflect.Indirect(p0))
				}
			}
		}
	default:
		return fmt.Errorf("Expected pointer to struct, got: %v", t)
	}

	return nil
}

// ExtractIntoSlicePtr will unmarshal the Result (r) into the provided
// interface{} (to).
//
// NOTE: For internal use only
//
// `to` must be a pointer to an underlying slice type
//
// If provided, `label` will be filtered out of the response
// body prior to `r` being unmarshalled into `to`.
func (r Result) ExtractIntoSlicePtr(to interface{}, label string) error {
	if r.Err != nil {
		return r.Err
	}

	t := reflect.TypeOf(to)
	if k := t.Kind(); k != reflect.Ptr {
		return fmt.Errorf("Expected pointer, got %v", k)
	}

	eltype := t.Elem()
	if eltype.Kind() != reflect.Slice {
		return fmt.Errorf("Expected pointer to slice, got: %v", t)
	}

	tov := reflect.ValueOf(to)

	var rawitems []interface{}
	switch bodyt := r.Body.(type) {
	case map[string][]interface{}:
		if label == "" {
			return fmt.Errorf("expected label for map[string]interface{} page")
		}
		rawitems = bodyt[label]
	case map[string]interface{}:
		if label == "" {
			return fmt.Errorf("expected label for map[string]interface{} page")
		}
		rawitems = bodyt[label].([]interface{})
	case []interface{}:
		rawitems = bodyt
	default:
		return fmt.Errorf("unsupported type for r.Body: %T", r.Body)
	}

	if len(rawitems) == 0 {
		return nil
	}

	switch rawitems[0].(type) {
	case map[string]interface{}:
	default:
		return fmt.Errorf("don't know how to handle extracting type: []%T", rawitems[0])
	}

	b, err := json.Marshal(rawitems)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, to)
	if err != nil {
		return err
	}

	eleltype := eltype.Elem()
	switch eleltype.Kind() {
	case reflect.Struct:
		for i := 0; i < eleltype.NumField(); i++ {
			if f := eleltype.Field(i); f.Anonymous {
				var jumer json.Unmarshaler
				s := reflect.PtrTo(f.Type)
				p0 := reflect.New(f.Type)
				for i, rawitem := range rawitems {
					itemb, err := json.Marshal(rawitem)
					if err != nil {
						return err
					}
					p1 := reflect.ValueOf(itemb)
					if s.Implements(reflect.TypeOf(&jumer).Elem()) {
						method, _ := s.MethodByName("UnmarshalJSON")
						params := []reflect.Value{p0, p1}
						rvs := method.Func.Call(params)
						if err, ok := rvs[0].Interface().(error); ok && err != nil {
							return err
						}
						tov.Elem().Index(i).FieldByName(f.Name).Set(reflect.Indirect(params[0]))
					} else {
						err := json.Unmarshal(itemb, p0.Interface())
						if err != nil {
							return err
						}
						tov.Elem().Index(i).FieldByName(f.Name).Set(reflect.Indirect(p0))
					}
				}
			}
		}
	}

	return nil
}

// PrettyPrintJSON creates a string containing the full response body as
// pretty-printed JSON. It's useful for capturing test fixtures and for
// debugging extraction bugs. If you include its output in an issue related to
// a buggy extraction function, we will all love you forever.
func (r Result) PrettyPrintJSON() string {
	pretty, err := json.MarshalIndent(r.Body, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	return string(pretty)
}

// ErrResult is an internal type to be used by individual resource packages, but
// its methods will be available on a wide variety of user-facing embedding
// types.
//
// It represents results that only contain a potential error and
// nothing else. Usually, if the operation executed successfully, the Err field
// will be nil; otherwise it will be stocked with a relevant error. Use the
// ExtractErr method
// to cleanly pull it out.
type ErrResult struct {
	Result
}

// ExtractErr is a function that extracts error information, or nil, from a result.
func (r ErrResult) ExtractErr() error {
	return r.Err
}

/*
HeaderResult is an internal type to be used by individual resource packages, but
its methods will be available on a wide variety of user-facing embedding types.

It represents a result that only contains an error (possibly nil) and an
http.Header. This is used, for example, by the objectstorage packages in
openstack, because most of the operations don't return response bodies, but do
have relevant information in headers.
*/
type HeaderResult struct {
	Result
}

// ExtractHeader will return the http.Header and error from the HeaderResult.
//
//   header, err := objects.Create(client, "my_container", objects.CreateOpts{}).ExtractHeader()
func (r HeaderResult) ExtractInto(to interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	tmpHeaderMap := map[string]string{}
	for k, v := range r.Header {
		if len(v) > 0 {
			tmpHeaderMap[k] = v[0]
		}
	}

	b, err := json.Marshal(tmpHeaderMap)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)

	return err
}

// RFC3339Milli describes a common time format used by some API responses.
const RFC3339Milli = "2006-01-02T15:04:05.999999Z"

type JSONRFC3339Milli time.Time

func (jt *JSONRFC3339Milli) UnmarshalJSON(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	t, err := time.Parse(RFC3339Milli, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339Milli(t)
	return nil
}

const RFC3339MilliNoZ = "2006-01-02T15:04:05.999999"

type JSONRFC3339MilliNoZ time.Time

func (jt *JSONRFC3339MilliNoZ) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339MilliNoZ, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339MilliNoZ(t)
	return nil
}

type JSONRFC1123 time.Time

func (jt *JSONRFC1123) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC1123, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC1123(t)
	return nil
}

type JSONUnix time.Time

func (jt *JSONUnix) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	unix, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	t = time.Unix(unix, 0)
	*jt = JSONUnix(t)
	return nil
}

// RFC3339NoZ is the time format used in Heat (Orchestration).
const RFC3339NoZ = "2006-01-02T15:04:05"

type JSONRFC3339NoZ time.Time

func (jt *JSONRFC3339NoZ) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339NoZ, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339NoZ(t)
	return nil
}

/*
Link is an internal type to be used in packages of collection resources that are
paginated in a certain way.

It's a response substructure common to many paginated collection results that is
used to point to related pages. Usually, the one we care about is the one with
Rel field set to "next".
*/
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

/*
ExtractNextURL is an internal function useful for packages of collection
resources that are paginated in a certain way.

It attempts to extract the "next" URL from slice of Link structs, or
"" if no such URL is present.
*/
func ExtractNextURL(links []Link) (string, error) {
	var url string

	for _, l := range links {
		if l.Rel == "next" {
			url = l.Href
		}
	}

	if url == "" {
		return "", nil
	}

	return url, nil
}
