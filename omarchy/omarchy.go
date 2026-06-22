package omarchy

import "strings"

type Info struct {
	Version string
	Branch  string
	Channel string
	Theme   string
	Font    string
}

func Available(r Reader) bool {
	_, err := r.cmd().Run("omarchy", "version")
	return err == nil
}

func (r Reader) Available() bool {
	return Available(r)
}

func GetVersion(r Reader) string {
	if !Available(r) {
		return ""
	}
	out, err := r.cmd().Run("omarchy", "version")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetVersion() string {
	return GetVersion(r)
}

func GetBranch(r Reader) string {
	if !Available(r) {
		return ""
	}
	out, err := r.cmd().Run("omarchy", "version", "branch")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetBranch() string {
	return GetBranch(r)
}

func GetChannel(r Reader) string {
	if !Available(r) {
		return ""
	}
	out, err := r.cmd().Run("omarchy", "version", "channel")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetChannel() string {
	return GetChannel(r)
}

func GetTheme(r Reader) string {
	if !Available(r) {
		return ""
	}
	out, err := r.cmd().Run("omarchy", "theme", "current")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetTheme() string {
	return GetTheme(r)
}

func GetFont(r Reader) string {
	if !Available(r) {
		return ""
	}
	out, err := r.cmd().Run("omarchy", "font", "current")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (r Reader) GetFont() string {
	return GetFont(r)
}

func IsToggleEnabled(r Reader, name string) bool {
	if !Available(r) || name == "" {
		return false
	}
	_, err := r.cmd().Run("omarchy", "toggle", "enabled", name)
	return err == nil
}

func (r Reader) IsToggleEnabled(name string) bool {
	return IsToggleEnabled(r, name)
}

func GetInfo(r Reader) *Info {
	if !Available(r) {
		return &Info{}
	}
	return &Info{
		Version: GetVersion(r),
		Branch:  GetBranch(r),
		Channel: GetChannel(r),
		Theme:   GetTheme(r),
		Font:    GetFont(r),
	}
}

func (r Reader) GetInfo() *Info {
	return GetInfo(r)
}
