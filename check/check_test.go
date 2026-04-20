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

func TestNonEmpty(t *testing.T) {
	tests := []struct {
		strings []string
		want    []string
	}{
		{[]string{""}, nil},
		{[]string{"", ""}, nil},
		{[]string{"a", ""}, []string{"a"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"", "b"}, []string{"b"}},
	}
	for _, test := range tests {
		if got := nonEmpty(test.strings...); !equal(got, test.want) {
			t.Errorf("nonEmpty(%v) = %v, want %v", test.strings, got, test.want)
		}
	}
}

func TestDefinitionsHaveUniqueNames(t *testing.T) {
	seen := make(map[string]struct{}, len(Definitions))

	for _, definition := range Definitions {
		if definition.Name == "" {
			t.Fatal("definition name must not be empty")
		}
		if _, ok := seen[definition.Name]; ok {
			t.Fatalf("duplicate definition name %q", definition.Name)
		}
		seen[definition.Name] = struct{}{}
	}
}

func TestWithoutActiveSkipsActiveDefinitions(t *testing.T) {
	definitions := []Definition{
		{Name: "passive-1"},
		{Name: "active", Active: true},
		{Name: "passive-2"},
	}

	filtered := WithoutActive(definitions)

	if len(filtered) != 2 {
		t.Fatalf("len(filtered) = %d, want 2", len(filtered))
	}
	if filtered[0].Name != "passive-1" {
		t.Fatalf("filtered[0].Name = %q, want %q", filtered[0].Name, "passive-1")
	}
	if filtered[1].Name != "passive-2" {
		t.Fatalf("filtered[1].Name = %q, want %q", filtered[1].Name, "passive-2")
	}
}

func TestDefinitionsMarkActiveChecks(t *testing.T) {
	tests := map[string]bool{
		"ping": true,
		"tls":  true,
	}

	for _, definition := range Definitions {
		wantActive, ok := tests[definition.Name]
		if !ok {
			continue
		}
		if definition.Active != wantActive {
			t.Fatalf("definition %q active = %t, want %t", definition.Name, definition.Active, wantActive)
		}
		delete(tests, definition.Name)
	}

	if len(tests) != 0 {
		t.Fatalf("missing definitions in test: %v", tests)
	}
}

// equal tells whether a and b contain the same elements. A nil argument is
// equivalent to an empty slice.
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
