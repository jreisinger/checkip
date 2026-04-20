package cli

import (
	"net"
	"sync"

	"github.com/jreisinger/checkip/check"
)

// Runner executes checks and memoizes cacheable results by IP address. It is
// safe for concurrent use.
type Runner struct {
	cachedDefinitions []check.Definition
	liveDefinitions   []check.Definition

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
	runner := &Runner{
		memo: make(map[string]*memoEntry),
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
	checks, errors := r.cachedRun(ipaddr)

	liveChecks, liveErrors := runDefinitions(r.liveDefinitions, ipaddr)
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

	checks, errors := runDefinitions(r.cachedDefinitions, ipaddr)
	entry.checks = cloneChecks(checks)
	entry.errors = cloneErrors(errors)

	return cloneChecks(checks), cloneErrors(errors)
}

func runDefinitions(definitions []check.Definition, ipaddr net.IP) (Checks, []error) {
	checks := make(Checks, 0, len(definitions))
	var errors []error

	for _, definition := range definitions {
		result, err := definition.Run(ipaddr)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		if result.Description == "" {
			result.Description = definition.Name
		}
		checks = append(checks, result)
	}

	return checks, errors
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
