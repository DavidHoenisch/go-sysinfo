//go:build integration

package clamav

import (
	"testing"

	gosysinfo "github.com/DavidHoenisch/go-sysinfo"
)

func TestIntegrationAvailable(t *testing.T) {
	r := Reader{FS: gosysinfo.Reader{}}
	if !Available(r) {
		t.Skip("clamd not available on this system")
	}
}

func TestIntegrationGetVersion(t *testing.T) {
	r := Reader{FS: gosysinfo.Reader{}}
	if !Available(r) {
		t.Skip("clamd not available on this system")
	}
	if version := GetVersion(r); version == "" {
		t.Fatal("expected non-empty version")
	}
}

func TestIntegrationGetDatabaseStats(t *testing.T) {
	r := Reader{FS: gosysinfo.Reader{}}
	if !Available(r) {
		t.Skip("clamd not available on this system")
	}
	stats := GetDatabaseStats(r)
	if stats.Raw == "" {
		t.Fatal("expected non-empty stats")
	}
}
