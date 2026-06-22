package omarchy

import "github.com/DavidHoenisch/go-sysinfo/internal/probe"

type Reader struct {
	Cmd probe.CommandRunner
}

func (r Reader) cmd() probe.CommandRunner {
	if r.Cmd != nil {
		return r.Cmd
	}
	return probe.OSCommandRunner{}
}
