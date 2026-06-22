//go:build integration

package hyprland

import "testing"

func TestIntegrationAvailable(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("hyprland not available on this system")
	}
}

func TestIntegrationGetVersion(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("hyprland not available on this system")
	}
	version := GetVersion(r)
	if version == nil || version.Version == "" {
		t.Fatal("expected non-empty version")
	}
}

func TestIntegrationGetMonitors(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("hyprland not available on this system")
	}
	monitors := GetMonitors(r)
	if len(monitors) == 0 {
		t.Fatal("expected at least one monitor")
	}
}

func TestIntegrationGetConfigErrors(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("hyprland not available on this system")
	}
	_ = GetConfigErrors(r)
}

func TestIntegrationGetOption(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("hyprland not available on this system")
	}
	option := GetOption(r, "decoration:rounding")
	if option == nil {
		t.Fatal("expected option for decoration:rounding")
	}
}

func TestIntegrationGetSessionInfo(t *testing.T) {
	r := Reader{}
	if !Available(r) {
		t.Skip("hyprland not available on this system")
	}
	info := GetSessionInfo(r)
	if info.Version == "" {
		t.Fatal("expected non-empty session version")
	}
	if info.FocusedMonitor == "" {
		t.Fatal("expected focused monitor")
	}
}
