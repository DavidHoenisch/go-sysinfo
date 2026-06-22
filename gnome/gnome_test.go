package gnome

import (
	"testing"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

func TestAvailable(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{
			"XDG_CURRENT_DESKTOP": "GNOME",
		},
		Cmd: probe.MockCommandRunner{
			"gnome-shell\x00--version": []byte("GNOME Shell 46.0\n"),
		},
	}
	if !Available(r) {
		t.Fatal("expected gnome available")
	}
}

func TestAvailableWithoutGnomeEnv(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{
			"XDG_CURRENT_DESKTOP": "KDE",
		},
		Cmd: probe.MockCommandRunner{
			"gnome-shell\x00--version": []byte("GNOME Shell 46.0\n"),
		},
	}
	if Available(r) {
		t.Fatal("expected gnome unavailable outside GNOME session")
	}
}

func TestGetShellVersion(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{
			"XDG_CURRENT_DESKTOP": "GNOME",
		},
		Cmd: probe.MockCommandRunner{
			"gnome-shell\x00--version": []byte("GNOME Shell 46.0\n"),
		},
	}
	version := GetShellVersion(r)
	if version != "GNOME Shell 46.0" {
		t.Fatalf("GetShellVersion() = %q", version)
	}
}

func TestGetSetting(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{
			"XDG_CURRENT_DESKTOP": "GNOME",
			"DESKTOP_SESSION":     "gnome",
		},
		Cmd: probe.MockCommandRunner{
			"gsettings\x00get\x00org.gnome.desktop.interface\x00gtk-theme": []byte("'Adwaita'\n"),
			"gnome-shell\x00--version":                                     []byte("GNOME Shell 46.0\n"),
		},
	}
	setting := GetSetting(r, "org.gnome.desktop.interface", "gtk-theme")
	if setting != "'Adwaita'" {
		t.Fatalf("GetSetting() = %q", setting)
	}
}

func TestGetSessionInfo(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{
			"XDG_CURRENT_DESKTOP": "GNOME",
			"DESKTOP_SESSION":     "gnome",
		},
		Cmd: probe.MockCommandRunner{
			"gnome-shell\x00--version": []byte("GNOME Shell 46.0\n"),
		},
	}
	info := GetSessionInfo(r)
	if info.ShellVersion != "GNOME Shell 46.0" {
		t.Fatalf("ShellVersion = %q", info.ShellVersion)
	}
	if info.Desktop != "GNOME" {
		t.Fatalf("Desktop = %q", info.Desktop)
	}
	if info.Session != "gnome" {
		t.Fatalf("Session = %q", info.Session)
	}
}

func TestUnavailableReturnsEmpty(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{},
		Cmd: probe.MockCommandRunner{},
	}
	if Available(r) {
		t.Fatal("expected unavailable")
	}
	if GetShellVersion(r) != "" {
		t.Fatal("expected empty shell version")
	}
	if GetSetting(r, "org.gnome.desktop.interface", "gtk-theme") != "" {
		t.Fatal("expected empty setting")
	}
	if info := GetSessionInfo(r); info.ShellVersion != "" || info.Desktop != "" {
		t.Fatal("expected empty session info")
	}
}
