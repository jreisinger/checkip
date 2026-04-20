package main

import "testing"

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
