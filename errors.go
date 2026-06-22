package gosysinfo

import "errors"

var (
	ErrNetworkNotFound     = errors.New("Network not found")
	ErrBlockDeviceNotFound = errors.New("Block device not found")
	ErrNotAvailable        = errors.New("integration not available")
)
