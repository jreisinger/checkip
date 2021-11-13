package check

import "regexp"

// Error is an error returned by a check.
type Error struct {
	err       error  // might contain secrets, like API keys
	ErrString string `json:"error"` // secrets redacted
}

func NewError(err error) *Error {
	return &Error{err: err, ErrString: redactSecrets(err.Error())}
}

func (e *Error) Error() string {
	return e.ErrString
}

func redactSecrets(s string) string {
	key := regexp.MustCompile(`(key|pass|password)=\w+`)
	return key.ReplaceAllString(s, "${1}=REDACTED")
}
