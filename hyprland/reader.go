package hyprland

import (
	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

type Reader struct {
	Cmd probe.CommandRunner
	Env probe.EnvReader
}

func (r Reader) cmd() probe.CommandRunner {
	if r.Cmd != nil {
		return r.Cmd
	}
	return probe.OSCommandRunner{}
}

func (r Reader) env() probe.EnvReader {
	if r.Env != nil {
		return r.Env
	}
	return probe.OSEnvReader{}
}
