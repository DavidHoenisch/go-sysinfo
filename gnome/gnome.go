package gnome

import (
	"strings"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

type SessionInfo struct {
	ShellVersion string
	Desktop      string
	Session      string
}

func isGnomeDesktop(env probe.EnvReader) bool {
	for _, key := range []string{"XDG_CURRENT_DESKTOP", "DESKTOP_SESSION"} {
		value := strings.ToUpper(env.Get(key))
		if strings.Contains(value, "GNOME") {
			return true
		}
	}
	return false
}

func hasGnomeBinary(cmd probe.CommandRunner) bool {
	for _, name := range []string{"gnome-shell", "gsettings"} {
		if _, err := cmd.Run(name, "--version"); err == nil {
			return true
		}
	}
	return false
}

func Available(r Reader) bool {
	if !isGnomeDesktop(r.env()) {
		return false
	}
	return hasGnomeBinary(r.cmd())
}

func (r Reader) Available() bool {
	return Available(r)
}

func GetShellVersion(r Reader) string {
	if !Available(r) {
		return ""
	}
	out, err := r.cmd().Run("gnome-shell", "--version")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetShellVersion() string {
	return GetShellVersion(r)
}

func GetSetting(r Reader, schema, key string) string {
	if !Available(r) {
		return ""
	}
	if schema == "" || key == "" {
		return ""
	}
	out, err := r.cmd().Run("gsettings", "get", schema, key)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetSetting(schema, key string) string {
	return GetSetting(r, schema, key)
}

func GetSessionInfo(r Reader) *SessionInfo {
	if !Available(r) {
		return &SessionInfo{}
	}
	return &SessionInfo{
		ShellVersion: GetShellVersion(r),
		Desktop:      r.env().Get("XDG_CURRENT_DESKTOP"),
		Session:      r.env().Get("DESKTOP_SESSION"),
	}
}

func (r Reader) GetSessionInfo() *SessionInfo {
	return GetSessionInfo(r)
}
