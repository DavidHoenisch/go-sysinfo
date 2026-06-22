package gosysinfo

import (
	"reflect"
	"testing"
)

func TestGetOSRelease(t *testing.T) {
	reader := mockReader{
		etcOSRelease: `NAME="Arch Linux"
PRETTY_NAME="Arch Linux"
ID=arch
BUILD_ID=rolling
ANSI_COLOR="38;2;23;147;209"
HOME_URL="https://archlinux.org/"
DOCUMENTATION_URL="https://wiki.archlinux.org/"
SUPPORT_URL="https://bbs.archlinux.org/"
BUG_REPORT_URL="https://gitlab.archlinux.org/groups/archlinux/-/issues"
PRIVACY_POLICY_URL="https://terms.archlinux.org/docs/privacy-policy/"
LOGO=archlinux-logo
`,
	}

	got := GetOSRelease(reader)
	want := &OSReleaseInfo{
		Name:             "Arch Linux",
		PrettyName:       "Arch Linux",
		ID:               "arch",
		BuildID:          "rolling",
		ANSIColor:        "38;2;23;147;209",
		HomeURL:          "https://archlinux.org/",
		DocumentationURL: "https://wiki.archlinux.org/",
		SupportURL:       "https://bbs.archlinux.org/",
		BugReportURL:     "https://gitlab.archlinux.org/groups/archlinux/-/issues",
		PrivacyPolicyURL: "https://terms.archlinux.org/docs/privacy-policy/",
		Logo:             "archlinux-logo",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetOSRelease() = %#v, want %#v", got, want)
	}
}

func TestGetOSReleaseEmpty(t *testing.T) {
	got := GetOSRelease(mockReader{})
	want := &OSReleaseInfo{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetOSRelease() = %#v, want %#v", got, want)
	}
}

func TestUnquoteOSReleaseValue(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`arch`, "arch"},
		{`"Arch Linux"`, "Arch Linux"},
		{`"Say \"hello\""`, `Say "hello"`},
		{`"line1\nline2"`, "line1\nline2"},
	}

	for _, tt := range tests {
		if got := unquoteOSReleaseValue(tt.in); got != tt.want {
			t.Errorf("unquoteOSReleaseValue(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
