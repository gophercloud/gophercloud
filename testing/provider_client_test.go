package testing

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestAuthenticatedHeaders(t *testing.T) {
	p := &gophercloud.ProviderClient{
		TokenID: "1234",
	}
	expected := map[string]string{"X-Auth-Token": "1234"}
	actual := p.AuthenticatedHeaders()
	th.CheckDeepEquals(t, expected, actual)
}

func TestUserAgent(t *testing.T) {
	p := &gophercloud.ProviderClient{}

	p.UserAgent.Prepend("custom-user-agent/2.4.0")
	expected := "custom-user-agent/2.4.0 " + gophercloud.DefaultUserAgent
	actual := p.UserAgent.Join()
	th.CheckEquals(t, expected, actual)

	p.UserAgent.Prepend("another-custom-user-agent/0.3.0", "a-third-ua/5.9.0")
	expected = "another-custom-user-agent/0.3.0 a-third-ua/5.9.0 custom-user-agent/2.4.0 " + gophercloud.DefaultUserAgent
	actual = p.UserAgent.Join()
	th.CheckEquals(t, expected, actual)

	p.UserAgent = gophercloud.UserAgent{}
	expected = gophercloud.DefaultUserAgent
	actual = p.UserAgent.Join()
	th.CheckEquals(t, expected, actual)
}

func TestConcurrentReauth(t *testing.T) {
	var info = struct {
		numreauths  int
		failedAuths int
		mut         *sync.RWMutex
	}{
		0,
		0,
		new(sync.RWMutex),
	}

	numconc := 20

	prereauthTok := client.TokenID
	postreauthTok := "12345678"

	p := new(gophercloud.ProviderClient)
	p.UseTokenLock()
	p.SetToken(prereauthTok)
	p.ReauthFunc = func(_ context.Context) error {
		p.SetThrowaway(true)
		time.Sleep(1 * time.Second)
		p.AuthenticatedHeaders()
		info.mut.Lock()
		info.numreauths++
		info.mut.Unlock()
		p.TokenID = postreauthTok
		p.SetThrowaway(false)
		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth-Token") != postreauthTok {
			w.WriteHeader(http.StatusUnauthorized)
			info.mut.Lock()
			info.failedAuths++
			info.mut.Unlock()
			return
		}
		info.mut.RLock()
		hasReauthed := info.numreauths != 0
		info.mut.RUnlock()

		if hasReauthed {
			th.CheckEquals(t, p.Token(), postreauthTok)
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{}`)
	})

	wg := new(sync.WaitGroup)
	reqopts := new(gophercloud.RequestOpts)
	reqopts.KeepResponseBody = true
	reqopts.MoreHeaders = map[string]string{
		"X-Auth-Token": prereauthTok,
	}

	for i := 0; i < numconc; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := p.Request(context.TODO(), "GET", fmt.Sprintf("%s/route", th.Endpoint()), reqopts)
			th.CheckNoErr(t, err)
			if resp == nil {
				t.Errorf("got a nil response")
				return
			}
			if resp.Body == nil {
				t.Errorf("response body was nil")
				return
			}
			defer resp.Body.Close()
			actual, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("error reading response body: %s", err)
				return
			}
			th.CheckByteArrayEquals(t, []byte(`{}`), actual)
		}()
	}

	wg.Wait()

	th.AssertEquals(t, 1, info.numreauths)
}

func TestReauthEndLoop(t *testing.T) {
	var info = struct {
		reauthAttempts   int
		maxReauthReached bool
		mut              *sync.RWMutex
	}{
		0,
		false,
		new(sync.RWMutex),
	}

	numconc := 20
	mut := new(sync.RWMutex)

	p := new(gophercloud.ProviderClient)
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.ReauthFunc = func(_ context.Context) error {
		info.mut.Lock()
		defer info.mut.Unlock()

		if info.reauthAttempts > 5 {
			info.maxReauthReached = true
			return fmt.Errorf("Max reauthentication attempts reached")
		}
		p.SetThrowaway(true)
		p.AuthenticatedHeaders()
		p.SetThrowaway(false)
		info.reauthAttempts++

		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		// route always return 401
		w.WriteHeader(http.StatusUnauthorized)
	})

	reqopts := new(gophercloud.RequestOpts)

	// counters for the upcoming errors
	errAfter := 0
	errUnable := 0

	wg := new(sync.WaitGroup)
	for i := 0; i < numconc; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Request(context.TODO(), "GET", fmt.Sprintf("%s/route", th.Endpoint()), reqopts)

			mut.Lock()
			defer mut.Unlock()

			// ErrErrorAfter... will happen after a successful reauthentication,
			// but the service still responds with a 401.
			if _, ok := err.(*gophercloud.ErrErrorAfterReauthentication); ok {
				errAfter++
			}

			// ErrErrorUnable... will happen when the custom reauth func reports
			// an error.
			if _, ok := err.(*gophercloud.ErrUnableToReauthenticate); ok {
				errUnable++
			}
		}()
	}

	wg.Wait()
	th.AssertEquals(t, info.reauthAttempts, 6)
	th.AssertEquals(t, info.maxReauthReached, true)
	th.AssertEquals(t, errAfter > 1, true)
	th.AssertEquals(t, errUnable < 20, true)
}

func TestRequestThatCameDuringReauthWaitsUntilItIsCompleted(t *testing.T) {
	var info = struct {
		numreauths  int
		failedAuths int
		reauthCh    chan struct{}
		mut         *sync.RWMutex
	}{
		0,
		0,
		make(chan struct{}),
		new(sync.RWMutex),
	}

	numconc := 20

	prereauthTok := client.TokenID
	postreauthTok := "12345678"

	p := new(gophercloud.ProviderClient)
	p.UseTokenLock()
	p.SetToken(prereauthTok)
	p.ReauthFunc = func(_ context.Context) error {
		info.mut.RLock()
		if info.numreauths == 0 {
			info.mut.RUnlock()
			close(info.reauthCh)
			time.Sleep(1 * time.Second)
		} else {
			info.mut.RUnlock()
		}
		p.SetThrowaway(true)
		p.AuthenticatedHeaders()
		info.mut.Lock()
		info.numreauths++
		info.mut.Unlock()
		p.TokenID = postreauthTok
		p.SetThrowaway(false)
		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth-Token") != postreauthTok {
			info.mut.Lock()
			info.failedAuths++
			info.mut.Unlock()
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		info.mut.RLock()
		hasReauthed := info.numreauths != 0
		info.mut.RUnlock()

		if hasReauthed {
			th.CheckEquals(t, p.Token(), postreauthTok)
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{}`)
	})

	wg := new(sync.WaitGroup)
	reqopts := new(gophercloud.RequestOpts)
	reqopts.KeepResponseBody = true
	reqopts.MoreHeaders = map[string]string{
		"X-Auth-Token": prereauthTok,
	}

	for i := 0; i < numconc; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i != 0 {
				<-info.reauthCh
			}
			resp, err := p.Request(context.TODO(), "GET", fmt.Sprintf("%s/route", th.Endpoint()), reqopts)
			th.CheckNoErr(t, err)
			if resp == nil {
				t.Errorf("got a nil response")
				return
			}
			if resp.Body == nil {
				t.Errorf("response body was nil")
				return
			}
			defer resp.Body.Close()
			actual, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("error reading response body: %s", err)
				return
			}
			th.CheckByteArrayEquals(t, []byte(`{}`), actual)
		}(i)
	}

	wg.Wait()

	th.AssertEquals(t, 1, info.numreauths)
	th.AssertEquals(t, 1, info.failedAuths)
}

func TestRequestReauthsAtMostOnce(t *testing.T) {
	// There was an issue where Gophercloud would go into an infinite
	// reauthentication loop with buggy services that send 401 even for fresh
	// tokens. This test simulates such a service and checks that a call to
	// ProviderClient.Request() will not try to reauthenticate more than once.

	reauthCounter := 0
	var reauthCounterMutex sync.Mutex

	p := new(gophercloud.ProviderClient)
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.ReauthFunc = func(_ context.Context) error {
		reauthCounterMutex.Lock()
		reauthCounter++
		reauthCounterMutex.Unlock()
		//The actual token value does not matter, the endpoint does not check it.
		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	requestCounter := 0
	var requestCounterMutex sync.Mutex

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		requestCounterMutex.Lock()
		requestCounter++
		//avoid infinite loop
		if requestCounter == 10 {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		requestCounterMutex.Unlock()

		//always reply 401, even immediately after reauthenticate
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	// The expected error message indicates that we reauthenticated once (that's
	// the part before the colon), but when encountering another 401 response, we
	// did not attempt reauthentication again and just passed that 401 response to
	// the caller as ErrDefault401.
	_, err := p.Request(context.TODO(), "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	expectedErrorRx := regexp.MustCompile(`^Successfully re-authenticated, but got error executing request: Expected HTTP response code \[200\] when accessing \[GET http://[^/]*//route\], but got 401 instead: unauthorized$`)
	if !expectedErrorRx.MatchString(err.Error()) {
		t.Errorf("expected error that looks like %q, but got %q", expectedErrorRx.String(), err.Error())
	}
}

func TestRequestWithContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	p := &gophercloud.ProviderClient{}

	res, err := p.Request(ctx, "GET", ts.URL, &gophercloud.RequestOpts{KeepResponseBody: true})
	th.AssertNoErr(t, err)
	_, err = io.ReadAll(res.Body)
	th.AssertNoErr(t, err)
	err = res.Body.Close()
	th.AssertNoErr(t, err)

	cancel()
	_, err = p.Request(ctx, "GET", ts.URL, &gophercloud.RequestOpts{})
	if err == nil {
		t.Fatal("expecting error, got nil")
	}
	if !strings.Contains(err.Error(), ctx.Err().Error()) {
		t.Fatalf("expecting error to contain: %q, got %q", ctx.Err().Error(), err.Error())
	}
}

func TestRequestConnectionReuse(t *testing.T) {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))

	// an amount of iterations
	var iter = 10000
	// connections tracks an amount of connections made
	var connections int64

	ts.Config.ConnState = func(_ net.Conn, s http.ConnState) {
		// track an amount of connections
		if s == http.StateNew {
			atomic.AddInt64(&connections, 1)
		}
	}
	ts.Start()
	defer ts.Close()

	p := &gophercloud.ProviderClient{}
	for i := 0; i < iter; i++ {
		_, err := p.Request(context.TODO(), "GET", ts.URL, &gophercloud.RequestOpts{KeepResponseBody: false})
		th.AssertNoErr(t, err)
	}

	th.AssertEquals(t, int64(1), connections)
}

func TestRequestConnectionClose(t *testing.T) {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))

	// an amount of iterations
	var iter = 10
	// connections tracks an amount of connections made
	var connections int64

	ts.Config.ConnState = func(_ net.Conn, s http.ConnState) {
		// track an amount of connections
		if s == http.StateNew {
			atomic.AddInt64(&connections, 1)
		}
	}
	ts.Start()
	defer ts.Close()

	p := &gophercloud.ProviderClient{}
	for i := 0; i < iter; i++ {
		_, err := p.Request(context.TODO(), "GET", ts.URL, &gophercloud.RequestOpts{KeepResponseBody: true})
		th.AssertNoErr(t, err)
	}

	th.AssertEquals(t, int64(iter), connections)
}

func retryBackoffTest(retryCounter *uint, t *testing.T) gophercloud.RetryBackoffFunc {
	return func(ctx context.Context, respErr *gophercloud.ErrUnexpectedResponseCode, e error, retries uint) error {
		retryAfter := respErr.ResponseHeader.Get("Retry-After")
		if retryAfter == "" {
			return e
		}

		var sleep time.Duration

		// Parse delay seconds or HTTP date
		if v, err := strconv.ParseUint(retryAfter, 10, 32); err == nil {
			sleep = time.Duration(v) * time.Second
		} else if v, err := time.Parse(http.TimeFormat, retryAfter); err == nil {
			sleep = time.Until(v)
		} else {
			return e
		}

		if ctx != nil {
			t.Logf("Context sleeping for %d milliseconds", sleep.Milliseconds())
			select {
			case <-time.After(sleep):
				t.Log("sleep is over")
			case <-ctx.Done():
				t.Log(ctx.Err())
				return e
			}
		} else {
			t.Logf("Sleeping for %d milliseconds", sleep.Milliseconds())
			time.Sleep(sleep)
			t.Log("sleep is over")
		}

		*retryCounter = *retryCounter + 1

		return nil
	}
}

func TestRequestRetry(t *testing.T) {
	var retryCounter uint

	p := &gophercloud.ProviderClient{}
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.MaxBackoffRetries = 3

	p.RetryBackoffFunc = retryBackoffTest(&retryCounter, t)

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")

		//always reply 429
		http.Error(w, "retry later", http.StatusTooManyRequests)
	})

	_, err := p.Request(context.TODO(), "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	if err == nil {
		t.Fatal("expecting error, got nil")
	}
	th.AssertEquals(t, retryCounter, p.MaxBackoffRetries)
}

func TestRequestRetryHTTPDate(t *testing.T) {
	var retryCounter uint

	p := &gophercloud.ProviderClient{}
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.MaxBackoffRetries = 3

	p.RetryBackoffFunc = retryBackoffTest(&retryCounter, t)

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", time.Now().Add(1*time.Second).UTC().Format(http.TimeFormat))

		//always reply 429
		http.Error(w, "retry later", http.StatusTooManyRequests)
	})

	_, err := p.Request(context.TODO(), "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	if err == nil {
		t.Fatal("expecting error, got nil")
	}
	th.AssertEquals(t, retryCounter, p.MaxBackoffRetries)
}

func TestRequestRetryError(t *testing.T) {
	var retryCounter uint

	p := &gophercloud.ProviderClient{}
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.MaxBackoffRetries = 3

	p.RetryBackoffFunc = retryBackoffTest(&retryCounter, t)

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "foo bar")

		//always reply 429
		http.Error(w, "retry later", http.StatusTooManyRequests)
	})

	_, err := p.Request(context.TODO(), "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	if err == nil {
		t.Fatal("expecting error, got nil")
	}
	th.AssertEquals(t, retryCounter, uint(0))
}

func TestRequestRetrySuccess(t *testing.T) {
	var retryCounter uint

	p := &gophercloud.ProviderClient{}
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.MaxBackoffRetries = 3

	p.RetryBackoffFunc = retryBackoffTest(&retryCounter, t)

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		//always reply 200
		http.Error(w, "retry later", http.StatusOK)
	})

	_, err := p.Request(context.TODO(), "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	if err != nil {
		t.Fatal(err)
	}
	th.AssertEquals(t, retryCounter, uint(0))
}

func TestRequestRetryContext(t *testing.T) {
	var retryCounter uint

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sleep := 2.5 * 1000 * time.Millisecond
		time.Sleep(sleep)
		cancel()
	}()

	p := &gophercloud.ProviderClient{}
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.MaxBackoffRetries = 3

	p.RetryBackoffFunc = retryBackoffTest(&retryCounter, t)

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")

		//always reply 429
		http.Error(w, "retry later", http.StatusTooManyRequests)
	})

	_, err := p.Request(ctx, "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	if err == nil {
		t.Fatal("expecting error, got nil")
	}
	t.Logf("retryCounter: %d, p.MaxBackoffRetries: %d", retryCounter, p.MaxBackoffRetries-1)
	th.AssertEquals(t, retryCounter, p.MaxBackoffRetries-1)
}

func TestRequestGeneralRetry(t *testing.T) {
	p := &gophercloud.ProviderClient{}
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.RetryFunc = func(context context.Context, method, url string, options *gophercloud.RequestOpts, err error, failCount uint) error {
		if failCount >= 5 {
			return err
		}
		return nil
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	count := 0
	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if count < 3 {
			http.Error(w, "bad gateway", http.StatusBadGateway)
			count += 1
		} else {
			fmt.Fprintln(w, "OK")
		}
	})

	_, err := p.Request(context.TODO(), "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	if err != nil {
		t.Fatal("expecting nil, got err")
	}
	th.AssertEquals(t, 3, count)
}

func TestRequestGeneralRetryAbort(t *testing.T) {
	p := &gophercloud.ProviderClient{}
	p.UseTokenLock()
	p.SetToken(client.TokenID)
	p.RetryFunc = func(context context.Context, method, url string, options *gophercloud.RequestOpts, err error, failCount uint) error {
		return err
	}

	th.SetupHTTP()
	defer th.TeardownHTTP()

	count := 0
	th.Mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if count < 3 {
			http.Error(w, "bad gateway", http.StatusBadGateway)
			count += 1
		} else {
			fmt.Fprintln(w, "OK")
		}
	})

	_, err := p.Request(context.TODO(), "GET", th.Endpoint()+"/route", &gophercloud.RequestOpts{})
	if err == nil {
		t.Fatal("expecting err, got nil")
	}
	th.AssertEquals(t, 1, count)
}

func TestRequestWrongOkCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
		// Returns 200 OK
	}))
	defer ts.Close()

	p := &gophercloud.ProviderClient{}

	_, err := p.Request(context.TODO(), "DELETE", ts.URL, &gophercloud.RequestOpts{})
	th.AssertErr(t, err)
	if urErr, ok := err.(gophercloud.ErrUnexpectedResponseCode); ok {
		// DELETE expects a 202 or 204 by default
		// Make sure returned error contains the expected OK codes
		th.AssertDeepEquals(t, []int{202, 204}, urErr.Expected)
	} else {
		t.Fatalf("expected error type gophercloud.ErrUnexpectedResponseCode but got %T", err)
	}
}
