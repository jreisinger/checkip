package checkip

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

// Highlight makes string more visible.
func Highlight(s string) string {
	return fmt.Sprintf("%s", aurora.Magenta(s))
}

// Lowlight makes string less visible.
func Lowlight(s string) string {
	return fmt.Sprintf("%s", aurora.Gray(11, s))
}
