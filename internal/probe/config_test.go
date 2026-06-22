package probe

import "testing"

func TestParseKeyValueConfig(t *testing.T) {
	content := `
# comment
LocalSocket /run/clamav/clamd.ctl
DatabaseDirectory /var/lib/clamav
`
	values := ParseKeyValueConfig(content)
	if values["LocalSocket"] != "/run/clamav/clamd.ctl" {
		t.Fatalf("LocalSocket = %q", values["LocalSocket"])
	}
	if values["DatabaseDirectory"] != "/var/lib/clamav" {
		t.Fatalf("DatabaseDirectory = %q", values["DatabaseDirectory"])
	}
}
