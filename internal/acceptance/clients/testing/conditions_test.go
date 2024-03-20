//go:build acceptance

package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
)

func TestIsCurrentAbove(t *testing.T) {
	cases := []struct {
		Current string
		Release string
		Result  bool
	}{
		{Current: "master", Release: "zed", Result: true},
		{Current: "master", Release: "2023.1", Result: true},
		{Current: "master", Release: "master", Result: false},
		{Current: "zed", Release: "master", Result: false},
		{Current: "zed", Release: "yoga", Result: true},
		{Current: "zed", Release: "2023.1", Result: false},
		{Current: "2023.1", Release: "2023.1", Result: false},
		{Current: "2023.2", Release: "stable/2023.1", Result: true},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%s above %s", tt.Current, tt.Release), func(t *testing.T) {
			os.Setenv("OS_BRANCH", tt.Current)
			got := clients.IsCurrentAbove(t, tt.Release)
			if got != tt.Result {
				t.Errorf("got %v want %v", got, tt.Result)
			}
		})

	}
}

func TestIsCurrentBelow(t *testing.T) {
	cases := []struct {
		Current string
		Release string
		Result  bool
	}{
		{Current: "master", Release: "zed", Result: false},
		{Current: "master", Release: "2023.1", Result: false},
		{Current: "master", Release: "master", Result: false},
		{Current: "zed", Release: "master", Result: true},
		{Current: "zed", Release: "yoga", Result: false},
		{Current: "zed", Release: "2023.1", Result: true},
		{Current: "2023.1", Release: "2023.1", Result: false},
		{Current: "2023.2", Release: "stable/2023.1", Result: false},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%s below %s", tt.Current, tt.Release), func(t *testing.T) {
			os.Setenv("OS_BRANCH", tt.Current)
			got := clients.IsCurrentBelow(t, tt.Release)
			if got != tt.Result {
				t.Errorf("got %v want %v", got, tt.Result)
			}
		})

	}
}
