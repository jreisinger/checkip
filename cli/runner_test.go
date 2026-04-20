package cli

import (
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jreisinger/checkip/check"
)

func TestRunnerCachesRepeatedIPs(t *testing.T) {
	var cachedRuns int32
	definitions := []check.Definition{
		{
			Name: "cached",
			Run: func(net.IP) (check.Check, error) {
				atomic.AddInt32(&cachedRuns, 1)
				return check.Check{Description: "cached", Type: check.Info}, nil
			},
		},
	}

	runner := NewRunner(definitions)
	ipaddr := net.ParseIP("1.1.1.1")

	for range 3 {
		checks, errors := runner.Run(ipaddr)
		if len(errors) != 0 {
			t.Fatalf("got %d errors, want 0", len(errors))
		}
		if len(checks) != 1 {
			t.Fatalf("got %d checks, want 1", len(checks))
		}
	}

	if got := atomic.LoadInt32(&cachedRuns); got != 1 {
		t.Fatalf("cached check ran %d times, want 1", got)
	}
}

func TestRunnerCachesOnlyCacheableChecks(t *testing.T) {
	var cachedRuns int32
	var liveRuns int32

	definitions := []check.Definition{
		{
			Name: "cached",
			Run: func(net.IP) (check.Check, error) {
				atomic.AddInt32(&cachedRuns, 1)
				return check.Check{Description: "cached", Type: check.Info}, nil
			},
		},
		{
			Name:  "live",
			Cache: check.CacheNone,
			Run: func(net.IP) (check.Check, error) {
				atomic.AddInt32(&liveRuns, 1)
				return check.Check{Description: "live", Type: check.Info}, nil
			},
		},
	}

	runner := NewRunner(definitions)
	ipaddr := net.ParseIP("1.1.1.1")

	for range 2 {
		checks, errors := runner.Run(ipaddr)
		if len(errors) != 0 {
			t.Fatalf("got %d errors, want 0", len(errors))
		}
		if len(checks) != 2 {
			t.Fatalf("got %d checks, want 2", len(checks))
		}
	}

	if got := atomic.LoadInt32(&cachedRuns); got != 1 {
		t.Fatalf("cached check ran %d times, want 1", got)
	}
	if got := atomic.LoadInt32(&liveRuns); got != 2 {
		t.Fatalf("live check ran %d times, want 2", got)
	}
}

func TestRunnerCollapsesConcurrentRepeatedIPs(t *testing.T) {
	var runs int32

	definitions := []check.Definition{
		{
			Name: "cached",
			Run: func(net.IP) (check.Check, error) {
				atomic.AddInt32(&runs, 1)
				time.Sleep(20 * time.Millisecond)
				return check.Check{Description: "cached", Type: check.Info}, nil
			},
		},
	}

	runner := NewRunner(definitions)
	ipaddr := net.ParseIP("1.1.1.1")

	start := make(chan struct{})
	var wg sync.WaitGroup
	errors := make(chan error, 4)

	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			checks, runErrors := runner.Run(ipaddr)
			if len(runErrors) != 0 {
				errors <- runErrors[0]
				return
			}
			if len(checks) != 1 {
				errors <- &testError{msg: "unexpected check count"}
				return
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}

	if got := atomic.LoadInt32(&runs); got != 1 {
		t.Fatalf("check ran %d times, want 1", got)
	}
}

// Sorting the returned slice must not affect subsequent cached results.
func TestRunnerReturnsIndependentCachedSlices(t *testing.T) {
	definitions := []check.Definition{
		{
			Name: "zeta",
			Run: func(net.IP) (check.Check, error) {
				return check.Check{Description: "zeta", Type: check.Info}, nil
			},
		},
		{
			Name: "alpha",
			Run: func(net.IP) (check.Check, error) {
				return check.Check{Description: "alpha", Type: check.Info}, nil
			},
		},
	}

	runner := NewRunner(definitions)
	ipaddr := net.ParseIP("1.1.1.1")

	checks, errors := runner.Run(ipaddr)
	if len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	checks.SortByName()

	cachedChecks, errors := runner.Run(ipaddr)
	if len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	if cachedChecks[0].Description != "zeta" {
		t.Fatalf("first cached check = %q, want %q", cachedChecks[0].Description, "zeta")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
