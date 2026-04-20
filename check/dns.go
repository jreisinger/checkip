package check

import "strings"

func trimTrailingDot(name string) string {
	return strings.TrimSuffix(name, ".")
}
