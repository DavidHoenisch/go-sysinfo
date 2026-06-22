package hyprland

import (
	"testing"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

const versionJSON = `{
    "branch": "v0.55.2",
    "commit": "39d7e209c79d451efab1b21151d5938289da838d",
    "version": "0.55.2",
    "dirty": false,
    "tag": "v0.55.2"
}`

const monitorsJSON = `[{
    "id": 0,
    "name": "eDP-1",
    "width": 2880,
    "height": 1800,
    "focused": true,
    "activeWorkspace": {"id": 1, "name": "1"}
}]`

const workspacesJSON = `[{
    "id": 1,
    "name": "1",
    "monitor": "eDP-1",
    "monitorID": 0,
    "windows": 1
}]`

const activeWorkspaceJSON = `{
    "id": 1,
    "name": "1",
    "monitor": "eDP-1",
    "monitorID": 0,
    "windows": 1
}`

const clientsJSON = `[{
    "address": "0x55a9685d2f10",
    "class": "Alacritty",
    "title": "Sitemap Builder",
    "workspace": {"id": 1, "name": "1"}
}]`

const activeWindowJSON = `{
    "address": "0x55a9685d2f10",
    "class": "Alacritty",
    "title": "Sitemap Builder",
    "workspace": {"id": 1, "name": "1"}
}`

const bindsJSON = `[{
    "key": "XF86AudioRaiseVolume",
    "description": "Volume up",
    "dispatcher": "exec",
    "arg": "swayosd-client --output-volume raise"
}]`

const configErrorsJSON = `["", "invalid syntax at line 5"]`

const optionJSON = `{"option": "decoration:rounding", "int": 0, "set": true}`

func mockHyprland() Reader {
	return Reader{
		Env: probe.MockEnvReader{
			"HYPRLAND_INSTANCE_SIGNATURE": "test-instance",
			"XDG_CURRENT_DESKTOP":         "Hyprland",
		},
		Cmd: probe.MockCommandRunner{
			"hyprctl\x00-j\x00version":           []byte(versionJSON),
			"hyprctl\x00-j\x00monitors":          []byte(monitorsJSON),
			"hyprctl\x00-j\x00workspaces":        []byte(workspacesJSON),
			"hyprctl\x00-j\x00activeworkspace":   []byte(activeWorkspaceJSON),
			"hyprctl\x00-j\x00clients":           []byte(clientsJSON),
			"hyprctl\x00-j\x00activewindow":      []byte(activeWindowJSON),
			"hyprctl\x00-j\x00binds":             []byte(bindsJSON),
			"hyprctl\x00-j\x00configerrors":      []byte(configErrorsJSON),
			"hyprctl\x00-j\x00getoption\x00decoration:rounding": []byte(optionJSON),
		},
	}
}

func TestAvailable(t *testing.T) {
	r := mockHyprland()
	if !Available(r) {
		t.Fatal("expected hyprland available")
	}
}

func TestAvailableWithoutSession(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{},
		Cmd: probe.MockCommandRunner{
			"hyprctl\x00-j\x00version": []byte(versionJSON),
		},
	}
	if Available(r) {
		t.Fatal("expected hyprland unavailable without session env")
	}
}

func TestAvailableWithoutHyprctl(t *testing.T) {
	r := Reader{
		Env: probe.MockEnvReader{
			"HYPRLAND_INSTANCE_SIGNATURE": "test-instance",
		},
		Cmd: probe.MockCommandRunner{},
	}
	if Available(r) {
		t.Fatal("expected hyprland unavailable without hyprctl")
	}
}

func TestGetVersion(t *testing.T) {
	r := mockHyprland()
	version := GetVersion(r)
	if version == nil {
		t.Fatal("expected version")
	}
	if version.Version != "0.55.2" {
		t.Fatalf("Version = %q", version.Version)
	}
}

func TestGetMonitors(t *testing.T) {
	r := mockHyprland()
	monitors := GetMonitors(r)
	if len(monitors) != 1 {
		t.Fatalf("len(monitors) = %d", len(monitors))
	}
	if monitors[0].Name != "eDP-1" {
		t.Fatalf("monitor name = %q", monitors[0].Name)
	}
}

func TestGetWorkspaces(t *testing.T) {
	r := mockHyprland()
	workspaces := GetWorkspaces(r)
	if len(workspaces) != 1 || workspaces[0].Name != "1" {
		t.Fatalf("workspaces = %+v", workspaces)
	}
}

func TestGetActiveWorkspace(t *testing.T) {
	r := mockHyprland()
	ws := GetActiveWorkspace(r)
	if ws == nil || ws.ID != 1 {
		t.Fatalf("active workspace = %+v", ws)
	}
}

func TestGetClients(t *testing.T) {
	r := mockHyprland()
	clients := GetClients(r)
	if len(clients) != 1 || clients[0].Class != "Alacritty" {
		t.Fatalf("clients = %+v", clients)
	}
}

func TestGetActiveWindow(t *testing.T) {
	r := mockHyprland()
	client := GetActiveWindow(r)
	if client == nil || client.Title != "Sitemap Builder" {
		t.Fatalf("active window = %+v", client)
	}
}

func TestGetBinds(t *testing.T) {
	r := mockHyprland()
	binds := GetBinds(r)
	if len(binds) != 1 || binds[0].Key != "XF86AudioRaiseVolume" {
		t.Fatalf("binds = %+v", binds)
	}
}

func TestGetConfigErrors(t *testing.T) {
	r := mockHyprland()
	errors := GetConfigErrors(r)
	if len(errors) != 1 || errors[0] != "invalid syntax at line 5" {
		t.Fatalf("config errors = %+v", errors)
	}
}

func TestGetOption(t *testing.T) {
	r := mockHyprland()
	option := GetOption(r, "decoration:rounding")
	if option == nil || option.Int != 0 || !option.Set {
		t.Fatalf("option = %+v", option)
	}
}

func TestOptionString(t *testing.T) {
	tests := []struct {
		name string
		opt  *Option
		want string
	}{
		{
			name: "custom",
			opt:  &Option{Custom: "5 5 5 5"},
			want: "5 5 5 5",
		},
		{
			name: "string",
			opt:  &Option{Str: "master"},
			want: "master",
		},
		{
			name: "float",
			opt:  &Option{Float: 1.5},
			want: "1.5",
		},
		{
			name: "int",
			opt:  &Option{Int: 10},
			want: "10",
		},
		{
			name: "nil",
			opt:  nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OptionString(tt.opt); got != tt.want {
				t.Fatalf("OptionString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetSessionInfo(t *testing.T) {
	r := mockHyprland()
	info := GetSessionInfo(r)
	if info.Version != "0.55.2" {
		t.Fatalf("Version = %q", info.Version)
	}
	if info.ActiveWorkspaceName != "1" {
		t.Fatalf("ActiveWorkspaceName = %q", info.ActiveWorkspaceName)
	}
	if info.FocusedMonitor != "eDP-1" {
		t.Fatalf("FocusedMonitor = %q", info.FocusedMonitor)
	}
	if info.ConfigErrorCount != 1 {
		t.Fatalf("ConfigErrorCount = %d", info.ConfigErrorCount)
	}
}

func TestUnavailableReturnsEmpty(t *testing.T) {
	r := Reader{Env: probe.MockEnvReader{}, Cmd: probe.MockCommandRunner{}}
	if Available(r) {
		t.Fatal("expected unavailable")
	}
	if GetVersion(r) != nil {
		t.Fatal("expected nil version")
	}
	if GetMonitors(r) != nil {
		t.Fatal("expected nil monitors")
	}
	if GetWorkspaces(r) != nil {
		t.Fatal("expected nil workspaces")
	}
	if GetActiveWorkspace(r) != nil {
		t.Fatal("expected nil active workspace")
	}
	if GetClients(r) != nil {
		t.Fatal("expected nil clients")
	}
	if GetActiveWindow(r) != nil {
		t.Fatal("expected nil active window")
	}
	if GetBinds(r) != nil {
		t.Fatal("expected nil binds")
	}
	if GetConfigErrors(r) != nil {
		t.Fatal("expected nil config errors")
	}
	if GetOption(r, "decoration:rounding") != nil {
		t.Fatal("expected nil option")
	}
	if info := GetSessionInfo(r); info.Version != "" || info.FocusedMonitor != "" {
		t.Fatal("expected empty session info")
	}
}
