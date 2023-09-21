package common

import (
	"net/url"
	"strconv"
)

type Modifiers interface {
	MultiOptString(name string, v ...string)
	MultiOptBool(name string, v ...bool)
	SingleOptString(name string, v string)
}

type NeutronListOptsConfig func(Modifiers)

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

func (opts NeutronListOpts) ToQueryString() string {
	q := &url.URL{RawQuery: opts.params.Encode()}
	return q.String()
}

func MultiOptString(name string, v ...string) NeutronListOptsConfig {
	return func(opts Modifiers) {
		opts.MultiOptString(name, v...)
	}
}

func MultiOptBool(name string, v ...bool) NeutronListOptsConfig {
	return func(opts Modifiers) {
		opts.MultiOptBool(name, v...)
	}
}

func SingleOptString(name string, v string) NeutronListOptsConfig {
	return func(opts Modifiers) {
		opts.SingleOptString(name, v)
	}
}

func Tags(v ...string) NeutronListOptsConfig {
	return MultiOptString("tags", v...)
}

func TagsAny(v ...string) NeutronListOptsConfig {
	return MultiOptString("tags-any", v...)
}

func NotTags(v ...string) NeutronListOptsConfig {
	return MultiOptString("not-tags", v...)
}

func NotTagsAny(v ...string) NeutronListOptsConfig {
	return MultiOptString("not-tags-any", v...)
}

func SortKey(v string) NeutronListOptsConfig {
	return SingleOptString("sort_key", v)
}

func SortDir(v string) NeutronListOptsConfig {
	return SingleOptString("sort_dir", v)
}

func Marker(v string) NeutronListOptsConfig {
	return SingleOptString("marker", v)
}

func Limit(v int) NeutronListOptsConfig {
	return SingleOptString("limit", strconv.Itoa(v))
}
