package checkip

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

// highlight makes string more visible.
func highlight(s string) string {
	return fmt.Sprintf("%s", aurora.Magenta(s))
}

// lowlight makes string less visible.
func lowlight(s string) string {
	return fmt.Sprintf("%s", aurora.Gray(11, s))
}
