//go:build integration

package gnome

import (
	"testing"

	gosysinfo "github.com/DavidHoenisch/go-sysinfo"
)

func TestIntegrationAvailable(t *testing.T) {
	r := Reader{FS: gosysinfo.Reader{}}
	if !Available(r) {
		t.Skip("GNOME not available on this system")
	}
}

func TestIntegrationGetSessionInfo(t *testing.T) {
	r := Reader{FS: gosysinfo.Reader{}}
	if !Available(r) {
		t.Skip("GNOME not available on this system")
	}
	info := GetSessionInfo(r)
	if info.Desktop == "" {
		t.Fatal("expected non-empty desktop")
	}
}

func TestIntegrationGetShellVersion(t *testing.T) {
	r := Reader{FS: gosysinfo.Reader{}}
	if !Available(r) {
		t.Skip("GNOME not available on this system")
	}
	if version := GetShellVersion(r); version == "" {
		t.Fatal("expected non-empty shell version")
	}
}
