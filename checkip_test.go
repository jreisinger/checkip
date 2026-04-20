package main

import (
	"testing"

	"github.com/jreisinger/checkip/check"
)

func TestValidateParallelism(t *testing.T) {
	tests := []struct {
		name        string
		parallelism int
		wantErr     bool
	}{
		{name: "negative", parallelism: -1, wantErr: true},
		{name: "zero", parallelism: 0, wantErr: true},
		{name: "positive", parallelism: 1, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateParallelism(tt.parallelism)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateParallelism(%d) error = %v, wantErr %t", tt.parallelism, err, tt.wantErr)
			}
		})
	}
}

func TestSelectedDefinitionsKeepsAllChecksByDefault(t *testing.T) {
	selected := selectedDefinitions(false)

	if len(selected) != len(check.Definitions) {
		t.Fatalf("len(selected) = %d, want %d", len(selected), len(check.Definitions))
	}
}

func TestSelectedDefinitionsCanDisableActiveChecks(t *testing.T) {
	selected := selectedDefinitions(true)

	if len(selected) >= len(check.Definitions) {
		t.Fatalf("len(selected) = %d, want less than %d", len(selected), len(check.Definitions))
	}

	for _, definition := range selected {
		if definition.Active {
			t.Fatalf("selected active definition %q", definition.Name)
		}
	}

	for _, active := range []string{"ping", "tls"} {
		if containsDefinition(selected, active) {
			t.Fatalf("selected definitions unexpectedly contain %q", active)
		}
	}
}

func containsDefinition(definitions []check.Definition, name string) bool {
	for _, definition := range definitions {
		if definition.Name == name {
			return true
		}
	}
	return false
}
