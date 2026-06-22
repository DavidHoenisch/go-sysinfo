package omarchy

import (
	"testing"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

func mockAvailable() Reader {
	return Reader{
		Cmd: probe.MockCommandRunner{
			"omarchy\x00version": []byte("3.8.2\n"),
		},
	}
}

func TestAvailable(t *testing.T) {
	r := mockAvailable()
	if !Available(r) {
		t.Fatal("expected omarchy available")
	}
}

func TestAvailableWithoutBinary(t *testing.T) {
	r := Reader{Cmd: probe.MockCommandRunner{}}
	if Available(r) {
		t.Fatal("expected omarchy unavailable")
	}
}

func TestGetVersion(t *testing.T) {
	r := mockAvailable()
	version := GetVersion(r)
	if version != "3.8.2" {
		t.Fatalf("GetVersion() = %q", version)
	}
}

func TestGetBranch(t *testing.T) {
	r := Reader{
		Cmd: probe.MockCommandRunner{
			"omarchy\x00version":        []byte("3.8.2\n"),
			"omarchy\x00version\x00branch": []byte("master\n"),
		},
	}
	if branch := GetBranch(r); branch != "master" {
		t.Fatalf("GetBranch() = %q", branch)
	}
}

func TestGetChannel(t *testing.T) {
	r := Reader{
		Cmd: probe.MockCommandRunner{
			"omarchy\x00version":         []byte("3.8.2\n"),
			"omarchy\x00version\x00channel": []byte("stable\n"),
		},
	}
	if channel := GetChannel(r); channel != "stable" {
		t.Fatalf("GetChannel() = %q", channel)
	}
}

func TestGetTheme(t *testing.T) {
	r := Reader{
		Cmd: probe.MockCommandRunner{
			"omarchy\x00version":              []byte("3.8.2\n"),
			"omarchy\x00theme\x00current": []byte("Catppuccin\n"),
		},
	}
	if theme := GetTheme(r); theme != "Catppuccin" {
		t.Fatalf("GetTheme() = %q", theme)
	}
}

func TestGetFont(t *testing.T) {
	r := Reader{
		Cmd: probe.MockCommandRunner{
			"omarchy\x00version":            []byte("3.8.2\n"),
			"omarchy\x00font\x00current": []byte("JetBrainsMono Nerd Font\n"),
		},
	}
	if font := GetFont(r); font != "JetBrainsMono Nerd Font" {
		t.Fatalf("GetFont() = %q", font)
	}
}

func TestIsToggleEnabled(t *testing.T) {
	r := Reader{
		Cmd: probe.MockCommandRunner{
			"omarchy\x00version":                      []byte("3.8.2\n"),
			"omarchy\x00toggle\x00enabled\x00nightlight": []byte(""),
		},
	}
	if !IsToggleEnabled(r, "nightlight") {
		t.Fatal("expected nightlight toggle enabled")
	}
	if IsToggleEnabled(r, "idle") {
		t.Fatal("expected idle toggle disabled")
	}
}

func TestGetInfo(t *testing.T) {
	r := Reader{
		Cmd: probe.MockCommandRunner{
			"omarchy\x00version":              []byte("3.8.2\n"),
			"omarchy\x00version\x00branch":    []byte("master\n"),
			"omarchy\x00version\x00channel":   []byte("stable\n"),
			"omarchy\x00theme\x00current":     []byte("Catppuccin\n"),
			"omarchy\x00font\x00current":      []byte("JetBrainsMono Nerd Font\n"),
		},
	}
	info := GetInfo(r)
	if info.Version != "3.8.2" {
		t.Fatalf("Version = %q", info.Version)
	}
	if info.Branch != "master" {
		t.Fatalf("Branch = %q", info.Branch)
	}
	if info.Channel != "stable" {
		t.Fatalf("Channel = %q", info.Channel)
	}
	if info.Theme != "Catppuccin" {
		t.Fatalf("Theme = %q", info.Theme)
	}
	if info.Font != "JetBrainsMono Nerd Font" {
		t.Fatalf("Font = %q", info.Font)
	}
}

func TestUnavailableReturnsEmpty(t *testing.T) {
	r := Reader{Cmd: probe.MockCommandRunner{}}
	if Available(r) {
		t.Fatal("expected unavailable")
	}
	if GetVersion(r) != "" {
		t.Fatal("expected empty version")
	}
	if GetBranch(r) != "" {
		t.Fatal("expected empty branch")
	}
	if GetChannel(r) != "" {
		t.Fatal("expected empty channel")
	}
	if GetTheme(r) != "" {
		t.Fatal("expected empty theme")
	}
	if GetFont(r) != "" {
		t.Fatal("expected empty font")
	}
	if IsToggleEnabled(r, "nightlight") {
		t.Fatal("expected toggle disabled when unavailable")
	}
	if info := GetInfo(r); info.Version != "" || info.Theme != "" {
		t.Fatal("expected empty info")
	}
}
