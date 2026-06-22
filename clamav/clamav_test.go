package clamav

import (
	"strings"
	"testing"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

type mockFS map[string]string

func (m mockFS) Read(path string) string {
	return m[path]
}

func mockClamdSocket(responses map[string]string) probe.MockSocketDialer {
	return probe.MockSocketDialer{
		DialFn: func(network, address string) (probe.SocketConn, error) {
			var pending string
			return &probe.MockSocketConn{
				WriteFn: func(b []byte) (int, error) {
					cmd := strings.TrimRight(string(b), "\x00")
					cmd = strings.TrimPrefix(cmd, "z")
					pending = cmd
					return len(b), nil
				},
				ReadFn: func(b []byte) (int, error) {
					response, ok := responses[pending]
					if !ok {
						return 0, nil
					}
					copy(b, response)
					return len(response), nil
				},
			}, nil
		},
	}
}

func TestAvailable(t *testing.T) {
	r := Reader{
		FS: mockFS{
			defaultClamdConf: "LocalSocket /tmp/clamd.sock\n",
		},
		Socket: mockClamdSocket(map[string]string{
			"PING": "PONG\x00",
		}),
	}
	if !Available(r) {
		t.Fatal("expected clamav available via socket ping")
	}
}

func TestAvailableViaBinary(t *testing.T) {
	r := Reader{
		FS:     mockFS{},
		Socket: probe.MockSocketDialer{},
		Cmd: probe.MockCommandRunner{
			"clamd\x00--version": []byte("ClamAV 1.4.0\n"),
		},
	}
	if !Available(r) {
		t.Fatal("expected available via clamd --version fallback")
	}
}

func TestGetVersionFromSocket(t *testing.T) {
	r := Reader{
		FS: mockFS{
			defaultClamdConf: "LocalSocket /tmp/clamd.sock\n",
		},
		Socket: mockClamdSocket(map[string]string{
			"VERSION": "ClamAV 1.4.0/26862/Mon Jan  1 00:00:00 2024\x00",
		}),
	}
	version := GetVersion(r)
	if version == "" || !strings.Contains(version, "ClamAV 1.4.0") {
		t.Fatalf("GetVersion() = %q", version)
	}
}

func TestGetVersionFallbackCLI(t *testing.T) {
	r := Reader{
		FS:     mockFS{},
		Socket: probe.MockSocketDialer{},
		Cmd: probe.MockCommandRunner{
			"clamd\x00--version": []byte("ClamAV 1.4.0/26862\n"),
		},
	}
	version := GetVersion(r)
	if version != "ClamAV 1.4.0/26862" {
		t.Fatalf("GetVersion() = %q", version)
	}
}

func TestGetDatabaseStats(t *testing.T) {
	r := Reader{
		FS: mockFS{
			defaultClamdConf: "LocalSocket /tmp/clamd.sock\n",
		},
		Socket: mockClamdSocket(map[string]string{
			"STATS": "POOLS: 1\nsignatures: 9876543\n",
		}),
	}
	stats := GetDatabaseStats(r)
	if stats.SignatureCount != "9876543" {
		t.Fatalf("SignatureCount = %q", stats.SignatureCount)
	}
}

func TestGetConfig(t *testing.T) {
	r := Reader{
		FS: mockFS{
			defaultClamdConf:     "LocalSocket /custom/clamd.sock\nDatabaseDirectory /var/lib/clamav\n",
			defaultFreshclamConf: "DatabaseDirectory /var/lib/clamav-fresh\n",
		},
	}
	cfg := GetConfig(r)
	if cfg.LocalSocket != "/custom/clamd.sock" {
		t.Fatalf("LocalSocket = %q", cfg.LocalSocket)
	}
	if cfg.DatabaseDirectory != "/var/lib/clamav" {
		t.Fatalf("DatabaseDirectory = %q", cfg.DatabaseDirectory)
	}
	if cfg.FreshClamDatabaseDir != "/var/lib/clamav-fresh" {
		t.Fatalf("FreshClamDatabaseDir = %q", cfg.FreshClamDatabaseDir)
	}
}

func TestGetInfo(t *testing.T) {
	r := Reader{
		FS: mockFS{
			defaultClamdConf: "LocalSocket /tmp/clamd.sock\nDatabaseDirectory /var/lib/clamav\n",
		},
		Socket: mockClamdSocket(map[string]string{
			"PING":    "PONG\x00",
			"VERSION": "ClamAV 1.4.0\x00",
			"STATS":   "signatures: 42\n",
		}),
	}
	info := GetInfo(r)
	if !info.Running {
		t.Fatal("expected Running true")
	}
	if info.Version != "ClamAV 1.4.0" {
		t.Fatalf("Version = %q", info.Version)
	}
	if info.SignatureCount != "42" {
		t.Fatalf("SignatureCount = %q", info.SignatureCount)
	}
	if info.DatabasePath != "/var/lib/clamav" {
		t.Fatalf("DatabasePath = %q", info.DatabasePath)
	}
}

func TestUnavailableReturnsEmpty(t *testing.T) {
	r := Reader{
		FS:     mockFS{},
		Socket: probe.MockSocketDialer{},
		Cmd:    probe.MockCommandRunner{},
	}
	if Available(r) {
		t.Fatal("expected unavailable")
	}
	if GetVersion(r) != "" {
		t.Fatal("expected empty version")
	}
	if stats := GetDatabaseStats(r); stats.SignatureCount != "" {
		t.Fatal("expected empty database stats")
	}
}
