package gophercloud

import (
	"net/url"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
)

func TestMaybeString(t *testing.T) {
	testString := ""
	var expected *string
	actual := MaybeString(testString)
	th.CheckDeepEquals(t, actual, expected)

	testString = "carol"
	expected = &testString
	actual = MaybeString(testString)
	th.CheckDeepEquals(t, actual, expected)
}

func TestMaybeInt(t *testing.T) {
	testInt := 0
	var expected *int
	actual := MaybeInt(testInt)
	th.CheckDeepEquals(t, actual, expected)

	testInt = 4
	expected = &testInt
	actual = MaybeInt(testInt)
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

	actual, err := BuildQueryString(&opts)
	if err != nil {
		t.Errorf("Error building query string: %v", err)
	}

	th.CheckDeepEquals(t, actual, expected)
}

func TestBuildQueryStringWithoutRequiredFieldSet(t *testing.T) {
	type Opts struct {
		J int    `q:"j"`
		R string `q:"r,required"`
		C bool
	}

	opts := Opts{J: 2, C: true}

	_, err := BuildQueryString(&opts)
	if err == nil {
		t.Error("Unexpected result: There should be an error thrown when a required field isn't set.")
	}

	t.Log(err)
}
