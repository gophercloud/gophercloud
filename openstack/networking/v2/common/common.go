package common

import (
	"net/url"
	"strconv"
)

type NeutronListOpts struct {
	params url.Values
}

func (opts NeutronListOpts) MultiOptString(name string, v ...string) {
	for _, val := range v {
		opts.params.Add(name, val)
	}
}

func (opts NeutronListOpts) MultiOptBool(name string, v ...bool) {
	for _, val := range v {
		opts.params.Add(name, strconv.FormatBool(val))
	}
}

func (opts NeutronListOpts) SingleOptString(name string, v string) {
	opts.params.Set(name, v)
}

func (opts NeutronListOpts) ToQueryString() (string, error) {
	q := &url.URL{RawQuery: opts.params.Encode()}
	return q.String(), nil
}
