package clamav

import (
	"bytes"
	"io"
	"strings"

	"github.com/DavidHoenisch/go-sysinfo/internal/probe"
)

func clamdCommand(name string) []byte {
	return append(append([]byte("z"), name...), 0)
}

func clamdRoundTrip(dialer probe.SocketDialer, socketPath, command string) (string, bool) {
	conn, err := dialer.Dial("unix", socketPath)
	if err != nil {
		return "", false
	}
	defer conn.Close()

	if _, err := conn.Write(clamdCommand(command)); err != nil {
		return "", false
	}

	var buf bytes.Buffer
	tmp := make([]byte, 4096)
	for {
		n, err := conn.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			if buf.Len() > 0 {
				break
			}
			return "", false
		}
		if n < len(tmp) {
			break
		}
	}

	response := strings.TrimRight(buf.String(), "\x00")
	response = strings.TrimSpace(response)
	if response == "" {
		return "", false
	}
	return response, true
}

func pingClamd(dialer probe.SocketDialer, socketPath string) bool {
	response, ok := clamdRoundTrip(dialer, socketPath, "PING")
	return ok && response == "PONG"
}

func parseStatsSignatures(stats string) string {
	for _, line := range strings.Split(stats, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "signatures:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "signatures:"))
		}
	}
	return ""
}
