package testing

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/internal"
)

func TestInvalidMicroversions(t *testing.T) {
	tests := []string{
		"",
		"2",
		"2.0.0",
		"a.0",
		"-1.3",
		"0.b",
		"3.-3",
	}
	for _, tt := range tests {
		if got, err := internal.New(tt); err == nil {
			t.Errorf("New(%q) = (%v, %v) want (nil, an error)", tt, got, err)
		}
	}
}

func TestValidMicroversions(t *testing.T) {
	tests := []struct {
		validInput string
		want       *internal.Microversion
	}{
		{
			validInput: "2.0",
			want:       &internal.Microversion{Major: 2, Minor: 0},
		},
		{
			validInput: "3.1",
			want:       &internal.Microversion{Major: 3, Minor: 1},
		},
	}
	for _, tt := range tests {
		if got, err := internal.New(tt.validInput); err != nil {
			t.Errorf("New(%q) = (%v, %v) want (%v, %v)", tt.validInput, got, err.Error(), tt.want, nil)
		} else {
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New(%q) = (%v, %v) want (%v, %v)", tt.validInput, got, err, tt.want, nil)
			}
		}
	}
}

func TestLessThan(t *testing.T) {
	tests := []struct {
		v1   *internal.Microversion
		v2   *internal.Microversion
		want bool
	}{
		{
			v1:   &internal.Microversion{Major: 0, Minor: 1},
			v2:   &internal.Microversion{Major: 1, Minor: 0},
			want: true,
		},
		{
			v1:   &internal.Microversion{Major: 2, Minor: 1},
			v2:   &internal.Microversion{Major: 1, Minor: 0},
			want: false,
		},
		{
			v1:   &internal.Microversion{Major: 2, Minor: 1},
			v2:   &internal.Microversion{Major: 2, Minor: 0},
			want: false,
		},
		{
			v1:   &internal.Microversion{Major: 2, Minor: 1},
			v2:   &internal.Microversion{Major: 2, Minor: 1},
			want: false,
		},
	}
	for _, tt := range tests {
		if got := tt.v1.LessThan(tt.v2); got != tt.want {
			t.Errorf("LessThan(%v, %v) = (%v) want (%v)", tt.v1, tt.v2, got, tt.want)
		}
	}
}

func TestString(t *testing.T) {
	want := "3.7"
	test, err := internal.New(want)
	if err != nil {
		t.Errorf("internal error: (%v)", err.Error())
	}
	if got := test.String(); got != want {
		t.Errorf("(%v).String() = %q want %q", test, got, want)
	}
}
