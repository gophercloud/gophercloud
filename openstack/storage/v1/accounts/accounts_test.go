package accounts

import (
	"net/http"
	"testing"
)

func TestExtractMetadata(t *testing.T) {
	getResult := &http.Response{}

	expected := map[string]string{}

	actual := ExtractMetadata(getResult)

	for key, value := range expected {
		if value != actual[key] {
			t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
		}
	}
}
