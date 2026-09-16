package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the grey struct methods
func TestGreyStruct(t *testing.T) {
	t.Run("Summary method formats correctly", func(t *testing.T) {
		g := grey{
			IP:      "1.2.3.4",
			Riot:    true,
			Message: "Test message",
		}
		summary := g.Summary()
		assert.Equal(t, "Test message", summary)
	})

	t.Run("Summary method handles empty message", func(t *testing.T) {
		g := grey{
			IP:   "1.2.3.4",
			Riot: false,
		}
		summary := g.Summary()
		assert.Equal(t, "n/a", summary)
	})
}
