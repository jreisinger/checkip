package check

import "regexp"

type Error struct {
	err       error
	ErrString string `json:"error"`
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
