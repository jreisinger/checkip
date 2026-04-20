package cli

import (
	"net"
	"sync"
	"time"

	"github.com/jreisinger/checkip/check"
)

// RunnerOptions configures runner caching behavior.
type RunnerOptions struct {
	DisableCache bool
	CacheDir     string
	Now          func() time.Time
}

// Runner executes checks and memoizes cacheable results by IP address. It is
// safe for concurrent use.
type Runner struct {
	cachedDefinitions []check.Definition
	liveDefinitions   []check.Definition
	disableCache      bool
	resultCache       *resultCache

	mu   sync.Mutex
	memo map[string]*memoEntry
}

type memoEntry struct {
	ready  chan struct{}
	checks Checks
	errors []error
}

// NewRunner creates a runner from check definitions.
func NewRunner(definitions []check.Definition) *Runner {
	return NewRunnerWithOptions(definitions, RunnerOptions{})
}

// NewRunnerWithOptions creates a runner from check definitions and options.
func NewRunnerWithOptions(definitions []check.Definition, options RunnerOptions) *Runner {
	runner := &Runner{
		disableCache: options.DisableCache,
	}
	if !options.DisableCache {
		runner.memo = make(map[string]*memoEntry)
		cache, err := newResultCache(options.CacheDir, options.Now)
		if err == nil {
			runner.resultCache = cache
		}
	}

	for _, definition := range definitions {
		if definition.Cache == check.CacheNone {
			runner.liveDefinitions = append(runner.liveDefinitions, definition)
			continue
		}
		runner.cachedDefinitions = append(runner.cachedDefinitions, definition)
	}

	return runner
}

// Run executes the configured checks for ipaddr.
func (r *Runner) Run(ipaddr net.IP) (Checks, []error) {
	if r.disableCache {
		return r.runAll(ipaddr)
	}

	checks, errors := r.cachedRun(ipaddr)

	liveChecks, liveErrors := r.runDefinitions(r.liveDefinitions, ipaddr)
	checks = append(checks, liveChecks...)
	errors = append(errors, liveErrors...)

	return checks, errors
}

func (r *Runner) runAll(ipaddr net.IP) (Checks, []error) {
	checks, errors := r.runDefinitions(r.cachedDefinitions, ipaddr)

	liveChecks, liveErrors := r.runDefinitions(r.liveDefinitions, ipaddr)
	checks = append(checks, liveChecks...)
	errors = append(errors, liveErrors...)

	return checks, errors
}

func (r *Runner) cachedRun(ipaddr net.IP) (Checks, []error) {
	if len(r.cachedDefinitions) == 0 {
		return nil, nil
	}

	key := ipaddr.String()

	r.mu.Lock()
	if entry, ok := r.memo[key]; ok {
		r.mu.Unlock()
		<-entry.ready
		return cloneChecks(entry.checks), cloneErrors(entry.errors)
	}

	entry := &memoEntry{ready: make(chan struct{})}
	r.memo[key] = entry
	r.mu.Unlock()
	defer close(entry.ready)

	checks, errors := r.runDefinitions(r.cachedDefinitions, ipaddr)
	entry.checks = cloneChecks(checks)
	entry.errors = cloneErrors(errors)

	return cloneChecks(checks), cloneErrors(errors)
}

func (r *Runner) runDefinitions(definitions []check.Definition, ipaddr net.IP) (Checks, []error) {
	checks := make(Checks, 0, len(definitions))
	var errors []error

	for _, definition := range definitions {
		result, err := r.runDefinition(definition, ipaddr)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		checks = append(checks, result)
	}

	return checks, errors
}

func (r *Runner) runDefinition(definition check.Definition, ipaddr net.IP) (check.Check, error) {
	ip := ipaddr.String()
	if r.resultCache != nil && definition.PersistentTTL > 0 {
		if result, ok := r.resultCache.load(definition, ip); ok {
			return result, nil
		}
	}

	result, err := definition.Run(ipaddr)
	if err != nil {
		return check.Check{}, err
	}
	result.Description = definition.Name

	if r.resultCache != nil && definition.PersistentTTL > 0 && result.MissingCredentials == "" {
		r.resultCache.store(definition, ip, result)
	}

	return result, nil
}

func cloneChecks(checks Checks) Checks {
	if len(checks) == 0 {
		return nil
	}

	cloned := make(Checks, len(checks))
	copy(cloned, checks)
	return cloned
}

func cloneErrors(errors []error) []error {
	if len(errors) == 0 {
		return nil
	}

	cloned := make([]error, len(errors))
	copy(cloned, errors)
	return cloned
}
