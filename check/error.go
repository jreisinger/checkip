package check

import "regexp"

type ResultError struct {
	err error
}

func NewResultError(err error) *ResultError {
	return &ResultError{err: err}
}

func (e *ResultError) Error() string {
	return redactSecrets(e.err.Error())
}

func redactSecrets(s string) string {
	key := regexp.MustCompile(`(key|pass|password)=\w+`)
	return key.ReplaceAllString(s, "${1}=REDACTED")
}
