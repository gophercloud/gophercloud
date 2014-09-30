package accounts

import (
	"net/http"
	"reflect"
	"testing"
)

func TestExtractAccountMetadata(t *testing.T) {
	getResult := &http.Response{}

	expected := map[string]string{}

	actual := ExtractMetadata(getResult)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
	}
}
