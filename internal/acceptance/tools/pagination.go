package tools

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CheckPagination checks that fresh EachPage and AllPages pagers retrieve every
// expected ID exactly once. Each nonempty page must contain one item. Unrelated
// resources are allowed, but duplicate or empty IDs are rejected throughout the
// collection. At least two distinct expected IDs are required.
//
// newPager must use the public list constructor with a page size of one.
// extractIDs must use the resource's public extraction function.
func CheckPagination(t *testing.T, ctx context.Context, newPager func() pagination.Pager, extractIDs func(pagination.Page) ([]string, error), expectedIDs []string) {
	t.Helper()
	for _, allPages := range []bool{false, true} {
		name := "EachPage"
		if allPages {
			name = "AllPages"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()
			if err := checkPagination(ctx, newPager(), extractIDs, expectedIDs, allPages); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func checkPagination(ctx context.Context, pager pagination.Pager, extractIDs func(pagination.Page) ([]string, error), expectedIDs []string, allPages bool) error {
	expected := make(map[string]bool, len(expectedIDs))
	for _, id := range expectedIDs {
		if id == "" || expected[id] {
			return fmt.Errorf("expected IDs must be nonempty and distinct: %q", expectedIDs)
		}
		expected[id] = true
	}
	if len(expected) < 2 {
		return fmt.Errorf("pagination requires at least two expected IDs")
	}

	seen := make(map[string]bool)
	collect := func(page pagination.Page) error {
		ids, err := extractIDs(page)
		if err != nil {
			return err
		}
		if !allPages && len(ids) != 1 {
			return fmt.Errorf("expected page size 1, got %d", len(ids))
		}
		for _, id := range ids {
			if id == "" {
				return fmt.Errorf("collection contains an empty ID")
			}
			if seen[id] {
				return fmt.Errorf("collection repeated ID %q", id)
			}
			seen[id] = true
		}
		return nil
	}

	if allPages {
		page, err := pager.AllPages(ctx)
		if err != nil {
			return err
		}
		if err := collect(page); err != nil {
			return err
		}
	} else {
		err := pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
			return true, collect(page)
		})
		if err != nil {
			return err
		}
	}
	for _, id := range expectedIDs {
		if !seen[id] {
			return fmt.Errorf("collection omitted expected ID %q", id)
		}
	}
	return nil
}
