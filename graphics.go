package gosysinfo

import (
	"strings"
)

const (
	fb0_base         string = "/sys/class/graphics/fb0"
	fb0_virtual_size string = "/sys/class/graphics/fb0/virtual_size"
	fb0_screen_modes string = "/sys/class/graphics/fb0/modes"
)

type ScreenSize struct {
	X string
	Y string
}

type ScreenMode struct {
	string
}

type Graphics interface {
	GetScreenVirtualSize() ScreenSize
	GetScreenModes() ScreenMode
}

func GetScreenVirtualSize(r SysReader) ScreenSize {
	content := r.Read(fb0_virtual_size)
	if content == "" {
		return ScreenSize{X: "0", Y: "0"}
	}

	res := strings.Split(content, ",")
	if len(res) != 2 {
		return ScreenSize{X: "0", Y: "0"}
	}

	return ScreenSize{
		X: res[0],
		Y: res[1],
	}
}

func GetScreenMode(r SysReader) ScreenMode {
	content := r.Read(fb0_virtual_size)
	if content == "" {
		return ScreenMode{content}
	}

	return ScreenMode{content}
}

func (r Reader) GetScreenVirtualSize() ScreenSize {
	return GetScreenVirtualSize(r)
}

func (r Reader) GetScreenMode() ScreenMode {
	return GetScreenMode(r)
}
