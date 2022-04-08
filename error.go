package checkip

import (
	"regexp"
	"runtime"
	"strings"
)

// Error is an error returned by a check.
type Error struct {
	err       error  // might contain secrets, like API keys
	ErrString string `json:"error"` // secrets redacted
}

func NewError(err error) *Error {
	callerName := "unknownCaller"
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		name := details.Name()
		parts := strings.Split(name, ".")
		callerName = parts[len(parts)-1]
	}
	return &Error{err: err, ErrString: callerName + ": " + redactSecrets(err.Error())}
}

func (e *Error) Error() string {
	return e.ErrString
}

func redactSecrets(s string) string {
	key := regexp.MustCompile(`(key|pass|password)=\w+`)
	return key.ReplaceAllString(s, "${1}=REDACTED")
}
