package util

import (
	"testing"
	"time"
)

func TestIsOlderThanWeek(t *testing.T) {
	y2k := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if !IsOlderThanOneWeek(y2k) {
		t.Errorf("%v is not older than a week", y2k)
	}
}
