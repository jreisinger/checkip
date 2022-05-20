package check

import "testing"

func TestNa(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"", "n/a"},
		{" ", "n/a"},
		{"\t\n", "n/a"},
		{"a", "a"},
		{"a\tb", "a\tb"},
	}
	for _, test := range tests {
		if got := na(test.s); got != test.want {
			t.Errorf("na(%s) = %s, want %s", test.s, got, test.want)
		}
	}
}
