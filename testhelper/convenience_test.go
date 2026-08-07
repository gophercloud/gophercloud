package testhelper

import (
	"reflect"
	"testing"
)

func TestAssertTagsEqual(t *testing.T) {
	expected := []string{"tag1", "tag2"}
	actual := []string{"tag2", "tag1"}

	AssertTagsEqual(t, expected, actual)

	if !reflect.DeepEqual(expected, []string{"tag1", "tag2"}) {
		t.Fatalf("expected slice was modified: %v", expected)
	}
	if !reflect.DeepEqual(actual, []string{"tag2", "tag1"}) {
		t.Fatalf("actual slice was modified: %v", actual)
	}
}
