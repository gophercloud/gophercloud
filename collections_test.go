package gophercloud

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func responseWithBody(body string) (http.Response, error) {
	return http.Response{
		Body: nopCloser{bytes.NewReader([]byte(body))},
	}, nil
}

// SinglePage sample and test cases.

func ExtractSingleInts(page Page) ([]int, error) {
	var response struct {
		Ints []int `mapstructure:"ints"`
	}

	err := mapstructure.Decode(page.(SinglePage).Body, &response)
	if err != nil {
		return nil, err
	}

	return response.Ints, nil
}

func setupSinglePaged() Pager {
	return NewSinglePager(func() (http.Response, error) {
		return responseWithBody(`{ "ints": [1, 2, 3] }`)
	})
}

func TestEnumerateSinglePaged(t *testing.T) {
	callCount := 0
	err := setupSinglePaged().EachPage(func(page Page) bool {
		callCount++

		expected := []int{1, 2, 3}
		actual, err := ExtractSingleInts(page)
		if err != nil {
			t.Fatalf("Unexpected error extracting ints: %v", err)
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %v, but was %v", expected, actual)
		}
		return true
	})
	if err != nil {
		t.Fatalf("Unexpected error calling EachPage: %v", err)
	}

	if callCount != 1 {
		t.Errorf("Callback was invoked %d times", callCount)
	}
}

// LinkedPager sample and test cases.

func ExtractLinkedInts(page Page) ([]int, error) {
	var response struct {
		Ints []int `mapstructure:"ints"`
	}

	err := mapstructure.Decode(page.(LinkedPage).Body, &response)
	if err != nil {
		return nil, err
	}

	return response.Ints, nil
}

func createLinked(t *testing.T) Pager {
	return NewLinkedPager("page1", func(url string) (http.Response, error) {
		switch url {
		case "page1":
			return responseWithBody(`{ "ints": [1, 2, 3], "links": { "next": "page2" } }`)
		case "page2":
			return responseWithBody(`{ "ints": [4, 5, 6], "links": { "next": "page3" } }`)
		case "page3":
			return responseWithBody(`{ "ints": [7, 8, 9], "links": { "next": null } }`)
		default:
			t.Fatalf("LinkedPager called with unexpected URL: %v", url)
			return http.Response{}, errors.New("Wat")
		}
	})
}

func TestEnumerateLinked(t *testing.T) {
	pager := createLinked(t)

	callCount := 0
	err := pager.EachPage(func(page Page) bool {
		actual, err := ExtractLinkedInts(page)
		if err != nil {
			t.Errorf("Unable to extract ints from page: %v", err)
			return false
		}

		t.Logf("Handler invoked with %v", actual)

		var expected []int
		switch callCount {
		case 0:
			expected = []int{1, 2, 3}
		case 1:
			expected = []int{4, 5, 6}
		case 2:
			expected = []int{7, 8, 9}
		default:
			t.Fatalf("Unexpected call count: %d", callCount)
			return false
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Call %d: Expected %#v, but was %#v", callCount, expected, actual)
		}

		callCount++
		return true
	})
	if err != nil {
		t.Errorf("Unexpected error for page iteration: %v", err)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, but was %d", callCount)
	}
}
