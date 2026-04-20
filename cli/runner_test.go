package cli

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jreisinger/checkip/check"
)

func TestRunnerPersistentCacheHitAcrossRunners(t *testing.T) {
	var runs int32
	definition := testDefinition("canonical", time.Hour, &runs)
	cacheDir := t.TempDir()
	ipaddr := net.ParseIP("1.1.1.1")
	now := time.Date(2026, time.April, 20, 12, 0, 0, 0, time.UTC)

	runner1 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now },
	})
	checks, errors := runner1.Run(ipaddr)
	if len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	if len(checks) != 1 {
		t.Fatalf("got %d checks, want 1", len(checks))
	}

	runner2 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now.Add(30 * time.Minute) },
	})
	cachedChecks, errors := runner2.Run(ipaddr)
	if len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	if len(cachedChecks) != 1 {
		t.Fatalf("got %d checks, want 1", len(cachedChecks))
	}
	if got := atomic.LoadInt32(&runs); got != 1 {
		t.Fatalf("check ran %d times, want 1", got)
	}
	if cachedChecks[0].Description != definition.Name {
		t.Fatalf("description = %q, want %q", cachedChecks[0].Description, definition.Name)
	}

	info, ok := cachedChecks[0].IpAddrInfo.(*testInfo)
	if !ok {
		t.Fatalf("cached info type = %T, want *testInfo", cachedChecks[0].IpAddrInfo)
	}
	if info.Value != "value-1" {
		t.Fatalf("cached info value = %q, want %q", info.Value, "value-1")
	}

	out := captureStdout(t, func() {
		Checks(cachedChecks).PrintJSON(ipaddr)
	})
	if !strings.Contains(out, `"ipAddrInfo":{"value":"value-1"}`) {
		t.Fatalf("PrintJSON output = %q, want cached ipAddrInfo JSON", out)
	}
}

func TestRunnerExpiredPersistentCacheForcesRefresh(t *testing.T) {
	var runs int32
	definition := testDefinition("refresh", time.Hour, &runs)
	cacheDir := t.TempDir()
	ipaddr := net.ParseIP("1.1.1.1")
	now := time.Date(2026, time.April, 20, 12, 0, 0, 0, time.UTC)

	runner1 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now },
	})
	if _, errors := runner1.Run(ipaddr); len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}

	runner2 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now.Add(2 * time.Hour) },
	})
	refreshedChecks, errors := runner2.Run(ipaddr)
	if len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	if got := atomic.LoadInt32(&runs); got != 2 {
		t.Fatalf("check ran %d times, want 2", got)
	}

	info := refreshedChecks[0].IpAddrInfo.(*testInfo)
	if info.Value != "value-2" {
		t.Fatalf("refreshed info value = %q, want %q", info.Value, "value-2")
	}

	runner3 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now.Add(2*time.Hour + 30*time.Minute) },
	})
	if _, errors := runner3.Run(ipaddr); len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	if got := atomic.LoadInt32(&runs); got != 2 {
		t.Fatalf("check ran %d times after refresh reuse, want 2", got)
	}
}

func TestRunnerIgnoresCorruptPersistentCacheEntries(t *testing.T) {
	var runs int32
	definition := testDefinition("corrupt", time.Hour, &runs)
	cacheDir := t.TempDir()
	ipaddr := net.ParseIP("1.1.1.1")
	now := time.Date(2026, time.April, 20, 12, 0, 0, 0, time.UTC)

	cache, err := newResultCache(cacheDir, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	cachePath := cache.path(definition.Name, ipaddr.String())
	if err := os.WriteFile(cachePath, []byte("{not-json"), 0600); err != nil {
		t.Fatal(err)
	}

	runner1 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now },
	})
	if _, errors := runner1.Run(ipaddr); len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	if got := atomic.LoadInt32(&runs); got != 1 {
		t.Fatalf("check ran %d times, want 1", got)
	}

	runner2 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now.Add(30 * time.Minute) },
	})
	if _, errors := runner2.Run(ipaddr); len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}
	if got := atomic.LoadInt32(&runs); got != 1 {
		t.Fatalf("check ran %d times after rewrite, want 1", got)
	}
}

func TestRunnerNoCacheBypassesProcessAndPersistentCache(t *testing.T) {
	var runs int32
	definition := testDefinition("no-cache", time.Hour, &runs)
	cacheDir := t.TempDir()
	ipaddr := net.ParseIP("1.1.1.1")
	now := time.Date(2026, time.April, 20, 12, 0, 0, 0, time.UTC)

	seedRunner := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		CacheDir: cacheDir,
		Now:      func() time.Time { return now },
	})
	if _, errors := seedRunner.Run(ipaddr); len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}

	noCacheRunner := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{
		DisableCache: true,
		CacheDir:     cacheDir,
		Now:          func() time.Time { return now.Add(10 * time.Minute) },
	})
	for range 2 {
		if _, errors := noCacheRunner.Run(ipaddr); len(errors) != 0 {
			t.Fatalf("got %d errors, want 0", len(errors))
		}
	}

	if got := atomic.LoadInt32(&runs); got != 3 {
		t.Fatalf("check ran %d times, want 3", got)
	}
}

func TestRunnerSkipsWritingErrorsAndMissingCredentials(t *testing.T) {
	tests := []struct {
		name string
		run  func(net.IP) (check.Check, error)
	}{
		{
			name: "error",
			run: func(net.IP) (check.Check, error) {
				return check.Check{}, io.EOF
			},
		},
		{
			name: "missing credentials",
			run: func(net.IP) (check.Check, error) {
				return check.Check{
					Description:        "wrong",
					Type:               check.InfoAndIsMalicious,
					MissingCredentials: "API_KEY",
				}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheDir := t.TempDir()
			runner := NewRunnerWithOptions([]check.Definition{{
				Name:          "persisted",
				Run:           tt.run,
				PersistentTTL: time.Hour,
				NewInfo: func() check.IpInfo {
					return &testInfo{}
				},
			}}, RunnerOptions{CacheDir: cacheDir})

			_, _ = runner.Run(net.ParseIP("1.1.1.1"))

			entries, err := os.ReadDir(cacheDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Fatalf("cache contains %d entries, want 0", len(entries))
			}
		})
	}
}

func TestRunnerSkipsPersistentCacheWhenTTLIsZero(t *testing.T) {
	var runs int32
	cacheDir := t.TempDir()
	ipaddr := net.ParseIP("1.1.1.1")
	definition := check.Definition{
		Name: "ttl-zero",
		Run: func(net.IP) (check.Check, error) {
			atomic.AddInt32(&runs, 1)
			return check.Check{
				Description: "wrong",
				Type:        check.Info,
				IpAddrInfo:  &testInfo{Value: "live"},
			}, nil
		},
		NewInfo: func() check.IpInfo {
			return &testInfo{}
		},
	}

	runner1 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{CacheDir: cacheDir})
	if _, errors := runner1.Run(ipaddr); len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}

	runner2 := NewRunnerWithOptions([]check.Definition{definition}, RunnerOptions{CacheDir: cacheDir})
	if _, errors := runner2.Run(ipaddr); len(errors) != 0 {
		t.Fatalf("got %d errors, want 0", len(errors))
	}

	if got := atomic.LoadInt32(&runs); got != 2 {
		t.Fatalf("check ran %d times, want 2", got)
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("cache contains %d entries, want 0", len(entries))
	}
}

type testInfo struct {
	Value string `json:"value"`
}

func (t testInfo) Summary() string {
	return t.Value
}

func (t testInfo) Json() ([]byte, error) {
	return json.Marshal(t)
}

func testDefinition(name string, ttl time.Duration, runs *int32) check.Definition {
	return check.Definition{
		Name:          name,
		PersistentTTL: ttl,
		NewInfo: func() check.IpInfo {
			return &testInfo{}
		},
		Run: func(net.IP) (check.Check, error) {
			run := atomic.AddInt32(runs, 1)
			return check.Check{
				Description:       "wrong",
				Type:              check.InfoAndIsMalicious,
				IpAddrIsMalicious: true,
				IpAddrInfo:        &testInfo{Value: "value-" + strconv.Itoa(int(run))},
			}, nil
		},
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w
	defer func() {
		os.Stdout = origStdout
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	return string(out)
}
