package gophercloud

import (
	"strings"
	"testing"
)

type OptForTest struct {
	Data0 map[string]string `q:"data0"`
	Data1 map[string]int    `q:"data1"`
	Data2 map[string]bool   `q:"data2"`
	Data3 map[string][]int  `q:"data3"`
}

func TestBuildQueryString(t *testing.T) {

	t.Run("success_map", func(t *testing.T) {
		opt := OptForTest{
			Data0: map[string]string{"k1": "success1"},
			Data1: map[string]int{"k2": 10},
			Data2: map[string]bool{"k3": true},
		}
		r, e := BuildQueryString(opt)
		if e != nil {
			t.Fatalf("unexpect err : %v", e)
		}
		c := "?data0=%7B%27k1%27%3A%27success1%27%7D&data1=%7B%27k2%27%3A%2710%27%7D&data2=%7B%27k3%27%3A%27true%27%7D"
		if r.String() != c {
			t.Fatalf("expect %q, but receive %q", c, r.String())
		}
	})

	t.Run("fail_map", func(t *testing.T) {
		opt := OptForTest{
			Data3: map[string][]int{"kops6": {1, 2}},
		}
		_, e := BuildQueryString(opt)
		if e == nil {
			t.Fatalf("unexpect success")
		}
		if !strings.HasSuffix(e.Error(), "not support value type []int") {
			t.Fatalf("receive unexpected err: %v", e)
		}
	})
}
