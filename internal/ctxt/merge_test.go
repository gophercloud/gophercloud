package ctxt_test

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/ctxt"
)

func TestMerge(t *testing.T) {
	t.Run("returns values from both parents", func(t *testing.T) {
		ctx1 := context.WithValue(context.Background(),
			"key1", "value1")

		ctx2 := context.WithValue(context.WithValue(context.Background(),
			"key1", "this value should be overridden"),
			"key2", "value2")

		ctx, cancel := ctxt.Merge(ctx1, ctx2)
		defer cancel()

		if v1 := ctx.Value("key1"); v1 != nil {
			if s1, ok := v1.(string); ok {
				if s1 != "value1" {
					t.Errorf("found value for key1 %q, expected %q", s1, "value1")
				}
			} else {
				t.Errorf("key1 is not the expected type string")
			}
		} else {
			t.Errorf("key1 returned nil")
		}

		if v2 := ctx.Value("key2"); v2 != nil {
			if s2, ok := v2.(string); ok {
				if s2 != "value2" {
					t.Errorf("found value for key2 %q, expected %q", s2, "value2")
				}
			} else {
				t.Errorf("key2 is not the expected type string")
			}
		} else {
			t.Errorf("key2 returned nil")
		}
	})

	t.Run("first parent cancels", func(t *testing.T) {
		ctx1, cancel1 := context.WithCancel(context.Background())
		ctx2, cancel2 := context.WithCancel(context.Background())
		defer cancel2()

		ctx, cancel := ctxt.Merge(ctx1, ctx2)
		defer cancel()

		if err := ctx.Err(); err != nil {
			t.Errorf("context unexpectedly done: %v", err)
		}

		cancel1()
		time.Sleep(1 * time.Millisecond)
		if err := ctx.Err(); err == nil {
			t.Errorf("context not done despite parent1 cancelled")
		}
	})

	t.Run("second parent cancels", func(t *testing.T) {
		ctx1, cancel1 := context.WithCancel(context.Background())
		ctx2, cancel2 := context.WithCancel(context.Background())
		defer cancel1()

		ctx, cancel := ctxt.Merge(ctx1, ctx2)
		defer cancel()

		if err := ctx.Err(); err != nil {
			t.Errorf("context unexpectedly done: %v", err)
		}

		cancel2()
		time.Sleep(1 * time.Millisecond)
		if err := ctx.Err(); err == nil {
			t.Errorf("context not done despite parent2 cancelled")
		}
	})

	t.Run("inherits deadline from first parent", func(t *testing.T) {
		now := time.Now()
		t1 := now.Add(time.Hour)
		t2 := t1.Add(time.Second)

		ctx1, cancel1 := context.WithDeadline(context.Background(), t1)
		ctx2, cancel2 := context.WithDeadline(context.Background(), t2)
		defer cancel1()
		defer cancel2()

		ctx, cancel := ctxt.Merge(ctx1, ctx2)
		defer cancel()

		if err := ctx.Err(); err != nil {
			t.Errorf("context unexpectedly done: %v", err)
		}

		if deadline, ok := ctx.Deadline(); ok {
			if deadline != t1 {
				t.Errorf("expected deadline to be %v, found %v", t1, deadline)
			}
		}
	})

	t.Run("inherits deadline from second parent", func(t *testing.T) {
		now := time.Now()
		t2 := now.Add(time.Hour)
		t1 := t2.Add(time.Second)

		ctx1, cancel1 := context.WithDeadline(context.Background(), t1)
		ctx2, cancel2 := context.WithDeadline(context.Background(), t2)
		defer cancel1()
		defer cancel2()

		ctx, cancel := ctxt.Merge(ctx1, ctx2)
		defer cancel()

		if err := ctx.Err(); err != nil {
			t.Errorf("context unexpectedly done: %v", err)
		}

		if deadline, ok := ctx.Deadline(); ok {
			if deadline != t2 {
				t.Errorf("expected deadline to be %v, found %v", t2, deadline)
			}
		}
	})
}
