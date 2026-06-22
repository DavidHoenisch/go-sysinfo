//go:build integration

package omarchy

import "testing"

func TestIntegrationAvailable(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("omarchy not available on this system")
	}
}

func TestIntegrationGetInfo(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("omarchy not available on this system")
	}
	info := GetInfo(r)
	if info.Version == "" {
		t.Fatal("expected non-empty version")
	}
}

func TestIntegrationGetTheme(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("omarchy not available on this system")
	}
	if theme := GetTheme(r); theme == "" {
		t.Fatal("expected non-empty theme")
	}
}
