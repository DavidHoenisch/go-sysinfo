package gosysinfo

import (
	"errors"
	"os"
)

const (
	tpm_base               string = "/sys/class/tpm/tpm0"
	tpm_version_major      string = "/sys/class/tpm/tpm0/tpm_version_major"
	tpm_device_description string = "/sys/class/tpm/tpm0/device/description"
)

type TpmInfo struct {
	TpmVersionMajor      string
	TpmDeviceDescription string
}

type Tpm interface {
	GetTpm() *TpmInfo
}

// GetTpm return info about the computers TPM, if it exists. If no
// TPM exists nil is returned
func GetTpm(r SysReader) *TpmInfo {
	// check if tmp dir exists
	_, err := os.Stat(tpm_base)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return &TpmInfo{
		TpmVersionMajor:      getTpmInfo(r, tpm_version_major),
		TpmDeviceDescription: getTpmInfo(r, tpm_device_description),
	}
}

// TODO: this is common logic that could be extracted into common logic
func getTpmInfo(r SysReader, path string) string {
	return r.Read(path)
}
