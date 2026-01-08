package uuidx

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUUID(t *testing.T) {
	u := NewUUID()
	if len(u) == 0 {
		t.Error("NewUUID returned empty string")
	}
	if !IsValidUUID(u) {
		t.Errorf("NewUUID returned invalid UUID: %s", u)
	}
}

func TestNewUUIDv4(t *testing.T) {
	u := NewUUIDv4()
	if len(u) == 0 {
		t.Error("NewUUIDv4 returned empty string")
	}
	if !IsValidUUID(u) {
		t.Errorf("NewUUIDv4 returned invalid UUID: %s", u)
	}
}

func TestNewUUIDv5(t *testing.T) {
	ns := uuid.NameSpaceDNS
	name := "example.com"
	u1 := NewUUIDv5(ns, name)
	u2 := NewUUIDv5(ns, name)

	if len(u1) == 0 {
		t.Error("NewUUIDv5 returned empty string")
	}
	if !IsValidUUID(u1) {
		t.Errorf("NewUUIDv5 returned invalid UUID: %s", u1)
	}
	if u1 != u2 {
		t.Error("NewUUIDv5 should be deterministic")
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"invalid", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsValidUUID(tt.in); got != tt.want {
			t.Errorf("IsValidUUID(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
