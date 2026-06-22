package hyprland

import (
	"encoding/json"
	"strings"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

func isHyprlandSession(env probe.EnvReader) bool {
	if env.Get("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return true
	}
	value := strings.ToUpper(env.Get("XDG_CURRENT_DESKTOP"))
	return strings.Contains(value, "HYPRLAND")
}

func Available(r Reader) bool {
	if !isHyprlandSession(r.env()) {
		return false
	}
	_, err := r.cmd().Run("hyprctl", "-j", "version")
	return err == nil
}

func (r Reader) Available() bool {
	return Available(r)
}

func hyprctlJSON(r Reader, args ...string) ([]byte, bool) {
	if !Available(r) {
		return nil, false
	}
	cmdArgs := append([]string{"-j"}, args...)
	out, err := r.cmd().Run("hyprctl", cmdArgs...)
	if err != nil {
		return nil, false
	}
	return out, true
}

func GetSessionInfo(r Reader) *SessionInfo {
	if !Available(r) {
		return &SessionInfo{}
	}

	info := &SessionInfo{}

	if version := GetVersion(r); version != nil {
		info.Version = version.Version
	}

	if ws := GetActiveWorkspace(r); ws != nil {
		info.ActiveWorkspaceID = ws.ID
		info.ActiveWorkspaceName = ws.Name
	}

	for _, monitor := range GetMonitors(r) {
		if monitor.Focused {
			info.FocusedMonitor = monitor.Name
			break
		}
	}

	errors := GetConfigErrors(r)
	info.ConfigErrorCount = len(errors)

	return info
}

func (r Reader) GetSessionInfo() *SessionInfo {
	return GetSessionInfo(r)
}

func decodeJSON[T any](data []byte, dest *T) bool {
	if len(data) == 0 {
		return false
	}
	return json.Unmarshal(data, dest) == nil
}
