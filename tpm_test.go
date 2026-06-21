package gosysinfo

import "testing"

func TestGetTpmVersion(t *testing.T) {
	oldExists := tpmBaseExists
	defer func() { tpmBaseExists = oldExists }()
	tpmBaseExists = func() bool { return true }

	reader := mockReader{
		tpmVersionMajor: "2",
	}
	if got := GetTpmVersion(reader); got != "2" {
		t.Errorf("GetTpmVersion() = %q, want %q", got, "2")
	}
}

func TestGetTpmVersionNoTPM(t *testing.T) {
	oldExists := tpmBaseExists
	defer func() { tpmBaseExists = oldExists }()
	tpmBaseExists = func() bool { return false }

	reader := mockReader{
		tpmVersionMajor: "2",
	}
	if got := GetTpmVersion(reader); got != "" {
		t.Errorf("GetTpmVersion() = %q, want empty", got)
	}
}

func TestGetTpmDescription(t *testing.T) {
	oldExists := tpmBaseExists
	defer func() { tpmBaseExists = oldExists }()
	tpmBaseExists = func() bool { return true }

	reader := mockReader{
		tpmDeviceDescription: "TPM 2.0 Device",
	}
	if got := GetTpmDescription(reader); got != "TPM 2.0 Device" {
		t.Errorf("GetTpmDescription() = %q, want %q", got, "TPM 2.0 Device")
	}
}

func TestGetTpmDescriptionNoTPM(t *testing.T) {
	oldExists := tpmBaseExists
	defer func() { tpmBaseExists = oldExists }()
	tpmBaseExists = func() bool { return false }

	reader := mockReader{
		tpmDeviceDescription: "TPM 2.0 Device",
	}
	if got := GetTpmDescription(reader); got != "" {
		t.Errorf("GetTpmDescription() = %q, want empty", got)
	}
}
