package check

import (
	"regexp"
	"runtime"
	"strings"
)

// checkError is an error that should be returned by a check.
type checkError struct {
	err       error  // might contain secrets, like API keys
	ErrString string `json:"error"` // secrets redacted, ok to print
}

// newCheckError returns an error that contains caller name and has potential
// secrets redacted.
func newCheckError(err error) *checkError {
	callerName := "unknownCaller"
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		name := details.Name()
		parts := strings.Split(name, ".")
		callerName = parts[len(parts)-1]
	}
	return &checkError{err: err, ErrString: callerName + ": " + redactSecrets(err.Error())}
}

func (e *checkError) Error() string {
	return e.ErrString
}

func redactSecrets(s string) string {
	key := regexp.MustCompile(`(key|pass|password)=\w+`)
	return key.ReplaceAllString(s, "${1}=REDACTED")
}
