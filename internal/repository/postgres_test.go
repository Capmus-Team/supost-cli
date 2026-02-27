package repository

import (
	"strings"
	"testing"
)

func TestClampRecentLimit(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{name: "negative", input: -1, want: maxRecentActivePosts},
		{name: "zero", input: 0, want: maxRecentActivePosts},
		{name: "within range", input: 12, want: 12},
		{name: "too high", input: 999, want: maxRecentActivePosts},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clampRecentLimit(tt.input)
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEnsurePoolerSafeConnectionString_URLStyle(t *testing.T) {
	input := "postgresql://u:p@localhost:5432/db?sslmode=require"
	out := ensurePoolerSafeConnectionString(input)

	if !strings.Contains(out, "default_query_exec_mode=simple_protocol") {
		t.Fatalf("missing default_query_exec_mode in %q", out)
	}
	if !strings.Contains(out, "statement_cache_capacity=0") {
		t.Fatalf("missing statement_cache_capacity in %q", out)
	}
	if !strings.Contains(out, "description_cache_capacity=0") {
		t.Fatalf("missing description_cache_capacity in %q", out)
	}
}

func TestEnsurePoolerSafeConnectionString_KeyValueStyle(t *testing.T) {
	input := "host=localhost port=5432 user=postgres password=secret dbname=postgres sslmode=disable"
	out := ensurePoolerSafeConnectionString(input)

	if !strings.Contains(out, "default_query_exec_mode=simple_protocol") {
		t.Fatalf("missing default_query_exec_mode in %q", out)
	}
	if !strings.Contains(out, "statement_cache_capacity=0") {
		t.Fatalf("missing statement_cache_capacity in %q", out)
	}
	if !strings.Contains(out, "description_cache_capacity=0") {
		t.Fatalf("missing description_cache_capacity in %q", out)
	}
}
