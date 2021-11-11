package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
)

func TestContext(t *testing.T) {
	t.Run("cancellation", func(t *testing.T) {
		t.Run("cancel is idempotent", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			cancelMergedCtx()
			cancel1()
			cancel2()
			cancelMergedCtx()

			time.Sleep(time.Millisecond)

			cancelMergedCtx()

			select {
			case <-mergedCtx.Done():
			default:
				t.Errorf("expected mergedCtx to have been cancelled")
			}
		})

		t.Run("nothing cancelled", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			time.Sleep(time.Millisecond)

			select {
			case <-mergedCtx.Done():
				t.Errorf("expected mergedCtx to not have been cancelled")
			default:
			}
		})

		t.Run("cancel ctx1", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			cancel1()
			time.Sleep(time.Millisecond)

			select {
			case <-mergedCtx.Done():
			default:
				t.Errorf("expected mergedCtx to have been cancelled")
			}
		})

		t.Run("cancel ctx2", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			cancel2()
			time.Sleep(time.Millisecond)

			select {
			case <-mergedCtx.Done():
			default:
				t.Errorf("expected mergedCtx to have been cancelled")
			}
		})

		t.Run("cancel mergeCtx", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			cancelMergedCtx()
			time.Sleep(time.Millisecond)

			select {
			case <-mergedCtx.Done():
			default:
				t.Errorf("expected mergedCtx to have been cancelled")
			}
		})
	})

	t.Run("deadline", func(t *testing.T) {
		t.Run("no deadline", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			if d, ok := mergedCtx.Deadline(); ok {
				t.Errorf("expected mergedCtx to not have deadline, found %q", d)
			}
		})

		t.Run("ctx1 has deadline", func(t *testing.T) {
			t1 := time.Date(1955, time.November, 5, 1, 21, 0, 0, time.Local)

			ctx1, cancel1 := context.WithDeadline(context.Background(), t1)
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			if d, ok := mergedCtx.Deadline(); ok {
				if d != t1 {
					t.Errorf("expected mergedCtx to have deadline %q, found %q", t1, d)
				}
			} else {
				t.Errorf("expected mergedCtx to have deadline, found %q", d)
			}
		})

		t.Run("ctx2 has deadline", func(t *testing.T) {
			t1 := time.Date(1955, time.November, 5, 1, 21, 0, 0, time.Local)

			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithDeadline(context.Background(), t1)
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			if d, ok := mergedCtx.Deadline(); ok {
				if d != t1 {
					t.Errorf("expected mergedCtx to have deadline %q, found %q", t1, d)
				}
			} else {
				t.Errorf("expected mergedCtx to have deadline, found %q", d)
			}
		})

		t.Run("two deadlines, closer is 1", func(t *testing.T) {
			t1 := time.Date(1955, time.November, 5, 1, 21, 0, 0, time.Local)
			t2 := time.Date(2015, time.November, 5, 1, 21, 0, 0, time.Local)

			ctx1, cancel1 := context.WithDeadline(context.Background(), t1)
			defer cancel1()
			ctx2, cancel2 := context.WithDeadline(context.Background(), t2)
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			if d, ok := mergedCtx.Deadline(); ok {
				if d != t1 {
					t.Errorf("expected mergedCtx to have deadline %q, found %q", t1, d)
				}
			} else {
				t.Errorf("expected mergedCtx to have deadline, found %q", d)
			}
		})

		t.Run("two deadlines, closer is 2", func(t *testing.T) {
			t1 := time.Date(1955, time.November, 5, 1, 21, 0, 0, time.Local)
			t2 := time.Date(2015, time.November, 5, 1, 21, 0, 0, time.Local)

			ctx1, cancel1 := context.WithDeadline(context.Background(), t2)
			defer cancel1()
			ctx2, cancel2 := context.WithDeadline(context.Background(), t1)
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			if d, ok := mergedCtx.Deadline(); ok {
				if d != t1 {
					t.Errorf("expected mergedCtx to have deadline %q, found %q", t1, d)
				}
			} else {
				t.Errorf("expected mergedCtx to have deadline, found %q", d)
			}
		})
	})

	t.Run("errors", func(t *testing.T) {
		t.Run("consistently returns the first error, ctx1", func(t *testing.T) {
			ctx1, cancel1 := context.WithDeadline(context.Background(), time.Now())
			defer cancel1()
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			timer := time.NewTimer(time.Millisecond)
			defer timer.Stop()

			select {
			case <-mergedCtx.Done():
			case <-timer.C:
				t.Fatal("mergedCtx did not stop")
			}

			if want, have := context.DeadlineExceeded, mergedCtx.Err(); want != have {
				t.Errorf("expected error %q, found %q", want, have)
			}

			cancel2()
			if want, have := context.DeadlineExceeded, mergedCtx.Err(); want != have {
				t.Errorf("expected error to stay %q after the second context is done, found %q", want, have)
			}
		})

		t.Run("consistently returns the first error, ctx2", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			ctx2, cancel2 := context.WithDeadline(context.Background(), time.Now())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, ctx2)
			defer cancelMergedCtx()

			timer := time.NewTimer(time.Millisecond)
			defer timer.Stop()

			select {
			case <-mergedCtx.Done():
			case <-timer.C:
				t.Fatal("mergedCtx did not stop")
			}

			if want, have := context.DeadlineExceeded, mergedCtx.Err(); want != have {
				t.Errorf("expected error %q, found %q", want, have)
			}

			cancel1()
			if want, have := context.DeadlineExceeded, mergedCtx.Err(); want != have {
				t.Errorf("expected error to stay %q after the second context is done, found %q", want, have)
			}
		})
	})

	t.Run("merge", func(t *testing.T) {
		t.Run("ctx2 nil", func(t *testing.T) {
			ctx2, cancel2 := context.WithCancel(context.Background())
			defer cancel2()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(nil, ctx2)
			defer cancelMergedCtx()

			cancel2()
			time.Sleep(time.Millisecond)

			select {
			case <-mergedCtx.Done():
			default:
				t.Errorf("expected mergedCtx to have been cancelled")
			}
		})

		t.Run("ctx1 nil", func(t *testing.T) {
			ctx1, cancel1 := context.WithCancel(context.Background())
			defer cancel1()
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctx1, nil)
			defer cancelMergedCtx()

			cancel1()
			time.Sleep(time.Millisecond)

			select {
			case <-mergedCtx.Done():
			default:
				t.Errorf("expected mergedCtx to have been cancelled")
			}
		})

		t.Run("nil nil", func(t *testing.T) {
			mergedCtx, cancelMergedCtx := gophercloud.MergeContext(nil, nil)
			defer cancelMergedCtx()

			cancelMergedCtx()

			select {
			case <-mergedCtx.Done():
			default:
				t.Errorf("expected mergedCtx to have been cancelled")
			}
		})
	})

	t.Run("value", func(t *testing.T) {
		type key string
		type value string
		var (
			key1 key = "this is a key"
			key2 key = "this is another key"
			key3 key = "this key is for testing precedence"

			value1           value = "this is a value"
			value2           value = "this is another value"
			value3a, value3b value = "this value should be returned", "this value will be overwritten"
		)

		ctxA := context.WithValue(context.Background(), key1, value1)
		ctxA = context.WithValue(ctxA, key3, value3a)

		ctxB := context.WithValue(context.Background(), key2, value2)
		ctxB = context.WithValue(ctxB, key3, value3b)

		mergedCtx, cancelMergedCtx := gophercloud.MergeContext(ctxA, ctxB)
		defer cancelMergedCtx()

		{
			var key, want = key1, value1
			v := mergedCtx.Value(key)
			if have, ok := v.(value); ok {
				if want != have {
					t.Errorf("expected value from child context %q, found %q", want, have)
				}
			} else {
				t.Errorf("expected value from child context %q, found %q", want, v)
			}
		}

		{
			var key, want = key2, value2
			v := mergedCtx.Value(key)
			if have, ok := v.(value); ok {
				if want != have {
					t.Errorf("expected value from parent context %q, found %q", want, have)
				}
			} else {
				t.Errorf("expected value from parent context %q, found %q", want, v)
			}
		}

		{
			var key, want = key3, value3a

			v := mergedCtx.Value(key)
			if have, ok := v.(value); ok {
				if want != have {
					t.Errorf("expected value from child context %q, found %q", want, have)
				}
			} else {
				t.Errorf("expected value from child context %q, found %q", want, v)
			}
		}
	})
}
