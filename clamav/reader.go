package clamav

import (
	gosysinfo "github.com/DavidHoenisch/go-sysinfo"
	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

const (
	defaultClamdConf     = "/etc/clamav/clamd.conf"
	defaultFreshclamConf = "/etc/clamav/freshclam.conf"
	defaultLocalSocket   = "/run/clamav/clamd.ctl"
)

type Reader struct {
	FS     gosysinfo.SysReader
	Socket probe.SocketDialer
	Cmd    probe.CommandRunner
}

func (r Reader) fs() gosysinfo.SysReader {
	if r.FS != nil {
		return r.FS
	}
	return gosysinfo.Reader{}
}

func (r Reader) socket() probe.SocketDialer {
	if r.Socket != nil {
		return r.Socket
	}
	return probe.OSSocketDialer{}
}

func (r Reader) cmd() probe.CommandRunner {
	if r.Cmd != nil {
		return r.Cmd
	}
	return probe.OSCommandRunner{}
}
