package gophercloud

import (
	"net/url"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
)

func TestMaybeStringWithNonEmptyString(t *testing.T) {
	testString := "carol"
	expected := &testString
	actual := MaybeString("carol")
	th.CheckDeepEquals(t, actual, expected)
}

func TestMaybeStringWithEmptyString(t *testing.T) {
	var expected *string
	actual := MaybeString("")
	th.CheckDeepEquals(t, actual, expected)
}

func TestBuildQueryStringWithPointerToStruct(t *testing.T) {
	expected := &url.URL{
		RawQuery: "j=2&r=red",
	}

	type Opts struct {
		J int    `q:"j"`
		R string `q:"r"`
		C bool
	}

	opts := Opts{J: 2, R: "red"}

	actual, err := BuildQueryString(opts)
	if err != nil {
		t.Errorf("Error building query string: %v", err)
	}

	th.CheckDeepEquals(t, actual, expected)
}
